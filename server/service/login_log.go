package service

import (
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// LoginLogService 登录日志服务
type LoginLogService struct {
	repo *repository.LoginLogRepository
}

// NewLoginLogService 创建登录日志服务实例
func NewLoginLogService() *LoginLogService {
	return &LoginLogService{
		repo: repository.NewLoginLogRepository(),
	}
}

// Create 创建登录日志
func (s *LoginLogService) Create(log *entity.LoginLog) error {
	return s.repo.Create(log)
}

// GetList 获取登录日志列表（分页）
func (s *LoginLogService) GetList(req request.LoginLogListRequest) ([]response.LoginLogResponse, int64, error) {
	logs, total, err := s.repo.FindWithPage(req)
	if err != nil {
		return nil, 0, err
	}

	var result []response.LoginLogResponse
	for _, log := range logs {
		result = append(result, s.toResponse(log))
	}

	return result, total, nil
}

// Clear 清空登录日志
func (s *LoginLogService) Clear() error {
	return s.repo.Clear()
}

// toResponse 将实体转换为响应结构
func (s *LoginLogService) toResponse(log entity.LoginLog) response.LoginLogResponse {
	return response.LoginLogResponse{
		ID:        log.ID,
		Username:  log.Username,
		IP:        log.IP,
		Location:  log.Location,
		Browser:   log.Browser,
		OS:        log.OS,
		Status:    log.Status,
		Msg:       log.Msg,
		CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
