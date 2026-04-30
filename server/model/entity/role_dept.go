package entity

// RoleDept 角色部门关联实体（数据权限自定义用）
type RoleDept struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	RoleID int64 `json:"roleId" gorm:"not null;uniqueIndex:idx_role_dept"`
	DeptID int64 `json:"deptId" gorm:"not null;uniqueIndex:idx_role_dept"`
}

// TableName 返回带前缀的表名
func (RoleDept) TableName() string {
	return TableName("role_dept")
}
