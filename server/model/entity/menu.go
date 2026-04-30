package entity

// Menu 菜单实体
type Menu struct {
	BaseModel
	ParentID  int64  `json:"parentId" gorm:"default:0"`
	Name      string `json:"name" gorm:"size:64;not null"`
	Path      string `json:"path" gorm:"size:128"`
	Component string `json:"component" gorm:"size:128"`
	Icon      string `json:"icon" gorm:"size:64"`
	Type      int    `json:"type"` // 0=目录 1=菜单 2=按钮
	Sort      int    `json:"sort" gorm:"default:0"`
	Visible   int    `json:"visible" gorm:"default:1"` // 0=隐藏 1=显示
	Status    int    `json:"status" gorm:"default:1"`  // 0=禁用 1=正常
	Perms     string `json:"perms" gorm:"size:128"`
}

// TableName 返回带前缀的表名
func (Menu) TableName() string {
	return TableName("menu")
}
