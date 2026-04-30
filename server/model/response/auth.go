package response

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`  // 访问Token
	RefreshToken string `json:"refreshToken"` // 刷新Token
	ExpiresIn    int    `json:"expiresIn"`    // 过期时间（秒）
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          int64    `json:"id"`          // 用户ID
	Username    string   `json:"username"`    // 用户名
	Nickname    string   `json:"nickname"`    // 昵称
	Email       string   `json:"email"`       // 邮箱
	Phone       string   `json:"phone"`       // 手机号
	Avatar      string   `json:"avatar"`      // 头像
	DeptID      int64    `json:"deptId"`      // 部门ID
	Roles       []string `json:"roles"`       // 角色标识列表
	Permissions []string `json:"permissions"` // 权限标识列表
}
