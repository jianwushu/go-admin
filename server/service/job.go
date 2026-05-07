package service

import (
	"fmt"

	jobcron "go-admin/cron"
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"

	robfigcron "github.com/robfig/cron/v3"
)

// JobService 定时任务服务
type JobService struct {
	jobRepo *repository.JobRepository
	logRepo *repository.JobLogRepository
}

// NewJobService 创建定时任务服务实例
func NewJobService() *JobService {
	return &JobService{
		jobRepo: repository.NewJobRepository(),
		logRepo: repository.NewJobLogRepository(),
	}
}

// GetList 获取任务列表（分页）
func (s *JobService) GetList(req request.JobListRequest) ([]response.JobResponse, int64, error) {
	jobs, total, err := s.jobRepo.FindWithPage(req)
	if err != nil {
		return nil, 0, err
	}
	var result []response.JobResponse
	for _, j := range jobs {
		result = append(result, s.toResponse(j))
	}
	return result, total, nil
}

// GetByID 根据ID获取任务详情
func (s *JobService) GetByID(id int64) (*response.JobResponse, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := s.toResponse(*job)
	return &resp, nil
}

// Create 创建任务
func (s *JobService) Create(req request.JobCreateRequest, createdBy int64) error {
	// 验证 cron 表达式
	if err := s.validateCronExpr(req.CronExpr); err != nil {
		return err
	}

	job := &entity.Job{
		Name:       req.Name,
		JobType:    req.JobType,
		CronExpr:   req.CronExpr,
		FuncName:   req.FuncName,
		HTTPUrl:    req.HTTPUrl,
		HTTPMethod: req.HTTPMethod,
		HTTPBody:   req.HTTPBody,
		Status:     req.Status,
		Remark:     req.Remark,
		CreatedBy:  createdBy,
		UpdatedBy:  createdBy,
	}

	if job.Status == 0 {
		job.Status = 1 // 默认启用
	}

	if err := s.jobRepo.Create(job); err != nil {
		return err
	}

	// 如果是启用状态，添加到 cron 调度器
	if job.Status == 1 {
		manager := jobcron.GetManager()
		if manager != nil {
			cronJob := &jobcron.JobEntity{
				ID:         job.ID,
				Name:       job.Name,
				JobType:    job.JobType,
				CronExpr:   job.CronExpr,
				FuncName:   job.FuncName,
				HTTPUrl:    job.HTTPUrl,
				HTTPMethod: job.HTTPMethod,
				HTTPBody:   job.HTTPBody,
				Status:     job.Status,
			}
			if err := manager.AddJob(cronJob); err != nil {
				return fmt.Errorf("任务创建成功但调度失败: %v", err)
			}
		}
	}

	return nil
}

// Update 更新任务
func (s *JobService) Update(req request.JobUpdateRequest, updatedBy int64) error {
	job, err := s.jobRepo.FindByID(req.ID)
	if err != nil {
		return fmt.Errorf("任务不存在")
	}

	// 验证 cron 表达式
	if err := s.validateCronExpr(req.CronExpr); err != nil {
		return err
	}

	job.Name = req.Name
	job.JobType = req.JobType
	job.CronExpr = req.CronExpr
	job.FuncName = req.FuncName
	job.HTTPUrl = req.HTTPUrl
	job.HTTPMethod = req.HTTPMethod
	job.HTTPBody = req.HTTPBody
	job.Status = req.Status
	job.Remark = req.Remark
	job.UpdatedBy = updatedBy

	if err := s.jobRepo.Update(job); err != nil {
		return err
	}

	// 更新 cron 调度器
	manager := jobcron.GetManager()
	if manager != nil {
		if job.Status == 1 {
			cronJob := &jobcron.JobEntity{
				ID:         job.ID,
				Name:       job.Name,
				JobType:    job.JobType,
				CronExpr:   job.CronExpr,
				FuncName:   job.FuncName,
				HTTPUrl:    job.HTTPUrl,
				HTTPMethod: job.HTTPMethod,
				HTTPBody:   job.HTTPBody,
				Status:     job.Status,
			}
			if err := manager.UpdateJob(cronJob); err != nil {
				return fmt.Errorf("任务更新成功但调度失败: %v", err)
			}
		} else {
			manager.RemoveJob(job.ID)
		}
	}

	return nil
}

// ChangeStatus 修改任务状态
func (s *JobService) ChangeStatus(req request.JobChangeStatusRequest, updatedBy int64) error {
	job, err := s.jobRepo.FindByID(req.ID)
	if err != nil {
		return fmt.Errorf("任务不存在")
	}

	if err := s.jobRepo.UpdateStatus(req.ID, req.Status, updatedBy); err != nil {
		return err
	}

	// 更新 cron 调度器
	manager := jobcron.GetManager()
	if manager != nil {
		if req.Status == 1 {
			cronJob := &jobcron.JobEntity{
				ID:         job.ID,
				Name:       job.Name,
				JobType:    job.JobType,
				CronExpr:   job.CronExpr,
				FuncName:   job.FuncName,
				HTTPUrl:    job.HTTPUrl,
				HTTPMethod: job.HTTPMethod,
				HTTPBody:   job.HTTPBody,
				Status:     req.Status,
			}
			if err := manager.AddJob(cronJob); err != nil {
				return fmt.Errorf("状态更新成功但调度失败: %v", err)
			}
		} else {
			manager.RemoveJob(job.ID)
		}
	}

	return nil
}

// Delete 删除任务
func (s *JobService) Delete(id int64) error {
	// 先从调度器移除
	manager := jobcron.GetManager()
	if manager != nil {
		manager.RemoveJob(id)
	}

	// 删除任务
	if err := s.jobRepo.Delete(id); err != nil {
		return err
	}

	// 删除关联的日志
	return s.logRepo.DeleteLogsByJobID(id)
}

// DeleteBatch 批量删除任务
func (s *JobService) DeleteBatch(ids []int64) error {
	manager := jobcron.GetManager()
	for _, id := range ids {
		if manager != nil {
			manager.RemoveJob(id)
		}
	}

	if err := s.jobRepo.DeleteBatch(ids); err != nil {
		return err
	}

	// 删除关联的日志
	for _, id := range ids {
		_ = s.logRepo.DeleteLogsByJobID(id)
	}

	return nil
}

// RunOnce 手动执行一次任务
func (s *JobService) RunOnce(id int64) error {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("任务不存在")
	}

	manager := jobcron.GetManager()
	if manager == nil {
		return fmt.Errorf("定时任务管理器未初始化")
	}

	cronJob := &jobcron.JobEntity{
		ID:         job.ID,
		Name:       job.Name,
		JobType:    job.JobType,
		CronExpr:   job.CronExpr,
		FuncName:   job.FuncName,
		HTTPUrl:    job.HTTPUrl,
		HTTPMethod: job.HTTPMethod,
		HTTPBody:   job.HTTPBody,
		Status:     job.Status,
	}

	// 异步执行
	go manager.CreateHandler(cronJob)()

	return nil
}

// ==================== 任务日志 ====================

// GetLogList 获取任务日志列表（分页）
func (s *JobService) GetLogList(req request.JobLogListRequest) ([]response.JobLogResponse, int64, error) {
	logs, total, err := s.logRepo.FindLogsWithPage(req)
	if err != nil {
		return nil, 0, err
	}
	var result []response.JobLogResponse
	for _, l := range logs {
		result = append(result, s.toLogResponse(l))
	}
	return result, total, nil
}

// CleanLogs 清理指定任务的日志
func (s *JobService) CleanLogs(jobID int64) error {
	return s.logRepo.DeleteLogsByJobID(jobID)
}

// validateCronExpr 验证 cron 表达式
func (s *JobService) validateCronExpr(expr string) error {
	parser := robfigcron.NewParser(robfigcron.Second | robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.Dow)
	_, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("无效的cron表达式: %v", err)
	}
	return nil
}

// toResponse 将实体转换为响应结构
func (s *JobService) toResponse(j entity.Job) response.JobResponse {
	return response.JobResponse{
		ID:         j.ID,
		Name:       j.Name,
		JobType:    j.JobType,
		CronExpr:   j.CronExpr,
		FuncName:   j.FuncName,
		HTTPUrl:    j.HTTPUrl,
		HTTPMethod: j.HTTPMethod,
		HTTPBody:   j.HTTPBody,
		Status:     j.Status,
		Remark:     j.Remark,
		CreatedBy:  j.CreatedBy,
		UpdatedBy:  j.UpdatedBy,
		CreatedAt:  j.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  j.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// toLogResponse 将日志实体转换为响应结构
func (s *JobService) toLogResponse(l entity.JobLog) response.JobLogResponse {
	return response.JobLogResponse{
		ID:        l.ID,
		JobID:     l.JobID,
		JobName:   l.JobName,
		Status:    l.Status,
		Result:    l.Result,
		ErrorMsg:  l.ErrorMsg,
		Duration:  l.Duration,
		CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
