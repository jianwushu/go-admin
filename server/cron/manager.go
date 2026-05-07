package cron

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go-admin/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Manager 定时任务管理器
type Manager struct {
	cron    *cron.Cron
	entries map[int64]cron.EntryID // jobID -> cron entryID
	mu      sync.RWMutex
}

// 全局定时任务管理器实例
var manager *Manager

// GetManager 获取全局定时任务管理器实例
func GetManager() *Manager {
	return manager
}

// Init 初始化定时任务管理器
func Init() {
	manager = &Manager{
		cron:    cron.New(cron.WithSeconds()),
		entries: make(map[int64]cron.EntryID),
	}
	manager.cron.Start()
	global.Log.Info("定时任务管理器已启动")
}

// Stop 停止定时任务管理器
func (m *Manager) Stop() {
	ctx := m.cron.Stop()
	<-ctx.Done()
	global.Log.Info("定时任务管理器已停止")
}

// AddJob 添加定时任务
func (m *Manager) AddJob(job *JobEntity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已存在，先移除
	if entryID, exists := m.entries[job.ID]; exists {
		m.cron.Remove(entryID)
		delete(m.entries, job.ID)
	}

	// 禁用状态不添加
	if job.Status != 1 {
		return nil
	}

	// 创建执行函数
	handler := m.CreateHandler(job)

	// 添加到 cron
	entryID, err := m.cron.AddFunc(job.CronExpr, handler)
	if err != nil {
		return fmt.Errorf("添加定时任务失败 [%s]: %v", job.CronExpr, err)
	}

	m.entries[job.ID] = entryID
	global.Log.Info("定时任务已添加", zap.Int64("jobID", job.ID), zap.String("name", job.Name), zap.String("cron", job.CronExpr))
	return nil
}

// RemoveJob 移除定时任务
func (m *Manager) RemoveJob(jobID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entryID, exists := m.entries[jobID]; exists {
		m.cron.Remove(entryID)
		delete(m.entries, jobID)
		global.Log.Info("定时任务已移除", zap.Int64("jobID", jobID))
	}
}

// UpdateJob 更新定时任务（先移除再添加）
func (m *Manager) UpdateJob(job *JobEntity) error {
	m.RemoveJob(job.ID)
	return m.AddJob(job)
}

// IsRunning 检查任务是否在运行
func (m *Manager) IsRunning(jobID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.entries[jobID]
	return exists
}

// CreateHandler 创建任务执行函数（导出供外部调用）
func (m *Manager) CreateHandler(job *JobEntity) func() {
	return func() {
		startTime := time.Now()
		logEntry := &JobLogEntry{
			JobID:   job.ID,
			JobName: job.Name,
		}

		var err error
		var result string

		switch job.JobType {
		case 1: // 内置函数
			result, err = m.executeBuiltinJob(job.FuncName)
		case 2: // HTTP请求
			result, err = m.executeHTTPJob(job)
		default:
			err = fmt.Errorf("未知的任务类型: %d", job.JobType)
		}

		duration := time.Since(startTime).Milliseconds()
		logEntry.Duration = duration

		if err != nil {
			logEntry.Status = 2 // 失败
			logEntry.ErrorMsg = err.Error()
			global.Log.Error("定时任务执行失败",
				zap.Int64("jobID", job.ID),
				zap.String("name", job.Name),
				zap.Error(err),
				zap.Int64("duration", duration),
			)
		} else {
			logEntry.Status = 1 // 成功
			logEntry.Result = result
			global.Log.Info("定时任务执行成功",
				zap.Int64("jobID", job.ID),
				zap.String("name", job.Name),
				zap.Int64("duration", duration),
			)
		}

		// 记录执行日志到数据库
		m.saveLog(logEntry)
	}
}

// saveLog 保存执行日志
func (m *Manager) saveLog(logEntry *JobLogEntry) {
	tableName := fmt.Sprintf("%sjob_log", global.Config.TablePrefix)
	if err := global.DB.Table(tableName).Create(logEntry).Error; err != nil {
		global.Log.Error("保存任务执行日志失败", zap.Error(err))
	}
}

// executeBuiltinJob 执行内置函数
func (m *Manager) executeBuiltinJob(funcName string) (string, error) {
	fn, exists := BuiltinJobs[funcName]
	if !exists {
		return "", fmt.Errorf("未找到内置函数: %s", funcName)
	}
	return fn()
}

// executeHTTPJob 执行HTTP请求
func (m *Manager) executeHTTPJob(job *JobEntity) (string, error) {
	method := strings.ToUpper(job.HTTPMethod)
	if method == "" {
		method = "GET"
	}

	var body io.Reader
	if job.HTTPBody != "" {
		body = bytes.NewBufferString(job.HTTPBody)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, job.HTTPUrl, body)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("执行HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取HTTP响应失败: %v", err)
	}

	result := fmt.Sprintf("HTTP %s %s -> %d", method, job.HTTPUrl, resp.StatusCode)
	if len(respBody) > 0 && len(respBody) < 2000 {
		result += ": " + string(respBody)
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP请求返回错误状态码: %d", resp.StatusCode)
	}

	return result, nil
}

// JobLogEntry 用于 cron 包内部的日志实体（避免循环依赖）
type JobLogEntry struct {
	JobID    int64  `gorm:"column:job_id"`
	JobName  string `gorm:"column:job_name"`
	Status   int    `gorm:"column:status"`
	Result   string `gorm:"column:result"`
	ErrorMsg string `gorm:"column:error_msg"`
	Duration int64  `gorm:"column:duration"`
}
