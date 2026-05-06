import request from '@/utils/request'
import type { ApiResponse, PageResponse, UserItem, UserFormData, UserListParams, RoleInfo, RoleFormData, RoleListParams, MenuItem, MenuFormData, DeptInfo, DeptFormData, SystemConfigItem, SystemConfigListParams, SystemConfigBatchUpdateItem } from '@/types/api'

// ==================== 用户管理 API ====================

/** 获取用户列表 */
export function getUserList(params: UserListParams) {
  return request.get<PageResponse<UserItem>>('/system/user/list', { params })
}

/** 获取用户详情 */
export function getUserById(id: number) {
  return request.get<ApiResponse<UserItem>>(`/system/user/${id}`)
}

/** 创建用户 */
export function createUser(data: UserFormData) {
  return request.post<ApiResponse<null>>('/system/user', data)
}

/** 更新用户 */
export function updateUser(data: UserFormData) {
  return request.put<ApiResponse<null>>('/system/user', data)
}

/** 删除用户（支持单个或多个ID） */
export function deleteUser(id: number | number[]) {
  const idStr = Array.isArray(id) ? id.join(',') : id
  return request.delete<ApiResponse<null>>(`/system/user/${idStr}`)
}

/** 重置用户密码 */
export function resetUserPassword(id: number, password: string) {
  return request.put<ApiResponse<null>>('/system/user/reset-password', { id, password })
}

/** 修改用户状态 */
export function changeUserStatus(id: number, status: number) {
  return request.put<ApiResponse<null>>('/system/user/change-status', { id, status })
}

// ==================== 角色管理 API ====================

/** 获取角色列表 */
export function getRoleList(params: RoleListParams) {
  return request.get<PageResponse<RoleInfo>>('/system/role/list', { params })
}

/** 获取所有角色（不分页） */
export function getAllRoles() {
  return request.get<ApiResponse<RoleInfo[]>>('/system/role/all')
}

/** 获取角色详情 */
export function getRoleById(id: number) {
  return request.get<ApiResponse<RoleInfo>>(`/system/role/${id}`)
}

/** 创建角色 */
export function createRole(data: RoleFormData) {
  return request.post<ApiResponse<null>>('/system/role', data)
}

/** 更新角色 */
export function updateRole(data: RoleFormData) {
  return request.put<ApiResponse<null>>('/system/role', data)
}

/** 删除角色（支持单个或多个ID） */
export function deleteRole(id: number | number[]) {
  const idStr = Array.isArray(id) ? id.join(',') : id
  return request.delete<ApiResponse<null>>(`/system/role/${idStr}`)
}

/** 修改角色状态 */
export function changeRoleStatus(id: number, status: number) {
  return request.put<ApiResponse<null>>('/system/role/change-status', { id, status })
}

// ==================== 菜单管理 API ====================

/** 获取菜单树形列表 */
export function getMenuTree() {
  return request.get<ApiResponse<MenuItem[]>>('/system/menu/tree')
}

/** 获取菜单详情 */
export function getMenuById(id: number) {
  return request.get<ApiResponse<MenuItem>>(`/system/menu/${id}`)
}

/** 创建菜单 */
export function createMenu(data: MenuFormData) {
  return request.post<ApiResponse<null>>('/system/menu', data)
}

/** 更新菜单 */
export function updateMenu(data: MenuFormData) {
  return request.put<ApiResponse<null>>('/system/menu', data)
}

/** 删除菜单 */
export function deleteMenu(id: number) {
  return request.delete<ApiResponse<null>>(`/system/menu/${id}`)
}

// ==================== 部门管理 API ====================

/** 获取部门树形列表 */
export function getDeptTree() {
  return request.get<ApiResponse<DeptInfo[]>>('/system/dept/tree')
}

/** 获取部门详情 */
export function getDeptById(id: number) {
  return request.get<ApiResponse<DeptInfo>>(`/system/dept/${id}`)
}

/** 创建部门 */
export function createDept(data: DeptFormData) {
  return request.post<ApiResponse<null>>('/system/dept', data)
}

/** 更新部门 */
export function updateDept(data: DeptFormData) {
  return request.put<ApiResponse<null>>('/system/dept', data)
}

/** 删除部门 */
export function deleteDept(id: number) {
  return request.delete<ApiResponse<null>>(`/system/dept/${id}`)
}

// ==================== 系统配置 API ====================

/** 获取所有系统配置 */
export function getAllSystemConfig() {
  return request.get<ApiResponse<SystemConfigItem[]>>('/system/config/all')
}

/** 获取系统配置列表（分页） */
export function getSystemConfigList(params: SystemConfigListParams) {
  return request.get<PageResponse<SystemConfigItem>>('/system/config/list', { params })
}

/** 根据键获取配置 */
export function getSystemConfigByKey(key: string) {
  return request.get<ApiResponse<SystemConfigItem>>(`/system/config/key/${key}`)
}

/** 批量获取配置 */
export function getSystemConfigByKeys(keys: string[]) {
  return request.get<ApiResponse<Record<string, string>>>('/system/config/keys', { params: { keys: keys.join(',') } })
}

/** 更新单个配置 */
export function updateSystemConfig(id: number, configValue: string) {
  return request.put<ApiResponse<null>>('/system/config', { id, configValue })
}

/** 批量更新配置 */
export function batchUpdateSystemConfig(list: SystemConfigBatchUpdateItem[]) {
  return request.put<ApiResponse<null>>('/system/config/batch', { list })
}

/** 上传Logo图片 */
export function uploadSystemLogo(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<ApiResponse<string>>('/system/config/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
