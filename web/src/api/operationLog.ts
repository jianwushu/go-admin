import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types/api'

/** 操作日志项 */
export interface OperationLogItem {
  id: number
  module: string
  action: string
  method: string
  url: string
  ip: string
  operator: string
  requestParam: string
  responseData: string
  status: number
  errorMsg: string
  duration: number
  createdAt: string
}

/** 操作日志查询参数 */
export interface OperationLogListParams {
  page: number
  pageSize: number
  module?: string
  operator?: string
  status?: number
  method?: string
}

/** 获取操作日志列表 */
export function getOperationLogList(params: OperationLogListParams) {
  return request.get<PageResponse<OperationLogItem>>('/monitor/operation-log/list', { params })
}

/** 清空操作日志 */
export function clearOperationLog() {
  return request.delete<ApiResponse<null>>('/monitor/operation-log/clear')
}
