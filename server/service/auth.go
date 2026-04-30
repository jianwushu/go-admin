package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
	"go-admin/utils"

	"github.com/google/uuid"
)

// Redis Key 常量
const (
	// Token 存储 Key: jwt:token:{userId}:{uuid} -> accessToken
	// TTL = Token 过期时间
	tokenKeyPrefix = "jwt:token:"

	// 用户当前 Token 映射 Key: jwt:user:{userId} -> uuid
	// 用于单设备登录，记录用户当前有效的 Token UUID
	userTokenKeyPrefix = "jwt:user:"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login 用户登录
func (s *AuthService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 校验用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 校验密码
	if !utils.BcryptCheck(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 Token
	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 将 Token 缓存到 Redis（有状态设计）
	ctx := context.Background()
	tokenUUID := uuid.New().String()
	jwtCfg := global.Config.JWT

	// 1. 删除用户旧的 Token（单设备登录：踢掉之前的登录）
	s.deleteOldUserToken(ctx, user.ID)

	// 2. 存储新的 Token: jwt:token:{userId}:{uuid} -> accessToken
	tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, user.ID, tokenUUID)
	err = global.Redis.Set(ctx, tokenKey, accessToken, time.Duration(jwtCfg.Expire)*time.Second).Err()
	if err != nil {
		return nil, errors.New("缓存Token失败")
	}

	// 3. 记录用户当前 Token UUID: jwt:user:{userId} -> uuid
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, user.ID)
	err = global.Redis.Set(ctx, userTokenKey, tokenUUID, time.Duration(jwtCfg.Refresh)*time.Second).Err()
	if err != nil {
		return nil, errors.New("记录用户Token失败")
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    jwtCfg.Expire,
	}, nil
}

// Logout 用户登出（从 Redis 删除 Token）
func (s *AuthService) Logout(accessToken string) error {
	// 解析 Token 获取用户信息
	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		// Token 已失效，直接返回成功
		return nil
	}

	ctx := context.Background()

	// 删除用户 Token 映射
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, claims.UserID)
	uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
	if err == nil {
		// 删除具体的 Token
		tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, claims.UserID, uuid)
		global.Redis.Del(ctx, tokenKey)
	}

	// 删除用户 Token 映射
	global.Redis.Del(ctx, userTokenKey)

	return nil
}

// RefreshAccessToken 刷新访问Token
func (s *AuthService) RefreshAccessToken(refreshToken string) (*response.LoginResponse, error) {
	// 解析 refreshToken
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("refreshToken 已失效，请重新登录")
	}

	ctx := context.Background()

	// 检查用户是否仍然有效（未被踢下线）
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, claims.UserID)
	uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
	if err != nil {
		return nil, errors.New("refreshToken 已失效，请重新登录")
	}

	// 检查对应的 Token 是否存在
	tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, claims.UserID, uuid)
	exists, err := global.Redis.Exists(ctx, tokenKey).Result()
	if err != nil || exists == 0 {
		return nil, errors.New("refreshToken 已失效，请重新登录")
	}

	// 生成新的 Token 对
	accessToken, newRefreshToken, err := utils.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	jwtCfg := global.Config.JWT

	// 更新 Redis 中的 Token
	err = global.Redis.Set(ctx, tokenKey, accessToken, time.Duration(jwtCfg.Expire)*time.Second).Err()
	if err != nil {
		return nil, errors.New("更新Token失败")
	}

	// 更新用户 Token 映射的过期时间
	err = global.Redis.Expire(ctx, userTokenKey, time.Duration(jwtCfg.Refresh)*time.Second).Err()
	if err != nil {
		return nil, errors.New("更新Token映射失败")
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    jwtCfg.Expire,
	}, nil
}

// GetUserInfo 获取当前登录用户信息
func (s *AuthService) GetUserInfo(userID int64) (*response.UserInfoResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 查询用户角色
	roles, err := s.getUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 查询用户权限
	perms, err := s.getUserPerms(userID)
	if err != nil {
		return nil, err
	}

	return &response.UserInfoResponse{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		DeptID:      user.DeptID,
		Roles:       roles,
		Permissions: perms,
	}, nil
}

// ForceLogout 强制踢用户下线
func (s *AuthService) ForceLogout(userID int64) error {
	ctx := context.Background()

	// 删除用户 Token 映射
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, userID)
	uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
	if err == nil {
		// 删除具体的 Token
		tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, userID, uuid)
		global.Redis.Del(ctx, tokenKey)
	}

	// 删除用户 Token 映射
	global.Redis.Del(ctx, userTokenKey)

	return nil
}

// IsTokenValid 检查 Token 是否在 Redis 中有效
func IsTokenValid(userID int64, accessToken string) bool {
	ctx := context.Background()

	// 获取用户当前的 Token UUID
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, userID)
	uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
	if err != nil {
		return false
	}

	// 检查对应的 Token 是否存在
	tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, userID, uuid)
	exists, err := global.Redis.Exists(ctx, tokenKey).Result()
	if err != nil || exists == 0 {
		return false
	}

	// 验证 Token 内容是否匹配
	storedToken, err := global.Redis.Get(ctx, tokenKey).Result()
	if err != nil {
		return false
	}

	return storedToken == accessToken
}

// deleteOldUserToken 删除用户旧的 Token（单设备登录）
func (s *AuthService) deleteOldUserToken(ctx context.Context, userID int64) {
	userTokenKey := fmt.Sprintf("%s%d", userTokenKeyPrefix, userID)
	uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
	if err == nil {
		// 删除旧的 Token
		tokenKey := fmt.Sprintf("%s%d:%s", tokenKeyPrefix, userID, uuid)
		global.Redis.Del(ctx, tokenKey)
	}
	// 删除用户 Token 映射
	global.Redis.Del(ctx, userTokenKey)
}

// getUserRoles 获取用户角色标识列表
func (s *AuthService) getUserRoles(userID int64) ([]string, error) {
	var roles []string
	prefix := global.Config.TablePrefix

	err := global.DB.Table(prefix+"role r").
		Select("r.code").
		Joins("JOIN "+prefix+"user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL", userID).
		Pluck("r.code", &roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// getUserPerms 获取用户权限标识列表
func (s *AuthService) getUserPerms(userID int64) ([]string, error) {
	var perms []string
	prefix := global.Config.TablePrefix

	// 查询用户关联的角色
	var roleIDs []int64
	global.DB.Table(prefix+"user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return perms, nil
	}

	// 超级管理员拥有所有权限（通配标识）
	if s.isAdminRole(roleIDs) {
		perms = append(perms, "*:*:*")
		return perms, nil
	}

	// 查询角色关联的菜单权限标识（只查按钮类型，type=2）
	err := global.DB.Table(prefix+"menu m").
		Select("DISTINCT m.perms").
		Joins("JOIN "+prefix+"role_menu rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ? AND m.type = 2 AND m.status = 1 AND m.deleted_at IS NULL", roleIDs).
		Where("m.perms IS NOT NULL AND m.perms != ''").
		Pluck("m.perms", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// isAdminRole 检查角色列表中是否包含超级管理员角色
func (s *AuthService) isAdminRole(roleIDs []int64) bool {
	prefix := global.Config.TablePrefix

	var count int64
	err := global.DB.Table(prefix+"role").
		Where("id IN ? AND code = ? AND status = 1 AND deleted_at IS NULL", roleIDs, "admin").
		Count(&count).Error

	return err == nil && count > 0
}

// GetUserMenus 获取用户的菜单列表（目录和菜单，不含按钮）
func (s *AuthService) GetUserMenus(userID int64) ([]entity.Menu, error) {
	var menus []entity.Menu
	prefix := global.Config.TablePrefix

	// 查询用户关联的角色
	var roleIDs []int64
	global.DB.Table(prefix+"user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return menus, nil
	}

	// 查询角色关联的菜单（type=0 目录, type=1 菜单，不含 type=2 按钮）
	err := global.DB.Table(prefix+"menu m").
		Select("DISTINCT m.*").
		Joins("JOIN "+prefix+"role_menu rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ? AND m.type IN (0, 1) AND m.status = 1 AND m.deleted_at IS NULL", roleIDs).
		Order("m.sort ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}
