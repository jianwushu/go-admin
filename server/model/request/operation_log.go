package request

// OperationLogListRequest 操作日志列表查询参数
type OperationLogListRequest struct {
	PageRequest
	Module    string `json:"module" form:"module"`       // 模块名
	Operator  string `json:"operator" form:"operator"`   // 操作人
	Status    *int   `json:"status" form:"status"`       // 状态 0=失败 1=成功
	Method    string `json:"method" form:"method"`       // 请求方法
	StartTime string `json:"startTime" form:"startTime"` // 开始时间
	EndTime   string `json:"endTime" form:"endTime"`     // 结束时间
}
