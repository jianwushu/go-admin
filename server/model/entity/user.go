package entity

// User 用户实体
type User struct {
	BaseModel
	Username string `json:"username" gorm:"size:64;not null;uniqueIndex"`
	Password string `json:"-" gorm:"size:128;not null"` // API 响应中不返回
	Nickname string `json:"nickname" gorm:"size:64"`
	Email    string `json:"email" gorm:"size:128"`
	Phone    string `json:"phone" gorm:"size:20"`
	Status   int    `json:"status" gorm:"default:1"` // 0=禁用 1=正常
	DeptID   int64  `json:"deptId" gorm:"index"`
	Avatar   string `json:"avatar" gorm:"size:256"`
	Remark   string `json:"remark" gorm:"size:512"`
}

// TableName 返回带前缀的表名
func (User) TableName() string {
	return TableName("user")
}
