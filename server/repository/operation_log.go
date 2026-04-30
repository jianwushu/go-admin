package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/request"
)

// OperationLogRepository 操作日志仓储层
type OperationLogRepository struct {
	BaseRepository
}

// NewOperationLogRepository 创建操作日志仓储实例
func NewOperationLogRepository() *OperationLogRepository {
	return &OperationLogRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindWithPage 分页查询操作日志
func (r *OperationLogRepository) FindWithPage(req request.OperationLogListRequest) ([]entity.OperationLog, int64, error) {
	var logs []entity.OperationLog
	var total int64

	query := r.DB.Model(&entity.OperationLog{})

	// 条件过滤
	if req.Module != "" {
		query = query.Where("module LIKE ?", "%"+req.Module+"%")
	}
	if req.Operator != "" {
		query = query.Where("operator LIKE ?", "%"+req.Operator+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Method != "" {
		query = query.Where("method = ?", req.Method)
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

// Clear 清空操作日志
func (r *OperationLogRepository) Clear() error {
	return r.DB.Where("1 = 1").Delete(&entity.OperationLog{}).Error
}
