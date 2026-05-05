package request

// UserListRequest 用户列表查询请求
type UserListRequest struct {
	PageRequest
	Username string `json:"username" form:"username"` // 用户名（模糊查询）
	Status   int    `json:"status" form:"status"`     // 状态：0=禁用 1=正常
	DeptID   int64  `json:"deptId" form:"deptId"`     // 部门ID
	Phone    string `json:"phone" form:"phone"`       // 手机号（模糊查询）
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string  `json:"username" binding:"required,min=2,max=64"`  // 用户名
	Password string  `json:"password" binding:"required,min=6,max=128"` // 密码
	Nickname string  `json:"nickname" binding:"max=64"`                 // 昵称
	Email    string  `json:"email" binding:"omitempty,email,max=128"`   // 邮箱
	Phone    string  `json:"phone" binding:"omitempty,max=20"`          // 手机号
	Status   int     `json:"status" binding:"oneof=0 1"`                // 状态：0=禁用 1=正常
	DeptID   int64   `json:"deptId"`                                    // 部门ID
	Avatar   string  `json:"avatar" binding:"max=256"`                  // 头像
	Remark   string  `json:"remark" binding:"max=512"`                  // 备注
	RoleIDs  []int64 `json:"roleIds"`                                   // 角色ID列表
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	ID       int64   `json:"id" binding:"required"`                     // 用户ID
	Nickname string  `json:"nickname" binding:"max=64"`                 // 昵称
	Email    string  `json:"email" binding:"omitempty,email,max=128"`   // 邮箱
	Phone    string  `json:"phone" binding:"omitempty,max=20"`          // 手机号
	Status   int     `json:"status" binding:"oneof=0 1"`                // 状态：0=禁用 1=正常
	DeptID   int64   `json:"deptId"`                                    // 部门ID
	Avatar   string  `json:"avatar" binding:"max=256"`                  // 头像
	Remark   string  `json:"remark" binding:"max=512"`                  // 备注
	RoleIDs  []int64 `json:"roleIds"`                                   // 角色ID列表
}

// UserResetPasswordRequest 重置密码请求
type UserResetPasswordRequest struct {
	ID       int64  `json:"id" binding:"required"`              // 用户ID
	Password string `json:"password" binding:"required,min=6,max=128"` // 新密码
}

// UserChangeStatusRequest 修改用户状态请求
type UserChangeStatusRequest struct {
	ID     int64 `json:"id" binding:"required"`      // 用户ID
	Status int   `json:"status" binding:"oneof=0 1"` // 状态：0=禁用 1=正常
}

// UserProfileUpdateRequest 个人资料更新请求
type UserProfileUpdateRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`               // 昵称
	Email    string `json:"email" binding:"omitempty,email,max=128"` // 邮箱
	Phone    string `json:"phone" binding:"omitempty,max=20"`        // 手机号
	Avatar   string `json:"avatar" binding:"max=256"`                // 头像
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=128"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required,min=6,max=128"` // 新密码
}
