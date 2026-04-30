package service

import (
	"context"
	"fmt"
	"go-admin/global"
	"go-admin/model/response"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

// MonitorService 服务器监控服务
type MonitorService struct{}

// NewMonitorService 创建监控服务实例
func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

// GetServerInfo 获取服务器监控信息
func (s *MonitorService) GetServerInfo() (*response.ServerMonitorResponse, error) {
	result := &response.ServerMonitorResponse{}

	// 获取 CPU 信息
	cpuInfo, err := s.getCPUInfo()
	if err != nil {
		global.Log.Warn("获取CPU信息失败", zap.Error(err))
	} else {
		result.CPU = *cpuInfo
	}

	// 获取内存信息
	memInfo, err := s.getMemoryInfo()
	if err != nil {
		global.Log.Warn("获取内存信息失败", zap.Error(err))
	} else {
		result.Memory = *memInfo
	}

	// 获取磁盘信息
	diskInfo, err := s.getDiskInfo()
	if err != nil {
		global.Log.Warn("获取磁盘信息失败", zap.Error(err))
	} else {
		result.Disk = diskInfo
	}

	// 获取 Go 运行时信息
	result.GoRuntime = s.getGoRuntimeInfo()

	// 获取数据库状态
	result.DB = s.getDBStatus()

	// 获取 Redis 状态
	result.Redis = s.getRedisStatus()

	return result, nil
}

// getCPUInfo 获取 CPU 信息
func (s *MonitorService) getCPUInfo() (*response.CPUInfo, error) {
	// 获取 CPU 核心数
	cores, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("获取CPU核心数失败: %w", err)
	}

	// 获取各核心使用率
	usage, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("获取CPU使用率失败: %w", err)
	}

	// 计算总体使用率
	totalUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("获取CPU总体使用率失败: %w", err)
	}

	var usedRate float64
	if len(totalUsage) > 0 {
		usedRate = totalUsage[0]
	}

	return &response.CPUInfo{
		Cores:    cores,
		Usage:    usage,
		UsedRate: usedRate,
	}, nil
}

// getMemoryInfo 获取内存信息
func (s *MonitorService) getMemoryInfo() (*response.MemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %w", err)
	}

	return &response.MemoryInfo{
		Total:    vmStat.Total / 1024 / 1024,  // 转换为 MB
		Used:     vmStat.Used / 1024 / 1024,   // 转换为 MB
		Free:     vmStat.Free / 1024 / 1024,    // 转换为 MB
		UsedRate: vmStat.UsedPercent,
	}, nil
}

// getDiskInfo 获取磁盘信息
func (s *MonitorService) getDiskInfo() ([]response.DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘分区失败: %w", err)
	}

	var diskInfos []response.DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		// 跳过特殊文件系统
		if usage.Total == 0 {
			continue
		}

		diskInfos = append(diskInfos, response.DiskInfo{
			MountPoint: p.Mountpoint,
			Total:      usage.Total / 1024 / 1024 / 1024,  // 转换为 GB
			Used:       usage.Used / 1024 / 1024 / 1024,   // 转换为 GB
			Free:       usage.Free / 1024 / 1024 / 1024,    // 转换为 GB
			UsedRate:   usage.UsedPercent,
		})
	}

	return diskInfos, nil
}

// getGoRuntimeInfo 获取 Go 运行时信息
func (s *MonitorService) getGoRuntimeInfo() response.GoRuntimeInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return response.GoRuntimeInfo{
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Goroutines: runtime.NumGoroutine(),
		HeapAlloc:  m.HeapAlloc / 1024 / 1024,   // 转换为 MB
		HeapSys:    m.HeapSys / 1024 / 1024,     // 转换为 MB
		HeapIdle:   m.HeapIdle / 1024 / 1024,    // 转换为 MB
		HeapInuse:  m.HeapInuse / 1024 / 1024,   // 转换为 MB
		NumGC:      m.NumGC,
		LastGC:     time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"),
	}
}

// getDBStatus 获取数据库连接状态
func (s *MonitorService) getDBStatus() response.ServiceStatus {
	if global.DB == nil {
		return response.ServiceStatus{
			Status:  "offline",
			Message: "数据库未初始化",
		}
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		return response.ServiceStatus{
			Status:  "offline",
			Message: fmt.Sprintf("获取数据库实例失败: %s", err.Error()),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return response.ServiceStatus{
			Status:  "offline",
			Message: fmt.Sprintf("数据库连接失败: %s", err.Error()),
		}
	}

	stats := sqlDB.Stats()
	return response.ServiceStatus{
		Status:  "online",
		Message: fmt.Sprintf("连接池: %d/%d", stats.InUse, stats.MaxOpenConnections),
	}
}

// getRedisStatus 获取 Redis 连接状态
func (s *MonitorService) getRedisStatus() response.ServiceStatus {
	if global.Redis == nil {
		return response.ServiceStatus{
			Status:  "offline",
			Message: "Redis 未初始化",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := global.Redis.Ping(ctx).Err(); err != nil {
		return response.ServiceStatus{
			Status:  "offline",
			Message: fmt.Sprintf("Redis 连接失败: %s", err.Error()),
		}
	}

	return response.ServiceStatus{
		Status:  "online",
		Message: "连接正常",
	}
}
