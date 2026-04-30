package request

// LoginLogListRequest 登录日志列表查询参数
type LoginLogListRequest struct {
	PageRequest
	Username string `json:"username" form:"username"` // 用户名
	IP       string `json:"ip" form:"ip"`             // IP地址
	Status   *int   `json:"status" form:"status"`     // 状态 0=失败 1=成功
}
