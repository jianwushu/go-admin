package response

// DashboardResponse 仪表盘统计数据响应
type DashboardResponse struct {
	UserCount  int64 `json:"userCount"`  // 用户总数
	RoleCount  int64 `json:"roleCount"`  // 角色总数
	MenuCount  int64 `json:"menuCount"`  // 菜单总数
	DeptCount  int64 `json:"deptCount"`  // 部门总数
	TodayLogin int64 `json:"todayLogin"` // 今日登录次数
}

// RecentLoginUser 最近登录用户
type RecentLoginUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	LoginIP   string `json:"loginIp"`
	LoginTime string `json:"loginTime"`
}
