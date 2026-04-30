import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

/** CPU 信息 */
export interface CpuInfo {
  cores: number
  usage: number[]
  usedRate: number
}

/** 内存信息 */
export interface MemoryInfo {
  total: number
  used: number
  free: number
  usedRate: number
}

/** 磁盘信息 */
export interface DiskInfo {
  mountPoint: string
  total: number
  used: number
  free: number
  usedRate: number
}

/** Go 运行时信息 */
export interface GoRuntimeInfo {
  goVersion: string
  os: string
  arch: string
  goroutines: number
  heapAlloc: number
  heapSys: number
  heapIdle: number
  heapInuse: number
  numGC: number
  lastGC: string
}

/** 服务状态 */
export interface ServiceStatus {
  status: 'online' | 'offline'
  message: string
}

/** 服务器监控信息 */
export interface ServerMonitorInfo {
  cpu: CpuInfo
  memory: MemoryInfo
  disk: DiskInfo[]
  goRuntime: GoRuntimeInfo
  db: ServiceStatus
  redis: ServiceStatus
}

/** 获取服务器监控信息 */
export function getServerMonitor() {
  return request.get<ApiResponse<ServerMonitorInfo>>('/monitor/server')
}
