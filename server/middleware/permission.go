package middleware

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// Permission 权限校验中间件
// 根据接口匹配权限标识，校验用户是否有权限访问
// permission 参数格式: "system:user:list" (模块:资源:操作)
func Permission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID（由 Auth 中间件设置）
		userID := GetCurrentUserID(c)
		if userID == 0 {
			utils.Fail(c, 401, "用户未登录")
			c.Abort()
			return
		}

		// 超级管理员拥有所有权限（admin 角色）
		if isSuperAdmin(userID) {
			c.Next()
			return
		}

		// 检查用户是否拥有该权限
		if !hasPermission(userID, permission) {
			utils.Fail(c, 403, "没有操作权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// isSuperAdmin 检查用户是否是超级管理员
func isSuperAdmin(userID int64) bool {
	prefix := global.Config.TablePrefix

	var count int64
	err := global.DB.Table(prefix+"user_role ur").
		Joins("JOIN "+prefix+"role r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.code = ? AND r.status = 1 AND r.deleted_at IS NULL", userID, "admin").
		Count(&count).Error

	return err == nil && count > 0
}

// hasPermission 检查用户是否拥有指定权限
func hasPermission(userID int64, permission string) bool {
	prefix := global.Config.TablePrefix

	// 查询用户关联的角色
	var roleIDs []int64
	global.DB.Table(prefix+"user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return false
	}

	// 查询角色关联的菜单权限
	var count int64
	err := global.DB.Table(prefix+"role_menu rm").
		Joins("JOIN "+prefix+"menu m ON m.id = rm.menu_id").
		Where("rm.role_id IN ? AND m.perms = ? AND m.type = 2 AND m.status = 1 AND m.deleted_at IS NULL", roleIDs, permission).
		Count(&count).Error

	return err == nil && count > 0
}

// GetUserPermissions 获取用户所有权限标识
func GetUserPermissions(userID int64) []string {
	prefix := global.Config.TablePrefix

	var perms []string

	// 查询用户关联的角色
	var roleIDs []int64
	global.DB.Table(prefix+"user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs)

	if len(roleIDs) == 0 {
		return perms
	}

	// 超级管理员拥有所有权限（代码层级全管理权限，无需维护角色菜单关联）
	if isAdminRole(roleIDs) {
		perms = append(perms, "*:*:*")
		return perms
	}

	// 普通角色：查询角色关联的菜单权限标识（只查按钮类型，type=2）
	global.DB.Table(prefix+"menu m").
		Select("DISTINCT m.perms").
		Joins("JOIN "+prefix+"role_menu rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ? AND m.type = 2 AND m.status = 1 AND m.deleted_at IS NULL", roleIDs).
		Where("m.perms IS NOT NULL AND m.perms != ''").
		Pluck("m.perms", &perms)

	return perms
}

// isAdminRole 检查角色列表中是否包含超级管理员角色
func isAdminRole(roleIDs []int64) bool {
	prefix := global.Config.TablePrefix

	var count int64
	err := global.DB.Table(prefix+"role").
		Where("id IN ? AND code = ? AND status = 1 AND deleted_at IS NULL", roleIDs, "admin").
		Count(&count).Error

	return err == nil && count > 0
}

// GetUserRoles 获取用户所有角色标识
func GetUserRoles(userID int64) []string {
	prefix := global.Config.TablePrefix

	var roles []string
	global.DB.Table(prefix+"role r").
		Select("r.code").
		Joins("JOIN "+prefix+"user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL", userID).
		Pluck("r.code", &roles)

	return roles
}

// GetUserRolesWithScope 获取用户角色及其数据权限范围
func GetUserRolesWithScope(userID int64) []entity.Role {
	prefix := global.Config.TablePrefix

	var roles []entity.Role
	global.DB.Table(prefix+"role r").
		Select("r.*").
		Joins("JOIN "+prefix+"user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL", userID).
		Find(&roles)

	return roles
}
