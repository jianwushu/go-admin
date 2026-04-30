package response

// ServerMonitorResponse 服务器监控响应
type ServerMonitorResponse struct {
	CPU       CPUInfo       `json:"cpu"`       // CPU 信息
	Memory    MemoryInfo    `json:"memory"`    // 内存信息
	Disk      []DiskInfo    `json:"disk"`      // 磁盘信息
	GoRuntime GoRuntimeInfo `json:"goRuntime"` // Go 运行时信息
	DB        ServiceStatus `json:"db"`        // 数据库状态
	Redis     ServiceStatus `json:"redis"`     // Redis 状态
}

// CPUInfo CPU 信息
type CPUInfo struct {
	Cores    int       `json:"cores"`    // CPU 核心数
	Usage    []float64 `json:"usage"`    // 各核心使用率
	UsedRate float64   `json:"usedRate"` // 总体使用率
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total     uint64  `json:"total"`     // 总内存 (MB)
	Used      uint64  `json:"used"`      // 已用内存 (MB)
	Free      uint64  `json:"free"`      // 空闲内存 (MB)
	UsedRate  float64 `json:"usedRate"`  // 使用率 (%)
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	MountPoint string  `json:"mountPoint"` // 挂载点
	Total      uint64  `json:"total"`      // 总容量 (GB)
	Used       uint64  `json:"used"`       // 已用 (GB)
	Free       uint64  `json:"free"`       // 可用 (GB)
	UsedRate   float64 `json:"usedRate"`   // 使用率 (%)
}

// GoRuntimeInfo Go 运行时信息
type GoRuntimeInfo struct {
	GoVersion    string `json:"goVersion"`    // Go 版本
	OS           string `json:"os"`           // 操作系统
	Arch         string `json:"arch"`         // 架构
	Goroutines   int    `json:"goroutines"`   // Goroutine 数量
	HeapAlloc    uint64 `json:"heapAlloc"`    // 堆内存分配 (MB)
	HeapSys      uint64 `json:"heapSys"`      // 堆系统内存 (MB)
	HeapIdle     uint64 `json:"heapIdle"`     // 堆空闲内存 (MB)
	HeapInuse    uint64 `json:"heapInuse"`    // 堆使用中内存 (MB)
	NumGC        uint32 `json:"numGC"`        // GC 次数
	LastGC       string `json:"lastGC"`       // 上次 GC 时间
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Status  string `json:"status"`  // 状态: online / offline
	Message string `json:"message"` // 附加信息
}
