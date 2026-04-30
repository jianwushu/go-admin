package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/request"
)

// LoginLogRepository 登录日志仓储层
type LoginLogRepository struct {
	BaseRepository
}

// NewLoginLogRepository 创建登录日志仓储实例
func NewLoginLogRepository() *LoginLogRepository {
	return &LoginLogRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// Create 创建登录日志
func (r *LoginLogRepository) Create(log *entity.LoginLog) error {
	return r.DB.Create(log).Error
}

// FindWithPage 分页查询登录日志
func (r *LoginLogRepository) FindWithPage(req request.LoginLogListRequest) ([]entity.LoginLog, int64, error) {
	var logs []entity.LoginLog
	var total int64

	query := r.DB.Model(&entity.LoginLog{})

	// 条件过滤
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.IP != "" {
		query = query.Where("ip LIKE ?", "%"+req.IP+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	if err := query.Offset(offset).Limit(req.GetPageSize()).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// Clear 清空登录日志
func (r *LoginLogRepository) Clear() error {
	return r.DB.Where("1 = 1").Delete(&entity.LoginLog{}).Error
}
