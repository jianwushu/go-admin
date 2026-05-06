package response

// OnlineUser 在线用户信息
type OnlineUser struct {
	UserID         int64  `json:"userId"`         // 用户ID
	Username       string `json:"username"`       // 用户名
	Nickname       string `json:"nickname"`       // 昵称
	DeptName       string `json:"deptName"`       // 部门名称
	IP             string `json:"ip"`             // 登录IP
	LoginTime      int64  `json:"loginTime"`      // 登录时间（Unix时间戳）
	OnlineDuration string `json:"onlineDuration"` // 在线时长（格式化字符串）
}
