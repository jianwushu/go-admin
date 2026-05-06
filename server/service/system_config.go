package service

import (
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
)

// SystemConfigService 系统配置服务
type SystemConfigService struct {
	repo *repository.SystemConfigRepository
}

// NewSystemConfigService 创建系统配置服务实例
func NewSystemConfigService() *SystemConfigService {
	return &SystemConfigService{
		repo: repository.NewSystemConfigRepository(),
	}
}

// GetAll 获取所有配置
func (s *SystemConfigService) GetAll() ([]response.SystemConfigResponse, error) {
	configs, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []response.SystemConfigResponse
	for _, c := range configs {
		result = append(result, s.toResponse(c))
	}
	return result, nil
}

// GetList 获取配置列表（分页）
func (s *SystemConfigService) GetList(req request.SystemConfigListRequest) ([]response.SystemConfigResponse, int64, error) {
	configs, total, err := s.repo.FindWithPage(req)
	if err != nil {
		return nil, 0, err
	}
	var result []response.SystemConfigResponse
	for _, c := range configs {
		result = append(result, s.toResponse(c))
	}
	return result, total, nil
}

// GetByKey 根据配置键获取配置值
func (s *SystemConfigService) GetByKey(key string) (*response.SystemConfigResponse, error) {
	config, err := s.repo.FindByKey(key)
	if err != nil {
		return nil, err
	}
	resp := s.toResponse(*config)
	return &resp, nil
}

// GetByKeys 根据多个配置键批量获取配置
func (s *SystemConfigService) GetByKeys(keys []string) (map[string]string, error) {
	configs, err := s.repo.FindByKeys(keys)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(configs))
	for _, c := range configs {
		result[c.ConfigKey] = c.ConfigValue
	}
	return result, nil
}

// Update 更新单个配置
func (s *SystemConfigService) Update(req request.SystemConfigUpdateRequest, updatedBy int64) error {
	return s.repo.Update(req.ID, req.ConfigValue, updatedBy)
}

// BatchUpdate 批量更新配置
func (s *SystemConfigService) BatchUpdate(req request.SystemConfigBatchUpdateRequest, updatedBy int64) error {
	updates := make([]map[string]interface{}, 0, len(req.List))
	for _, item := range req.List {
		updates = append(updates, map[string]interface{}{
			"id":           item.ID,
			"config_value": item.ConfigValue,
			"updated_by":   updatedBy,
		})
	}
	return s.repo.BatchUpdate(updates)
}

// toResponse 将实体转换为响应结构
func (s *SystemConfigService) toResponse(c entity.SystemConfig) response.SystemConfigResponse {
	return response.SystemConfigResponse{
		ID:          c.ID,
		ConfigKey:   c.ConfigKey,
		ConfigValue: c.ConfigValue,
		ConfigType:  c.ConfigType,
		Remark:      c.Remark,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
