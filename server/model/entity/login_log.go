package entity

import "time"

// LoginLog 登录日志实体
type LoginLog struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:64;index"`
	IP        string    `json:"ip" gorm:"size:64"`
	Location  string    `json:"location" gorm:"size:128"`
	Browser   string    `json:"browser" gorm:"size:64"`
	OS        string    `json:"os" gorm:"size:64"`
	Status    int       `json:"status"` // 0=失败 1=成功
	Msg       string    `json:"msg" gorm:"size:256"`
	CreatedAt time.Time `json:"createdAt" gorm:"index;autoCreateTime"`
}

// TableName 返回带前缀的表名
func (LoginLog) TableName() string {
	return TableName("login_log")
}
