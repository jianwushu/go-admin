package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
)

// MenuRepository 菜单仓储层
type MenuRepository struct {
	BaseRepository
}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByID 根据ID查询菜单
func (r *MenuRepository) FindByID(id int64) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindAll 查询所有菜单（用于构建树形结构）
func (r *MenuRepository) FindAll() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// FindByParentID 根据父菜单ID查询子菜单
func (r *MenuRepository) FindByParentID(parentID int64) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.DB.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// FindByIDs 根据ID列表查询菜单
func (r *MenuRepository) FindByIDs(ids []int64) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.DB.Where("id IN ?", ids).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *entity.Menu) error {
	return r.DB.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *entity.Menu) error {
	return r.DB.Save(menu).Error
}

// Delete 删除菜单（软删除）
func (r *MenuRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.Menu{}, id).Error
}

// HasChildren 检查菜单是否有子菜单
func (r *MenuRepository) HasChildren(id int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// IsMenuAssignedToRole 检查菜单是否已分配给角色
func (r *MenuRepository) IsMenuAssignedToRole(menuID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.RoleMenu{}).Where("menu_id = ?", menuID).Count(&count).Error
	return count > 0, err
}

// FindByRoleIDs 根据角色ID列表查询菜单（用于获取角色的菜单权限）
func (r *MenuRepository) FindByRoleIDs(roleIDs []int64) ([]entity.Menu, error) {
	var menus []entity.Menu
	roleMenuTable := entity.RoleMenu{}.TableName()
	menuTable := entity.Menu{}.TableName()
	err := r.DB.Table(menuTable).
		Distinct(menuTable+".*").
		Joins("JOIN "+roleMenuTable+" ON "+menuTable+".id = "+roleMenuTable+".menu_id").
		Where(roleMenuTable+".role_id IN ?", roleIDs).
		Order(menuTable + ".sort ASC, " + menuTable + ".id ASC").
		Find(&menus).Error
	return menus, err
}
