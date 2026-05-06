package request

// SystemConfigListRequest 系统配置列表查询参数
type SystemConfigListRequest struct {
	PageRequest
	ConfigKey  string `json:"configKey" form:"configKey"`   // 配置键
	ConfigType string `json:"configType" form:"configType"` // 配置类型 text/image/json
}

// SystemConfigUpdateRequest 系统配置更新请求
type SystemConfigUpdateRequest struct {
	ID          int64  `json:"id" binding:"required"` // 配置ID
	ConfigValue string `json:"configValue"`           // 配置值
}

// SystemConfigBatchUpdateRequest 批量更新配置请求
type SystemConfigBatchUpdateRequest struct {
	List []SystemConfigUpdateRequest `json:"list" binding:"required,min=1"`
}
