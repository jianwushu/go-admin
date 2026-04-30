package request

// DeptCreateRequest 创建部门请求
type DeptCreateRequest struct {
	ParentID int64  `json:"parentId"`                       // 父部门ID
	Name     string `json:"name" binding:"required,max=64"` // 部门名称
	Sort     int    `json:"sort"`                           // 排序
	Status   int    `json:"status" binding:"oneof=0 1"`     // 状态：0=禁用 1=正常
	Leader   string `json:"leader" binding:"max=64"`        // 负责人
	Phone    string `json:"phone" binding:"max=20"`         // 联系电话
	Email    string `json:"email" binding:"max=128"`        // 邮箱
}

// DeptUpdateRequest 更新部门请求
type DeptUpdateRequest struct {
	ID       int64  `json:"id" binding:"required"`           // 部门ID
	ParentID int64  `json:"parentId"`                        // 父部门ID
	Name     string `json:"name" binding:"required,max=64"`  // 部门名称
	Sort     int    `json:"sort"`                            // 排序
	Status   int    `json:"status" binding:"oneof=0 1"`      // 状态：0=禁用 1=正常
	Leader   string `json:"leader" binding:"max=64"`         // 负责人
	Phone    string `json:"phone" binding:"max=20"`          // 联系电话
	Email    string `json:"email" binding:"max=128"`         // 邮箱
}
