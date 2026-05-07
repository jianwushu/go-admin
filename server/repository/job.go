package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/request"
)

// JobRepository 定时任务仓储层
type JobRepository struct {
	BaseRepository
}

// NewJobRepository 创建定时任务仓储实例
func NewJobRepository() *JobRepository {
	return &JobRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByID 根据ID查询任务
func (r *JobRepository) FindByID(id int64) (*entity.Job, error) {
	var job entity.Job
	if err := r.DB.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// FindWithPage 分页查询任务
func (r *JobRepository) FindWithPage(req request.JobListRequest) ([]entity.Job, int64, error) {
	var jobs []entity.Job
	var total int64

	query := r.DB.Model(&entity.Job{})

	// 条件过滤
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.JobType > 0 {
		query = query.Where("job_type = ?", req.JobType)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	if err := query.Offset(offset).Limit(req.GetPageSize()).Order("id DESC").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// FindAllEnabled 查询所有启用的任务（用于启动时加载）
func (r *JobRepository) FindAllEnabled() ([]entity.Job, error) {
	var jobs []entity.Job
	if err := r.DB.Where("status = ?", 1).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// Create 创建任务
func (r *JobRepository) Create(job *entity.Job) error {
	return r.DB.Create(job).Error
}

// Update 更新任务
func (r *JobRepository) Update(job *entity.Job) error {
	return r.DB.Save(job).Error
}

// UpdateStatus 更新任务状态
func (r *JobRepository) UpdateStatus(id int64, status int, updatedBy int64) error {
	return r.DB.Model(&entity.Job{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_by": updatedBy,
		}).Error
}

// Delete 删除任务（软删除）
func (r *JobRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.Job{}, id).Error
}

// DeleteBatch 批量删除任务（软删除）
func (r *JobRepository) DeleteBatch(ids []int64) error {
	return r.DB.Delete(&entity.Job{}, ids).Error
}

// ==================== 任务日志 ====================

// JobLogRepository 任务日志仓储层
type JobLogRepository struct {
	BaseRepository
}

// NewJobLogRepository 创建任务日志仓储实例
func NewJobLogRepository() *JobLogRepository {
	return &JobLogRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// CreateLog 创建任务执行日志
func (r *JobLogRepository) CreateLog(log *entity.JobLog) error {
	return r.DB.Create(log).Error
}

// FindLogsWithPage 分页查询任务日志
func (r *JobLogRepository) FindLogsWithPage(req request.JobLogListRequest) ([]entity.JobLog, int64, error) {
	var logs []entity.JobLog
	var total int64

	query := r.DB.Model(&entity.JobLog{})

	// 条件过滤
	if req.JobID > 0 {
		query = query.Where("job_id = ?", req.JobID)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
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

// DeleteLogsByJobID 删除指定任务的所有日志
func (r *JobLogRepository) DeleteLogsByJobID(jobID int64) error {
	return r.DB.Where("job_id = ?", jobID).Delete(&entity.JobLog{}).Error
}

// CleanOldLogs 清理N天前的日志
func (r *JobLogRepository) CleanOldLogs(days int) (int64, error) {
	result := r.DB.Where("created_at < datetime('now', ?)", "-"+string(rune(days+'0'))+" days").
		Delete(&entity.JobLog{})
	return result.RowsAffected, result.Error
}
