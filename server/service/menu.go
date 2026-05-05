package service

import (
	"errors"

	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// MenuService 菜单服务
type MenuService struct {
	menuRepo *repository.MenuRepository
}

// NewMenuService 创建菜单服务实例
func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo: repository.NewMenuRepository(),
	}
}

// GetTree 获取菜单树形列表
func (s *MenuService) GetTree() ([]response.MenuResponse, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.buildTree(menus, 0), nil
}

// GetByID 根据ID获取菜单详情
func (s *MenuService) GetByID(id int64) (*response.MenuResponse, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("菜单不存在")
	}
	return s.toMenuResponse(menu), nil
}

// Create 创建菜单
func (s *MenuService) Create(req request.MenuCreateRequest) error {
	// 如果有父菜单，检查父菜单是否存在
	if req.ParentID > 0 {
		parent, err := s.menuRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		// 按钮类型不能作为父菜单
		if parent.Type == 2 {
			return errors.New("按钮类型菜单不能作为父菜单")
		}
	}

	// 创建菜单实体
	menu := &entity.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		I18nKey:   req.I18nKey,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Type:      req.Type,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
		Perms:     req.Perms,
	}

	// 默认状态为正常
	if menu.Status == 0 && req.Status != 0 {
		menu.Status = 1
	}
	// 默认可见
	if menu.Visible == 0 && req.Visible != 0 {
		menu.Visible = 1
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return errors.New("创建菜单失败：" + err.Error())
	}

	return nil
}

// Update 更新菜单
func (s *MenuService) Update(req request.MenuUpdateRequest) error {
	// 检查菜单是否存在
	menu, err := s.menuRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 不能将自己设为自己的父菜单
	if req.ParentID == req.ID {
		return errors.New("父菜单不能是自己")
	}

	// 如果有父菜单，检查父菜单是否存在
	if req.ParentID > 0 {
		parent, err := s.menuRepo.FindByID(req.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		// 按钮类型不能作为父菜单
		if parent.Type == 2 {
			return errors.New("按钮类型菜单不能作为父菜单")
		}
		// 不能将子菜单设为父菜单（防止循环引用）
		if s.isChildOf(req.ID, req.ParentID) {
			return errors.New("不能将子菜单设为父菜单")
		}
	}

	// 更新菜单信息
	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.I18nKey = req.I18nKey
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Icon = req.Icon
	menu.Type = req.Type
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	menu.Status = req.Status
	menu.Perms = req.Perms

	if err := s.menuRepo.Update(menu); err != nil {
		return errors.New("更新菜单失败：" + err.Error())
	}

	return nil
}

// Delete 删除菜单
func (s *MenuService) Delete(id int64) error {
	// 检查菜单是否存在
	_, err := s.menuRepo.FindByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否有子菜单
	hasChildren, err := s.menuRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("该菜单下有子菜单，无法删除")
	}

	// 检查菜单是否已分配给角色
	assigned, err := s.menuRepo.IsMenuAssignedToRole(id)
	if err != nil {
		return err
	}
	if assigned {
		return errors.New("该菜单已分配给角色，无法删除")
	}

	// 删除菜单
	if err := s.menuRepo.Delete(id); err != nil {
		return errors.New("删除菜单失败：" + err.Error())
	}

	return nil
}

// GetMenuTreeByRoleIDs 根据角色ID列表获取菜单树（用于角色表单中的菜单选择）
func (s *MenuService) GetMenuTreeByRoleIDs(roleIDs []int64) ([]response.MenuResponse, error) {
	menus, err := s.menuRepo.FindByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	return s.buildTree(menus, 0), nil
}

// buildTree 构建菜单树
func (s *MenuService) buildTree(menus []entity.Menu, parentID int64) []response.MenuResponse {
	var tree []response.MenuResponse
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := response.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				I18nKey:   menu.I18nKey,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Type:      menu.Type,
				Sort:      menu.Sort,
				Visible:   menu.Visible,
				Status:    menu.Status,
				Perms:     menu.Perms,
			}
			children := s.buildTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// toMenuResponse 将菜单实体转换为响应结构
func (s *MenuService) toMenuResponse(menu *entity.Menu) *response.MenuResponse {
	return &response.MenuResponse{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Name:      menu.Name,
		I18nKey:   menu.I18nKey,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		Type:      menu.Type,
		Sort:      menu.Sort,
		Visible:   menu.Visible,
		Status:    menu.Status,
		Perms:     menu.Perms,
	}
}

// isChildOf 检查 targetID 是否是 parentID 的子菜单（递归检查）
func (s *MenuService) isChildOf(parentID, targetID int64) bool {
	children, err := s.menuRepo.FindByParentID(parentID)
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
