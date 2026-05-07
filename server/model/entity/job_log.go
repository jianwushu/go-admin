package entity

import "time"

// JobLog 定时任务执行日志实体
type JobLog struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	JobID     int64     `json:"jobId" gorm:"column:job_id;not null;index"`       // 任务ID
	JobName   string    `json:"jobName" gorm:"column:job_name;size:128"`         // 任务名称（冗余）
	Status    int       `json:"status" gorm:"column:status;default:1"`           // 1=成功 2=失败
	Result    string    `json:"result" gorm:"column:result;type:text"`           // 执行结果
	ErrorMsg  string    `json:"errorMsg" gorm:"column:error_msg;type:text"`      // 错误信息
	Duration  int64     `json:"duration" gorm:"column:duration"`                 // 执行耗时ms
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName 返回带前缀的表名
func (JobLog) TableName() string {
	return TableName("job_log")
}
