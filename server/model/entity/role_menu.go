package entity

// RoleMenu 角色菜单关联实体
type RoleMenu struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	RoleID int64 `json:"roleId" gorm:"not null;uniqueIndex:idx_role_menu"`
	MenuID int64 `json:"menuId" gorm:"not null;uniqueIndex:idx_role_menu"`
}

// TableName 返回带前缀的表名
func (RoleMenu) TableName() string {
	return TableName("role_menu")
}
