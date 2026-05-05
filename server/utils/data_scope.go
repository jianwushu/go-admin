package utils

import (
	"go-admin/global"
	"go-admin/model/entity"

	"gorm.io/gorm"
)

// DataScope 数据权限范围常量
const (
	DataScopeAll          = 1 // 全部数据
	DataScopeDept         = 2 // 本部门数据
	DataScopeDeptAndChild = 3 // 本部门及下级部门数据
	DataScopeSelf         = 4 // 仅本人数据
	DataScopeCustom       = 5 // 自定义部门数据
)

// DataScopeInfo 数据权限信息
type DataScopeInfo struct {
	UserID    int64         // 用户ID
	DeptID    int64         // 用户部门ID
	DataScope int           // 数据权限范围（取最高权限）
	RoleIDs   []int64       // 用户角色ID列表
	Roles     []entity.Role // 用户角色列表
}

// GetDataScopeInfo 获取用户数据权限信息
// 根据用户ID查询其角色，确定数据权限范围
func GetDataScopeInfo(userID int64) *DataScopeInfo {
	prefix := global.Config.TablePrefix

	info := &DataScopeInfo{
		UserID:    userID,
		DataScope: DataScopeAll, // 默认全部数据
	}

	// 查询用户部门ID
	var user entity.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err == nil {
		info.DeptID = user.DeptID
	}

	// 查询用户角色
	var roles []entity.Role
	global.DB.Table(prefix+"role r").
		Select("r.*").
		Joins("JOIN "+prefix+"user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL", userID).
		Find(&roles)

	if len(roles) == 0 {
		// 没有角色，仅本人数据
		info.DataScope = DataScopeSelf
		return info
	}

	info.Roles = roles
	for _, role := range roles {
		info.RoleIDs = append(info.RoleIDs, role.ID)
	}

	// 超级管理员拥有全部数据权限（代码层级全数据权限，无需维护角色数据权限关联）
	if isAdminRole(info.RoleIDs) {
		info.DataScope = DataScopeAll
		return info
	}

	// 确定数据权限范围（取最小值，即最高权限）
	// 1=全部 > 2=本部门 > 3=本部门及下级 > 4=仅本人 > 5=自定义
	minScope := DataScopeCustom
	for _, role := range roles {
		if role.DataScope > 0 && role.DataScope < minScope {
			minScope = role.DataScope
		}
	}
	info.DataScope = minScope

	return info
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

// ApplyDataScope 对 GORM 查询应用数据权限过滤
// db: 原始查询
// scopeInfo: 数据权限信息
// tableName: 需要过滤的表别名（如 "u" 表示 user 表，空字符串表示使用默认表）
// deptColumn: 部门ID字段名（默认 "dept_id"）
func ApplyDataScope(db *gorm.DB, scopeInfo *DataScopeInfo, tableName string, deptColumn string) *gorm.DB {
	if scopeInfo == nil {
		return db
	}

	// 全部数据，不过滤
	if scopeInfo.DataScope == DataScopeAll {
		return db
	}

	// 构建字段前缀
	deptField := "dept_id"
	if deptColumn != "" {
		deptField = deptColumn
	}
	if tableName != "" {
		deptField = tableName + "." + deptField
	}

	// 创建者字段（用于仅本人数据）
	creatorField := "created_by"
	if tableName != "" {
		creatorField = tableName + "." + creatorField
	}

	switch scopeInfo.DataScope {
	case DataScopeDept:
		// 本部门数据
		return db.Where(deptField+" = ?", scopeInfo.DeptID)

	case DataScopeDeptAndChild:
		// 本部门及下级部门数据
		childDeptIDs := GetChildDeptIDs(scopeInfo.DeptID)
		childDeptIDs = append(childDeptIDs, scopeInfo.DeptID)
		return db.Where(deptField+" IN ?", childDeptIDs)

	case DataScopeSelf:
		// 仅本人数据
		return db.Where(creatorField+" = ?", scopeInfo.UserID)

	case DataScopeCustom:
		// 自定义部门数据
		customDeptIDs := GetCustomDeptIDs(scopeInfo.RoleIDs)
		if len(customDeptIDs) > 0 {
			return db.Where(deptField+" IN ?", customDeptIDs)
		}
		// 没有自定义部门，仅本人数据
		return db.Where(creatorField+" = ?", scopeInfo.UserID)
	}

	return db
}

// GetChildDeptIDs 获取指定部门的所有下级部门ID（递归）
func GetChildDeptIDs(deptID int64) []int64 {
	var childIDs []int64
	var directChildren []int64

	// 查询直接子部门
	global.DB.Model(&entity.Dept{}).Where("parent_id = ?", deptID).Pluck("id", &directChildren)

	childIDs = append(childIDs, directChildren...)

	// 递归查询子部门的子部门
	for _, childID := range directChildren {
		subChildren := GetChildDeptIDs(childID)
		childIDs = append(childIDs, subChildren...)
	}

	return childIDs
}

// GetCustomDeptIDs 获取自定义数据权限的部门ID列表
func GetCustomDeptIDs(roleIDs []int64) []int64 {
	if len(roleIDs) == 0 {
		return nil
	}

	var deptIDs []int64
	global.DB.Model(&entity.RoleDept{}).Where("role_id IN ?", roleIDs).Pluck("dept_id", &deptIDs)

	return deptIDs
}
