import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UserInfo } from '@/types/api'
import { getUserInfo } from '@/api/auth'
import { setToken, setRefreshToken, clearTokens } from '@/utils/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  const userInfo = ref<UserInfo | null>(null)
  const roles = ref<string[]>([])
  const permissions = ref<string[]>([])

  /** 设置 Token */
  function setTokenState(accessToken: string, refresh?: string) {
    token.value = accessToken
    setToken(accessToken)
    if (refresh) {
      refreshToken.value = refresh
      setRefreshToken(refresh)
    }
  }

  /** 清除 Token */
  function clearToken() {
    token.value = ''
    refreshToken.value = ''
    clearTokens()
  }

  /** 设置用户信息 */
  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    roles.value = info.roles || []
    permissions.value = info.permissions || []
  }

  /** 清除用户信息 */
  function clearUserInfo() {
    userInfo.value = null
    roles.value = []
    permissions.value = []
  }

  /** 获取用户信息 */
  async function fetchUserInfo() {
    const { data: res } = await getUserInfo()
    setUserInfo(res.data)
    return res.data
  }

  /** 登出 */
  function logout() {
    clearToken()
    clearUserInfo()
  }

  return {
    token,
    refreshToken,
    userInfo,
    roles,
    permissions,
    setToken: setTokenState,
    clearToken,
    setUserInfo,
    clearUserInfo,
    fetchUserInfo,
    logout,
  }
})
