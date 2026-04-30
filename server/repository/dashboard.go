package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"time"
)

// DashboardRepository 仪表盘仓储层
type DashboardRepository struct {
	BaseRepository
}

// NewDashboardRepository 创建仪表盘仓储实例
func NewDashboardRepository() *DashboardRepository {
	return &DashboardRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// CountUsers 统计用户总数
func (r *DashboardRepository) CountUsers() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// CountRoles 统计角色总数
func (r *DashboardRepository) CountRoles() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Role{}).Count(&count).Error
	return count, err
}

// CountMenus 统计菜单总数
func (r *DashboardRepository) CountMenus() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Menu{}).Count(&count).Error
	return count, err
}

// CountDepts 统计部门总数
func (r *DashboardRepository) CountDepts() (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Dept{}).Count(&count).Error
	return count, err
}

// CountTodayLogins 统计今日登录次数
func (r *DashboardRepository) CountTodayLogins() (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.DB.Model(&entity.LoginLog{}).
		Where("status = 1 AND DATE(created_at) = ?", today).
		Count(&count).Error
	return count, err
}
