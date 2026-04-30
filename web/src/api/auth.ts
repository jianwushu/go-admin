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
