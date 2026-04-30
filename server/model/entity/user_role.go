package entity

// UserRole 用户角色关联实体
type UserRole struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"userId" gorm:"not null;uniqueIndex:idx_user_role"`
	RoleID int64 `json:"roleId" gorm:"not null;uniqueIndex:idx_user_role"`
}

// TableName 返回带前缀的表名
func (UserRole) TableName() string {
	return TableName("user_role")
}
