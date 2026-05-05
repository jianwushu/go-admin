package response

// UserResponse 用户响应
type UserResponse struct {
	ID        int64       `json:"id"`        // 用户ID
	Username  string      `json:"username"`  // 用户名
	Nickname  string      `json:"nickname"`  // 昵称
	Email     string      `json:"email"`     // 邮箱
	Phone     string      `json:"phone"`     // 手机号
	Status    int         `json:"status"`    // 状态：0=禁用 1=正常
	DeptID    int64       `json:"deptId"`    // 部门ID
	DeptName  string      `json:"deptName"`  // 部门名称
	Avatar    string      `json:"avatar"`    // 头像
	Remark    string      `json:"remark"`    // 备注
	Roles     []RoleBrief `json:"roles"`     // 角色列表
	CreatedAt string      `json:"createdAt"` // 创建时间
}

// RoleBrief 角色简要信息（用于用户详情中的角色列表）
type RoleBrief struct {
	ID   int64  `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名称
	Code string `json:"code"` // 角色标识
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID        int64   `json:"id"`        // 角色ID
	Name      string  `json:"name"`      // 角色名称
	Code      string  `json:"code"`      // 角色标识
	DataScope int     `json:"dataScope"` // 数据权限范围
	Sort      int     `json:"sort"`      // 排序
	Status    int     `json:"status"`    // 状态：0=禁用 1=正常
	Remark    string  `json:"remark"`    // 备注
	MenuIDs   []int64 `json:"menuIds"`   // 菜单ID列表
	DeptIDs   []int64 `json:"deptIds"`   // 部门ID列表
	CreatedAt string  `json:"createdAt"` // 创建时间
}

// UserProfileResponse 个人资料响应
type UserProfileResponse struct {
	ID        int64  `json:"id"`        // 用户ID
	Username  string `json:"username"`  // 用户名
	Nickname  string `json:"nickname"`  // 昵称
	Email     string `json:"email"`     // 邮箱
	Phone     string `json:"phone"`     // 手机号
	Avatar    string `json:"avatar"`    // 头像
	DeptID    int64  `json:"deptId"`    // 部门ID
	DeptName  string `json:"deptName"`  // 部门名称
	Remark    string `json:"remark"`    // 备注
	CreatedAt string `json:"createdAt"` // 创建时间
}
