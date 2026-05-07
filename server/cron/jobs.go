package cron

import (
	"fmt"
	"time"

	"go-admin/global"

	"go.uber.org/zap"
)

// BuiltinJobs 预定义的内置任务函数
// 返回值: (执行结果描述, 错误)
var BuiltinJobs = map[string]func() (string, error){
	"clean_operation_log": CleanOperationLog,
	"clean_login_log":     CleanLoginLog,
}

// CleanOperationLog 清理操作日志（保留最近30天）
func CleanOperationLog() (string, error) {
	days := 30
	result := global.DB.Exec(
		fmt.Sprintf("DELETE FROM %soperation_log WHERE created_at < datetime('now', '-%d days')", global.Config.TablePrefix, days),
	)
	if result.Error != nil {
		return "", fmt.Errorf("清理操作日志失败: %v", result.Error)
	}
	msg := fmt.Sprintf("清理操作日志完成，删除 %d 条记录", result.RowsAffected)
	global.Log.Info(msg, zap.Int64("rows", result.RowsAffected))
	return msg, nil
}

// CleanLoginLog 清理登录日志（保留最近30天）
func CleanLoginLog() (string, error) {
	days := 30
	result := global.DB.Exec(
		fmt.Sprintf("DELETE FROM %slogin_log WHERE created_at < datetime('now', '-%d days')", global.Config.TablePrefix, days),
	)
	if result.Error != nil {
		return "", fmt.Errorf("清理登录日志失败: %v", result.Error)
	}
	msg := fmt.Sprintf("清理登录日志完成，删除 %d 条记录", result.RowsAffected)
	global.Log.Info(msg, zap.Int64("rows", result.RowsAffected))
	return msg, nil
}

// InitJobs 从数据库加载所有启用的任务并启动
func InitJobs() {
	if manager == nil {
		global.Log.Warn("定时任务管理器未初始化，跳过加载任务")
		return
	}

	// 等待数据库就绪
	time.Sleep(500 * time.Millisecond)

	var jobs []struct {
		ID         int64  `gorm:"column:id"`
		Name       string `gorm:"column:name"`
		JobType    int    `gorm:"column:job_type"`
		CronExpr   string `gorm:"column:cron_expr"`
		FuncName   string `gorm:"column:func_name"`
		HTTPUrl    string `gorm:"column:http_url"`
		HTTPMethod string `gorm:"column:http_method"`
		HTTPBody   string `gorm:"column:http_body"`
		Status     int    `gorm:"column:status"`
	}

	tableName := fmt.Sprintf("%sjob", global.Config.TablePrefix)
	if err := global.DB.Table(tableName).Where("status = ?", 1).Find(&jobs).Error; err != nil {
		global.Log.Error("加载定时任务失败", zap.Error(err))
		return
	}

	for _, j := range jobs {
		job := &JobEntity{
			ID:         j.ID,
			Name:       j.Name,
			JobType:    j.JobType,
			CronExpr:   j.CronExpr,
			FuncName:   j.FuncName,
			HTTPUrl:    j.HTTPUrl,
			HTTPMethod: j.HTTPMethod,
			HTTPBody:   j.HTTPBody,
			Status:     j.Status,
		}
		if err := manager.AddJob(job); err != nil {
			global.Log.Error("添加定时任务失败",
				zap.Int64("jobID", j.ID),
				zap.String("name", j.Name),
				zap.Error(err),
			)
		}
	}

	global.Log.Info("定时任务加载完成", zap.Int("count", len(jobs)))
}

// JobEntity 用于 cron 包内部的任务实体（避免循环依赖）
type JobEntity struct {
	ID         int64
	Name       string
	JobType    int
	CronExpr   string
	FuncName   string
	HTTPUrl    string
	HTTPMethod string
	HTTPBody   string
	Status     int
}
