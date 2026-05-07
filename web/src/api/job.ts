import request from '@/utils/request'
import type { ApiResponse, PageResponse, JobItem, JobFormData, JobListParams, JobLogItem, JobLogListParams } from '@/types/api'

// ==================== 定时任务 API ====================

/** 获取任务列表 */
export function getJobList(params: JobListParams) {
  return request.get<PageResponse<JobItem>>('/tool/job/list', { params })
}

/** 获取任务详情 */
export function getJobById(id: number) {
  return request.get<ApiResponse<JobItem>>(`/tool/job/${id}`)
}

/** 创建任务 */
export function createJob(data: JobFormData) {
  return request.post<ApiResponse<null>>('/tool/job', data)
}

/** 更新任务 */
export function updateJob(data: JobFormData) {
  return request.put<ApiResponse<null>>('/tool/job', data)
}

/** 删除任务（支持单个或多个ID） */
export function deleteJob(id: number | number[]) {
  const idStr = Array.isArray(id) ? id.join(',') : id
  return request.delete<ApiResponse<null>>(`/tool/job/${idStr}`)
}

/** 修改任务状态 */
export function changeJobStatus(id: number, status: number) {
  return request.put<ApiResponse<null>>('/tool/job/change-status', { id, status })
}

/** 手动执行一次任务 */
export function runJobOnce(id: number) {
  return request.post<ApiResponse<null>>(`/tool/job/run-once/${id}`)
}

// ==================== 任务日志 API ====================

/** 获取任务日志列表 */
export function getJobLogList(params: JobLogListParams) {
  return request.get<PageResponse<JobLogItem>>('/tool/job/log/list', { params })
}

/** 清理指定任务的日志 */
export function cleanJobLogs(jobId: number) {
  return request.delete<ApiResponse<null>>('/tool/job/log/clean', { params: { jobId } })
}
