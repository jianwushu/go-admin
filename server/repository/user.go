package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/utils"

	"gorm.io/gorm"
)

// UserRepository 用户仓储层
type UserRepository struct {
	BaseRepository
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByUsername 根据用户名查询用户
func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id int64) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindWithPage 分页查询用户列表（支持数据权限过滤）
func (r *UserRepository) FindWithPage(page, pageSize int, scopeInfo *utils.DataScopeInfo, username string, status int, deptID int64) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.DB.Model(&entity.User{})

	// 应用数据权限过滤
	if scopeInfo != nil {
		query = r.ApplyDataScope(query, scopeInfo, "", "dept_id")
	}

	// 条件过滤
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status == 0 || status == 1 {
		query = query.Where("status = ?", status)
	}
	// status 为 -1 时不筛选状态
	if deptID > 0 {
		query = query.Where("dept_id = ?", deptID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindWithQuery 自定义查询条件的分页查询（支持数据权限）
func (r *UserRepository) FindWithQuery(query *gorm.DB, scopeInfo *utils.DataScopeInfo, page, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	// 应用数据权限过滤
	if scopeInfo != nil {
		query = r.ApplyDataScope(query, scopeInfo, "", "dept_id")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *entity.User) error {
	return r.DB.Save(user).Error
}

// Delete 删除用户（软删除），支持单个或多个ID
func (r *UserRepository) Delete(ids []int64) error {
	return r.DB.Delete(&entity.User{}, ids).Error
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(id int64, status int) error {
	return r.DB.Model(&entity.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(id int64, password string) error {
	return r.DB.Model(&entity.User{}).Where("id = ?", id).Update("password", password).Error
}

// GetUserRoles 获取用户关联的角色ID列表
func (r *UserRepository) GetUserRoles(userID int64) ([]int64, error) {
	var roleIDs []int64
	err := r.DB.Model(&entity.UserRole{}).Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// SetUserRoles 设置用户角色关联（先删后增）
func (r *UserRepository) SetUserRoles(userID int64, roleIDs []int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("user_id = ?", userID).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}

		// 创建新的关联
		if len(roleIDs) > 0 {
			var userRoles []entity.UserRole
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, entity.UserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *UserRepository) ExistsByUsernameExcludeID(username string, excludeID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("username = ? AND id != ?", username, excludeID).Count(&count).Error
	return count > 0, err
}
