package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
)

// DeptRepository 部门仓储层
type DeptRepository struct {
	BaseRepository
}

// NewDeptRepository 创建部门仓储实例
func NewDeptRepository() *DeptRepository {
	return &DeptRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByID 根据ID查询部门
func (r *DeptRepository) FindByID(id int64) (*entity.Dept, error) {
	var dept entity.Dept
	err := r.DB.Where("id = ?", id).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// FindAll 查询所有部门（用于构建树形结构）
func (r *DeptRepository) FindAll() ([]entity.Dept, error) {
	var depts []entity.Dept
	err := r.DB.Order("sort ASC, id ASC").Find(&depts).Error
	return depts, err
}

// FindByParentID 根据父部门ID查询子部门
func (r *DeptRepository) FindByParentID(parentID int64) ([]entity.Dept, error) {
	var depts []entity.Dept
	err := r.DB.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&depts).Error
	return depts, err
}

// Create 创建部门
func (r *DeptRepository) Create(dept *entity.Dept) error {
	return r.DB.Create(dept).Error
}

// Update 更新部门
func (r *DeptRepository) Update(dept *entity.Dept) error {
	return r.DB.Save(dept).Error
}

// Delete 删除部门（软删除）
func (r *DeptRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.Dept{}, id).Error
}

// HasChildren 检查部门是否有子部门
func (r *DeptRepository) HasChildren(id int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Dept{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// HasUsers 检查部门下是否有用户
func (r *DeptRepository) HasUsers(id int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("dept_id = ?", id).Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查同级部门下是否存在同名部门
func (r *DeptRepository) ExistsByName(parentID int64, name string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Dept{}).Where("parent_id = ? AND name = ?", parentID, name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcludeID 检查同级部门下是否存在同名部门（排除指定ID）
func (r *DeptRepository) ExistsByNameExcludeID(parentID int64, name string, excludeID int64) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.Dept{}).Where("parent_id = ? AND name = ? AND id != ?", parentID, name, excludeID).Count(&count).Error
	return count > 0, err
}
