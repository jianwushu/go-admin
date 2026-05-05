import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types/api'

/** 登录日志项 */
export interface LoginLogItem {
  id: number
  username: string
  ip: string
  location: string
  browser: string
  os: string
  status: number
  msg: string
  createdAt: string
}

/** 登录日志查询参数 */
export interface LoginLogListParams {
  page: number
  pageSize: number
  username?: string
  ip?: string
  status?: number
  startTime?: string
  endTime?: string
}

/** 获取登录日志列表 */
export function getLoginLogList(params: LoginLogListParams) {
  return request.get<PageResponse<LoginLogItem>>('/monitor/login-log/list', { params })
}

/** 清空登录日志 */
export function clearLoginLog() {
  return request.delete<ApiResponse<null>>('/monitor/login-log/clear')
}
