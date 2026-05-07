package entity

import (
	"time"

	"gorm.io/gorm"
)

// Job 定时任务实体
type Job struct {
	ID         int64          `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"column:name;size:128;not null"`          // 任务名称
	JobType    int            `json:"jobType" gorm:"column:job_type;default:1"`            // 1=内置函数 2=HTTP请求
	CronExpr   string         `json:"cronExpr" gorm:"column:cron_expr;size:64;not null"`   // cron表达式
	FuncName   string         `json:"funcName" gorm:"column:func_name;size:128"`           // 内置函数名
	HTTPUrl    string         `json:"httpUrl" gorm:"column:http_url;size:512"`             // HTTP地址
	HTTPMethod string         `json:"httpMethod" gorm:"column:http_method;size:10;default:GET"` // GET/POST
	HTTPBody   string         `json:"httpBody" gorm:"column:http_body;type:text"`          // 请求体
	Status     int            `json:"status" gorm:"column:status;default:1"`               // 1=启用 0=禁用
	Remark     string         `json:"remark" gorm:"column:remark;size:512"`
	CreatedBy  int64          `json:"createdBy" gorm:"column:created_by;default:0"`
	UpdatedBy  int64          `json:"updatedBy" gorm:"column:updated_by;default:0"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 返回带前缀的表名
func (Job) TableName() string {
	return TableName("job")
}
