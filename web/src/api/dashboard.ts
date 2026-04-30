import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

/** 仪表盘统计数据 */
export interface DashboardStats {
  userCount: number
  roleCount: number
  menuCount: number
  deptCount: number
  todayLogin: number
}

/** 获取仪表盘统计数据 */
export function getDashboardStats() {
  return request.get<ApiResponse<DashboardStats>>('/dashboard/stats')
}
