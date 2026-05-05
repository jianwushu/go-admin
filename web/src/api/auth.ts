import request from '@/utils/request'
import type { ApiResponse, LoginRequest, LoginResult, UserInfo, MenuItem } from '@/types/api'

/** 登录 */
export function login(data: LoginRequest) {
  return request.post<ApiResponse<LoginResult>>('/auth/login', data)
}

/** 登出 */
export function logout() {
  return request.post<ApiResponse<null>>('/auth/logout')
}

/** 刷新 Token */
export function refreshToken(refreshToken: string) {
  return request.post<ApiResponse<{ accessToken: string }>>('/auth/refresh', { refreshToken })
}

/** 获取当前用户信息 */
export function getUserInfo() {
  return request.get<ApiResponse<UserInfo>>('/user/info')
}

/** 获取当前用户菜单 */
export function getUserMenus() {
  return request.get<ApiResponse<MenuItem[]>>('/user/menus')
}

/** 个人资料响应 */
export interface UserProfile {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  deptId: number
  deptName: string
  remark: string
  createdAt: string
}

/** 个人资料更新请求 */
export interface UserProfileUpdateData {
  nickname: string
  email: string
  phone: string
  avatar: string
}

/** 修改密码请求 */
export interface ChangePasswordData {
  oldPassword: string
  newPassword: string
}

/** 获取个人资料 */
export function getUserProfile() {
  return request.get<ApiResponse<UserProfile>>('/user/profile')
}

/** 更新个人资料 */
export function updateUserProfile(data: UserProfileUpdateData) {
  return request.put<ApiResponse<null>>('/user/profile', data)
}

/** 修改密码 */
export function changePassword(data: ChangePasswordData) {
  return request.put<ApiResponse<null>>('/user/change-password', data)
}
