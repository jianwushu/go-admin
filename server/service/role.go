package service

import (
	"errors"
	"fmt"

	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo *repository.RoleRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
	}
}

// GetList 获取角色列表（分页）
func (s *RoleService) GetList(req request.RoleListRequest) ([]response.RoleResponse, int64, error) {
	var status int
	if req.Status != nil {
		status = *req.Status
	} else {
		status = -1 // -1 表示不筛选状态
	}
	roles, total, err := s.roleRepo.FindWithPage(
		req.GetPage(),
		req.GetPageSize(),
		req.Name,
		req.Code,
		status,
	)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应结构
	var roleResponses []response.RoleResponse
	for _, role := range roles {
		roleResp := s.toRoleResponse(&role)
		roleResponses = append(roleResponses, *roleResp)
	}

	return roleResponses, total, nil
}

// GetAll 获取所有角色（不分页，用于下拉选择）
func (s *RoleService) GetAll() ([]response.RoleResponse, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var roleResponses []response.RoleResponse
	for _, role := range roles {
		roleResp := s.toRoleResponse(&role)
		roleResponses = append(roleResponses, *roleResp)
	}

	return roleResponses, nil
}

// GetByID 根据ID获取角色详情
func (s *RoleService) GetByID(id int64) (*response.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	return s.toRoleResponse(role), nil
}

// Create 创建角色
func (s *RoleService) Create(req request.RoleCreateRequest) error {
	// 检查角色标识是否已存在
	exists, err := s.roleRepo.ExistsByCode(req.Code)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("角色标识已存在")
	}

	// 创建角色实体
	role := &entity.Role{
		Name:      req.Name,
		Code:      req.Code,
		DataScope: req.DataScope,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	// 默认状态为正常
	if role.Status == 0 {
		role.Status = 1
	}

	// admin 角色：代码层级全管理权限和全数据权限，无需维护关联关系
	if req.Code == "admin" {
		role.DataScope = 1 // 全部数据权限
	}

	// 创建角色
	if err := s.roleRepo.Create(role); err != nil {
		return errors.New("创建角色失败：" + err.Error())
	}

	// admin 角色无需维护角色菜单关联和角色部门关联（代码层级自动拥有全部权限）
	if req.Code == "admin" {
		return nil
	}

	// 设置角色菜单关联
	if len(req.MenuIDs) > 0 {
		if err := s.roleRepo.SetRoleMenus(role.ID, req.MenuIDs); err != nil {
			return errors.New("设置角色菜单失败：" + err.Error())
		}
	}

	// 设置角色部门关联（数据权限自定义时）
	if req.DataScope == 5 && len(req.DeptIDs) > 0 {
		if err := s.roleRepo.SetRoleDepts(role.ID, req.DeptIDs); err != nil {
			return errors.New("设置角色部门失败：" + err.Error())
		}
	}

	return nil
}

// Update 更新角色
func (s *RoleService) Update(req request.RoleUpdateRequest) error {
	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查角色标识是否重复（排除自身）
	exists, err := s.roleRepo.ExistsByCodeExcludeID(req.Code, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("角色标识已存在")
	}

	// 更新角色信息
	role.Name = req.Name
	role.Code = req.Code
	role.Sort = req.Sort
	role.Status = req.Status
	role.Remark = req.Remark

	// admin 角色：代码层级全管理权限和全数据权限，强制 DataScope=1
	if role.Code == "admin" {
		role.DataScope = 1
	} else {
		role.DataScope = req.DataScope
	}

	if err := s.roleRepo.Update(role); err != nil {
		return errors.New("更新角色失败：" + err.Error())
	}

	// admin 角色无需维护角色菜单关联和角色部门关联（代码层级自动拥有全部权限）
	if role.Code == "admin" {
		// 清除可能存在的历史关联数据
		_ = s.roleRepo.SetRoleMenus(role.ID, nil)
		_ = s.roleRepo.SetRoleDepts(role.ID, nil)
		return nil
	}

	// 更新角色菜单关联
	if req.MenuIDs != nil {
		if err := s.roleRepo.SetRoleMenus(role.ID, req.MenuIDs); err != nil {
			return errors.New("更新角色菜单失败：" + err.Error())
		}
	}

	// 更新角色部门关联
	if req.DeptIDs != nil {
		if err := s.roleRepo.SetRoleDepts(role.ID, req.DeptIDs); err != nil {
			return errors.New("更新角色部门失败：" + err.Error())
		}
	}

	return nil
}

// Delete 删除角色
func (s *RoleService) Delete(id int64) error {
	// 检查角色是否存在
	_, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 不允许删除超级管理员角色
	if id == 1 {
		return errors.New("不允许删除超级管理员角色")
	}

	// 检查角色是否已分配给用户
	assigned, err := s.roleRepo.IsRoleAssignedToUser(id)
	if err != nil {
		return err
	}
	if assigned {
		return errors.New("该角色已分配给用户，无法删除")
	}

	// 删除角色
	if err := s.roleRepo.Delete(id); err != nil {
		return errors.New("删除角色失败：" + err.Error())
	}

	// 清除角色菜单关联
	_ = s.roleRepo.SetRoleMenus(id, nil)

	// 清除角色部门关联
	_ = s.roleRepo.SetRoleDepts(id, nil)

	return nil
}

// ChangeStatus 修改角色状态
func (s *RoleService) ChangeStatus(req request.RoleChangeStatusRequest) error {
	// 检查角色是否存在
	_, err := s.roleRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 不允许禁用超级管理员角色
	if req.ID == 1 && req.Status == 0 {
		return errors.New("不允许禁用超级管理员角色")
	}

	if err := s.roleRepo.UpdateStatus(req.ID, req.Status); err != nil {
		return fmt.Errorf("修改角色状态失败：%v", err)
	}

	return nil
}

// toRoleResponse 将角色实体转换为响应结构
func (s *RoleService) toRoleResponse(role *entity.Role) *response.RoleResponse {
	resp := &response.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		Code:      role.Code,
		DataScope: role.DataScope,
		Sort:      role.Sort,
		Status:    role.Status,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// admin 角色：代码层级全管理权限和全数据权限，无需返回关联数据
	if role.Code == "admin" {
		return resp
	}

	// 查询角色菜单
	menuIDs, err := s.roleRepo.GetRoleMenus(role.ID)
	if err == nil {
		resp.MenuIDs = menuIDs
	}

	// 查询角色部门（数据权限自定义时）
	deptIDs, err := s.roleRepo.GetRoleDepts(role.ID)
	if err == nil {
		resp.DeptIDs = deptIDs
	}

	return resp
}
