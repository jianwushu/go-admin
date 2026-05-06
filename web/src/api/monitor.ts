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

/** 在线用户信息 */
export interface OnlineUser {
  userId: number
  username: string
  nickname: string
  deptName: string
  ip: string
  loginTime: number
  onlineDuration: string
}

/** 获取在线用户列表 */
export function getOnlineUsers(username?: string) {
  return request.get<ApiResponse<OnlineUser[]>>('/monitor/online', { params: { username } })
}

/** 强制用户下线 */
export function forceLogoutUser(userId: number) {
  return request.delete<ApiResponse<null>>(`/monitor/online/${userId}`)
}

/** 获取服务器监控信息 */
export function getServerMonitor() {
  return request.get<ApiResponse<ServerMonitorInfo>>('/monitor/server')
}
