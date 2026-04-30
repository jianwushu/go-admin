package entity

// Dept 部门实体
type Dept struct {
	BaseModel
	ParentID int64  `json:"parentId" gorm:"default:0"`
	Name     string `json:"name" gorm:"size:64;not null"`
	Sort     int    `json:"sort" gorm:"default:0"`
	Status   int    `json:"status" gorm:"default:1"` // 0=禁用 1=正常
	Leader   string `json:"leader" gorm:"size:64"`
	Phone    string `json:"phone" gorm:"size:20"`
	Email    string `json:"email" gorm:"size:128"`
}

// TableName 返回带前缀的表名
func (Dept) TableName() string {
	return TableName("dept")
}
