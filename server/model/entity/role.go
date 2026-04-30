package entity

// Role 角色实体
type Role struct {
	BaseModel
	Name      string `json:"name" gorm:"size:64;not null"`
	Code      string `json:"code" gorm:"size:64;not null;uniqueIndex"`
	DataScope int    `json:"dataScope" gorm:"default:1"` // 1=全部 2=本部门 3=本部门及下级 4=仅本人 5=自定义
	Sort      int    `json:"sort" gorm:"default:0"`
	Status    int    `json:"status" gorm:"default:1"` // 0=禁用 1=正常
	Remark    string `json:"remark" gorm:"size:512"`
}

// TableName 返回带前缀的表名
func (Role) TableName() string {
	return TableName("role")
}
