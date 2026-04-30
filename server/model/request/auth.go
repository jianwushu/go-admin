package request

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// RefreshTokenRequest 刷新Token请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新Token
}
