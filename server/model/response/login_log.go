package response

// LoginLogResponse 登录日志响应
type LoginLogResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Status    int    `json:"status"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"createdAt"`
}
