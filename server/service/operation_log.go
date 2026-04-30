package service

import (
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// OperationLogService 操作日志服务
type OperationLogService struct {
	repo *repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志服务实例
func NewOperationLogService() *OperationLogService {
	return &OperationLogService{
		repo: repository.NewOperationLogRepository(),
	}
}

// GetList 获取操作日志列表（分页）
func (s *OperationLogService) GetList(req request.OperationLogListRequest) ([]response.OperationLogResponse, int64, error) {
	logs, total, err := s.repo.FindWithPage(req)
	if err != nil {
		return nil, 0, err
	}

	var result []response.OperationLogResponse
	for _, log := range logs {
		result = append(result, s.toResponse(log))
	}

	return result, total, nil
}

// Clear 清空操作日志
func (s *OperationLogService) Clear() error {
	return s.repo.Clear()
}

// toResponse 将实体转换为响应结构
func (s *OperationLogService) toResponse(log entity.OperationLog) response.OperationLogResponse {
	return response.OperationLogResponse{
		ID:           log.ID,
		Module:       log.Module,
		Action:       log.Action,
		Method:       log.Method,
		URL:          log.URL,
		IP:           log.IP,
		Operator:     log.Operator,
		RequestParam: log.RequestParam,
		ResponseData: log.ResponseData,
		Status:       log.Status,
		ErrorMsg:     log.ErrorMsg,
		Duration:     log.Duration,
		CreatedAt:    log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
