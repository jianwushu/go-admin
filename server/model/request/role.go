package request

// RoleListRequest 角色列表查询请求
type RoleListRequest struct {
	PageRequest
	Name   string `json:"name" form:"name"`     // 角色名称（模糊查询）
	Code   string `json:"code" form:"code"`     // 角色标识（模糊查询）
	Status *int   `json:"status" form:"status"` // 状态：0=禁用 1=正常（指针类型，nil表示不筛选）
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name      string  `json:"name" binding:"required,max=64"`               // 角色名称
	Code      string  `json:"code" binding:"required,max=64"`               // 角色标识
	DataScope int     `json:"dataScope" binding:"required,oneof=1 2 3 4 5"` // 数据权限范围
	Sort      int     `json:"sort"`                                         // 排序
	Status    int     `json:"status" binding:"oneof=0 1"`                   // 状态：0=禁用 1=正常
	Remark    string  `json:"remark" binding:"max=512"`                     // 备注
	MenuIDs   []int64 `json:"menuIds"`                                      // 菜单ID列表
	DeptIDs   []int64 `json:"deptIds"`                                      // 部门ID列表（数据权限自定义时使用）
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	ID        int64   `json:"id" binding:"required"`                        // 角色ID
	Name      string  `json:"name" binding:"required,max=64"`               // 角色名称
	Code      string  `json:"code" binding:"required,max=64"`               // 角色标识
	DataScope int     `json:"dataScope" binding:"required,oneof=1 2 3 4 5"` // 数据权限范围
	Sort      int     `json:"sort"`                                         // 排序
	Status    int     `json:"status" binding:"oneof=0 1"`                   // 状态：0=禁用 1=正常
	Remark    string  `json:"remark" binding:"max=512"`                     // 备注
	MenuIDs   []int64 `json:"menuIds"`                                      // 菜单ID列表
	DeptIDs   []int64 `json:"deptIds"`                                      // 部门ID列表（数据权限自定义时使用）
}

// RoleChangeStatusRequest 修改角色状态请求
type RoleChangeStatusRequest struct {
	ID     int64 `json:"id" binding:"required"`      // 角色ID
	Status int   `json:"status" binding:"oneof=0 1"` // 状态：0=禁用 1=正常
}
