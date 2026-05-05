package repository

import (
	"go-admin/global"
	"go-admin/model/entity"

	"gorm.io/gorm"
)

// RoleRepository 角色仓储层
type RoleRepository struct {
	BaseRepository
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByID 根据ID查询角色
func (r *RoleRepository) FindByID(id int64) (*entity.Role, error) {
	var role entity.Role
	err := r.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByCode 根据角色标识查询角色
func (r *RoleRepository) FindByCode(code string) (*entity.Role, error) {
	var role entity.Role
	err := r.DB.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindAll 查询所有角色（不分页）
func (r *RoleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.DB.Order("sort ASC, id ASC").Find(&roles).Error
	return roles, err
}

// FindWithPage 分页查询角色列表
func (r *RoleRepository) FindWithPage(page, pageSize int, name string, code string, status int) ([]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	query := r.DB.Model(&entity.Role{})

	// 条件过滤
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status == 0 || status == 1 {
		query = query.Where("status = ?", status)
	}
	// status 为 -1 时不筛选状态

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Create 创建角色
func (r *RoleRepository) Create(role *entity.Role) error {
	return r.DB.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *entity.Role) error {
	return r.DB.Save(role).Error
}

// Delete 删除角色（软删除）
func (r *RoleRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.Role{}, id).Error
}

// DeleteBatch 批量删除角色（软删除）
func (r *RoleRepository) DeleteBatch(ids []int64) error {
	return r.DB.Delete(&entity.Role{}, ids).Error
}

// UpdateStatus 更新角色状态
func (r *RoleRepository) UpdateStatus(id int64, status int) error {
	return r.DB.Model(&entity.Role{}).Where("id = ?", id).Update("status", status).Error
}

// GetRoleMenus 获取角色关联的菜单ID列表
func (r *RoleRepository) GetRoleMenus(roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := r.DB.Model(&entity.RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// SetRoleMenus 设置角色菜单关联（先删后增）
func (r *RoleRepository) SetRoleMenus(roleID int64, menuIDs []int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}

		// 创建新的关联
		if len(menuIDs) > 0 {
			var roleMenus []entity.RoleMenu
			for _, menuID := range menuIDs {
				roleMenus = append(roleMenus, entity.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				})
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRoleDepts 获取角色关联的部门ID列表（数据权限自定义用）
func (r *RoleRepository) GetRoleDepts(roleID int64) ([]int64, error) {
	var deptIDs []int64
	err := r.DB.Model(&entity.RoleDept{}).Where("role_id = ?", roleID).Pluck("dept_id", &deptIDs).Error
	return deptIDs, err
}

// SetRoleDepts 设置角色部门关联（先删后增）
func (r *RoleRepository) SetRoleDepts(roleID int64, deptIDs []int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RoleDept{}).Error; err != nil {
			return err
		}

		// 创建新的关联
		if len(deptIDs) > 0 {
			var roleDepts []entity.RoleDept
			for _, deptID := range deptIDs {
				roleDepts = append(roleDepts, entity.RoleDept{
					RoleID: roleID,
					DeptID: deptID,
				})
			}
			if err := tx.Create(&roleDepts).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ExistsByCode 检查角色标识是否存在
func (r *RoleRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Role{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// ExistsByCodeExcludeID 检查角色标识是否存在（排除指定ID）
func (r *RoleRepository) ExistsByCodeExcludeID(code string, excludeID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Role{}).Where("code = ? AND id != ?", code, excludeID).Count(&count).Error
	return count > 0, err
}

// IsRoleAssignedToUser 检查角色是否已分配给用户
func (r *RoleRepository) IsRoleAssignedToUser(roleID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	return count > 0, err
}
