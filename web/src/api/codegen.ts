import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// ==================== 代码生成 类型定义 ====================

/** 表信息 */
export interface TableInfo {
  tableName: string
  tableComment: string
  engine: string
  createTime: string
  updateTime: string
}

/** 列信息 */
export interface ColumnInfo {
  columnName: string
  columnType: string
  columnComment: string
  isPk: boolean
  isNull: boolean
  maxLength: number
  columnDefault: string
  extra: string
}

/** 字段配置 */
export interface ColumnConfig {
  columnName: string
  columnType: string
  goType: string
  goField: string
  tsType: string
  label: string
  htmlType: string
  queryType: string
  isList: boolean
  isQuery: boolean
  isRequired: boolean
  isEdit: boolean
  dictType: string
  comment: string
  isPk: boolean
  isNull: boolean
  maxLength: number
  sort: number
}

/** 代码生成配置 */
export interface CodegenConfig {
  id?: number
  tableName: string
  tableComment: string
  className: string
  businessName: string
  functionName: string
  moduleName: string
  packageName: string
  author: string
  fields: ColumnConfig[]
  createdAt?: string
}

/** 代码预览请求 */
export interface CodegenPreviewRequest {
  tableName: string
  tableComment: string
  className: string
  businessName: string
  functionName: string
  moduleName: string
  packageName: string
  author: string
  fields: ColumnConfig[]
}

/** 代码文件 */
export interface CodegenFile {
  fileName: string
  filePath: string
  content: string
}

/** 代码预览响应 */
export interface CodegenPreviewResponse {
  files: CodegenFile[]
}

// ==================== 代码生成 API ====================

/** 获取数据库表列表 */
export function getTableList() {
  return request.get<ApiResponse<TableInfo[]>>('/codegen/tables')
}

/** 获取表的列信息 */
export function getColumnList(tableName: string) {
  return request.get<ApiResponse<ColumnInfo[]>>(`/codegen/columns/${tableName}`)
}

/** 代码预览 */
export function previewCode(data: CodegenPreviewRequest) {
  return request.post<ApiResponse<CodegenPreviewResponse>>('/codegen/preview', data)
}

/** 生成代码并下载 */
export function generateCode(data: CodegenPreviewRequest) {
  return request.post('/codegen/generate', data, { responseType: 'blob' })
}

/** 保存代码生成配置 */
export function saveConfig(data: CodegenConfig) {
  return request.post<ApiResponse<null>>('/codegen/config', data)
}

/** 获取代码生成配置 */
export function getConfig(tableName: string) {
  return request.get<ApiResponse<CodegenConfig>>(`/codegen/config/${tableName}`)
}

/** 获取所有代码生成配置 */
export function getAllConfigs() {
  return request.get<ApiResponse<CodegenConfig[]>>('/codegen/configs')
}

/** 删除代码生成配置 */
export function deleteConfig(id: number) {
  return request.delete<ApiResponse<null>>(`/codegen/config/${id}`)
}
