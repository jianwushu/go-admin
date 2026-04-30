package entity

import "time"

// OperationLog 操作日志实体
type OperationLog struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Module       string    `json:"module" gorm:"size:64"`
	Action       string    `json:"action" gorm:"size:64"`
	Method       string    `json:"method" gorm:"size:10"`
	URL          string    `json:"url" gorm:"size:256"`
	IP           string    `json:"ip" gorm:"size:64"`
	Operator     string    `json:"operator" gorm:"size:64;index"`
	RequestParam string    `json:"requestParam" gorm:"type:text"`
	ResponseData string    `json:"responseData" gorm:"type:text"`
	Status       int       `json:"status"` // 0=失败 1=成功
	ErrorMsg     string    `json:"errorMsg" gorm:"size:512"`
	Duration     int64     `json:"duration"` // 耗时（毫秒）
	CreatedAt    time.Time `json:"createdAt" gorm:"index;autoCreateTime"`
}

// TableName 返回带前缀的表名
func (OperationLog) TableName() string {
	return TableName("operation_log")
}
