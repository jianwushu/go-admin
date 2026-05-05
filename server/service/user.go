package service

import (
	"errors"
	"fmt"

	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
	"go-admin/utils"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
		roleRepo: repository.NewRoleRepository(),
	}
}

// GetList 获取用户列表（分页）
func (s *UserService) GetList(req request.UserListRequest, scopeInfo *utils.DataScopeInfo) ([]response.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindWithPage(
		req.GetPage(),
		req.GetPageSize(),
		scopeInfo,
		req.Username,
		req.Status,
		req.DeptID,
	)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应结构
	var userResponses []response.UserResponse
	for _, user := range users {
		userResp := s.toUserResponse(&user)
		userResponses = append(userResponses, *userResp)
	}

	return userResponses, total, nil
}

// GetByID 根据ID获取用户详情
func (s *UserService) GetByID(id int64) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return s.toUserResponse(user), nil
}

// Create 创建用户
func (s *UserService) Create(req request.UserCreateRequest) error {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.BcryptHash(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 创建用户实体
	user := &entity.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		DeptID:   req.DeptID,
		Avatar:   req.Avatar,
		Remark:   req.Remark,
	}

	// 默认状态为正常
	if user.Status == 0 {
		user.Status = 1
	}

	// 创建用户
	if err := s.userRepo.Create(user); err != nil {
		return errors.New("创建用户失败：" + err.Error())
	}

	// 设置用户角色关联
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.SetUserRoles(user.ID, req.RoleIDs); err != nil {
			return errors.New("设置用户角色失败：" + err.Error())
		}
	}

	return nil
}

// Update 更新用户
func (s *UserService) Update(req request.UserUpdateRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Status = req.Status
	user.DeptID = req.DeptID
	user.Avatar = req.Avatar
	user.Remark = req.Remark

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户失败：" + err.Error())
	}

	// 更新用户角色关联
	if req.RoleIDs != nil {
		if err := s.userRepo.SetUserRoles(user.ID, req.RoleIDs); err != nil {
			return errors.New("更新用户角色失败：" + err.Error())
		}
	}

	return nil
}

// Delete 删除用户
func (s *UserService) Delete(id int64) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 不允许删除超级管理员
	if id == 1 {
		return errors.New("不允许删除超级管理员")
	}

	// 删除用户
	if err := s.userRepo.Delete(id); err != nil {
		return errors.New("删除用户失败：" + err.Error())
	}

	// 清除用户角色关联
	_ = s.userRepo.SetUserRoles(id, nil)

	return nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(req request.UserResetPasswordRequest) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 加密密码
	hashedPassword, err := utils.BcryptHash(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(req.ID, hashedPassword); err != nil {
		return errors.New("重置密码失败：" + err.Error())
	}

	return nil
}

// ChangeStatus 修改用户状态
func (s *UserService) ChangeStatus(req request.UserChangeStatusRequest) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 不允许禁用超级管理员
	if req.ID == 1 && req.Status == 0 {
		return errors.New("不允许禁用超级管理员")
	}

	if err := s.userRepo.UpdateStatus(req.ID, req.Status); err != nil {
		return fmt.Errorf("修改用户状态失败：%v", err)
	}

	return nil
}

// GetProfile 获取当前用户的个人资料
func (s *UserService) GetProfile(userID int64) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	profile := &response.UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		DeptID:    user.DeptID,
		Remark:    user.Remark,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 查询部门名称
	if user.DeptID > 0 {
		deptRepo := repository.NewDeptRepository()
		dept, err := deptRepo.FindByID(user.DeptID)
		if err == nil {
			profile.DeptName = dept.Name
		}
	}

	return profile, nil
}

// UpdateProfile 更新当前用户的个人资料
func (s *UserService) UpdateProfile(userID int64, req request.UserProfileUpdateRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Avatar = req.Avatar

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新个人资料失败：" + err.Error())
	}

	return nil
}

// ChangePassword 修改当前用户的密码
func (s *UserService) ChangePassword(userID int64, req request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.BcryptCheck(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.BcryptHash(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return errors.New("修改密码失败：" + err.Error())
	}

	return nil
}

// toUserResponse 将用户实体转换为响应结构
func (s *UserService) toUserResponse(user *entity.User) *response.UserResponse {
	resp := &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		DeptID:    user.DeptID,
		Avatar:    user.Avatar,
		Remark:    user.Remark,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 查询用户角色
	roleIDs, err := s.userRepo.GetUserRoles(user.ID)
	if err == nil && len(roleIDs) > 0 {
		for _, roleID := range roleIDs {
			role, err := s.roleRepo.FindByID(roleID)
			if err == nil {
				resp.Roles = append(resp.Roles, response.RoleBrief{
					ID:   role.ID,
					Name: role.Name,
					Code: role.Code,
				})
			}
		}
	}

	return resp
}
