package service

import (
	"errors"

	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// DeptService 部门服务
type DeptService struct {
	deptRepo *repository.DeptRepository
}

// NewDeptService 创建部门服务实例
func NewDeptService() *DeptService {
	return &DeptService{
		deptRepo: repository.NewDeptRepository(),
	}
}

// GetTree 获取部门树形列表
func (s *DeptService) GetTree() ([]response.DeptResponse, error) {
	depts, err := s.deptRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.buildTree(depts, 0), nil
}

// GetByID 根据ID获取部门详情
func (s *DeptService) GetByID(id int64) (*response.DeptResponse, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("部门不存在")
	}
	return s.toDeptResponse(dept), nil
}

// Create 创建部门
func (s *DeptService) Create(req request.DeptCreateRequest) error {
	// 如果有父部门，检查父部门是否存在
	if req.ParentID > 0 {
		_, err := s.deptRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父部门不存在")
		}
	}

	// 检查同级部门下是否存在同名部门
	exists, err := s.deptRepo.ExistsByName(req.ParentID, req.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("同级部门下已存在同名部门")
	}

	// 创建部门实体
	dept := &entity.Dept{
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     req.Sort,
		Status:   req.Status,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	// 默认状态为正常
	if dept.Status == 0 && req.Status != 0 {
		dept.Status = 1
	}

	if err := s.deptRepo.Create(dept); err != nil {
		return errors.New("创建部门失败：" + err.Error())
	}

	return nil
}

// Update 更新部门
func (s *DeptService) Update(req request.DeptUpdateRequest) error {
	// 检查部门是否存在
	dept, err := s.deptRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("部门不存在")
	}

	// 不能将自己设为自己的父部门
	if req.ParentID == req.ID {
		return errors.New("父部门不能是自己")
	}

	// 如果有父部门，检查父部门是否存在
	if req.ParentID > 0 {
		_, err := s.deptRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父部门不存在")
		}
		// 不能将子部门设为父部门（防止循环引用）
		if s.isChildOf(req.ID, req.ParentID) {
			return errors.New("不能将子部门设为父部门")
		}
	}

	// 检查同级部门下是否存在同名部门（排除自身）
	exists, err := s.deptRepo.ExistsByNameExcludeID(req.ParentID, req.Name, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("同级部门下已存在同名部门")
	}

	// 更新部门信息
	dept.ParentID = req.ParentID
	dept.Name = req.Name
	dept.Sort = req.Sort
	dept.Status = req.Status
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email

	if err := s.deptRepo.Update(dept); err != nil {
		return errors.New("更新部门失败：" + err.Error())
	}

	return nil
}

// Delete 删除部门
func (s *DeptService) Delete(id int64) error {
	// 检查部门是否存在
	_, err := s.deptRepo.FindByID(id)
	if err != nil {
		return errors.New("部门不存在")
	}

	// 检查是否有子部门
	hasChildren, err := s.deptRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("该部门下有子部门，无法删除")
	}

	// 检查部门下是否有用户
	hasUsers, err := s.deptRepo.HasUsers(id)
	if err != nil {
		return err
	}
	if hasUsers {
		return errors.New("该部门下有用户，无法删除")
	}

	// 删除部门
	if err := s.deptRepo.Delete(id); err != nil {
		return errors.New("删除部门失败：" + err.Error())
	}

	return nil
}

// buildTree 构建部门树
func (s *DeptService) buildTree(depts []entity.Dept, parentID int64) []response.DeptResponse {
	var tree []response.DeptResponse
	for _, dept := range depts {
		if dept.ParentID == parentID {
			node := response.DeptResponse{
				ID:       dept.ID,
				ParentID: dept.ParentID,
				Name:     dept.Name,
				Sort:     dept.Sort,
				Status:   dept.Status,
				Leader:   dept.Leader,
				Phone:    dept.Phone,
				Email:    dept.Email,
			}
			children := s.buildTree(depts, dept.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// toDeptResponse 将部门实体转换为响应结构
func (s *DeptService) toDeptResponse(dept *entity.Dept) *response.DeptResponse {
	return &response.DeptResponse{
		ID:       dept.ID,
		ParentID: dept.ParentID,
		Name:     dept.Name,
		Sort:     dept.Sort,
		Status:   dept.Status,
		Leader:   dept.Leader,
		Phone:    dept.Phone,
		Email:    dept.Email,
	}
}

// isChildOf 检查 targetID 是否是 parentID 的子部门（递归检查）
func (s *DeptService) isChildOf(parentID, targetID int64) bool {
	children, err := s.deptRepo.FindByParentID(parentID)
	if err != nil {
		return false
	}
	for _, child := range children {
		if child.ID == targetID {
			return true
		}
		if s.isChildOf(child.ID, targetID) {
			return true
		}
	}
	return false
}
