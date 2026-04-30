package service

import (
	"go-admin/model/response"
	"go-admin/repository"
)

// DashboardService 仪表盘服务层
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

// NewDashboardService 创建仪表盘服务实例
func NewDashboardService() *DashboardService {
	return &DashboardService{
		dashboardRepo: repository.NewDashboardRepository(),
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (s *DashboardService) GetDashboardStats() (*response.DashboardResponse, error) {
	userCount, err := s.dashboardRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	roleCount, err := s.dashboardRepo.CountRoles()
	if err != nil {
		return nil, err
	}

	menuCount, err := s.dashboardRepo.CountMenus()
	if err != nil {
		return nil, err
	}

	deptCount, err := s.dashboardRepo.CountDepts()
	if err != nil {
		return nil, err
	}

	todayLogin, err := s.dashboardRepo.CountTodayLogins()
	if err != nil {
		return nil, err
	}

	return &response.DashboardResponse{
		UserCount:  userCount,
		RoleCount:  roleCount,
		MenuCount:  menuCount,
		DeptCount:  deptCount,
		TodayLogin: todayLogin,
	}, nil
}
