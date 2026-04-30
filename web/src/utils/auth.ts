/** Token 存取工具 */

const TOKEN_KEY = 'token'
const REFRESH_TOKEN_KEY = 'refreshToken'

/** 获取 Token */
export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

/** 设置 Token */
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

/** 移除 Token */
export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

/** 获取 Refresh Token */
export function getRefreshToken(): string {
  return localStorage.getItem(REFRESH_TOKEN_KEY) || ''
}

/** 设置 Refresh Token */
export function setRefreshToken(token: string): void {
  localStorage.setItem(REFRESH_TOKEN_KEY, token)
}

/** 移除 Refresh Token */
export function removeRefreshToken(): void {
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

/** 清除所有 Token */
export function clearTokens(): void {
  removeToken()
  removeRefreshToken()
}
