package response

// DeptResponse 部门响应
type DeptResponse struct {
	ID       int64          `json:"id"`                // 部门ID
	ParentID int64          `json:"parentId"`          // 父部门ID
	Name     string         `json:"name"`              // 部门名称
	Sort     int            `json:"sort"`              // 排序
	Status   int            `json:"status"`            // 状态：0=禁用 1=正常
	Leader   string         `json:"leader"`            // 负责人
	Phone    string         `json:"phone"`             // 联系电话
	Email    string         `json:"email"`             // 邮箱
	Children []DeptResponse `json:"children,omitempty"` // 子部门
}
