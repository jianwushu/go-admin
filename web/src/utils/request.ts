import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, setToken, getRefreshToken, clearTokens } from '@/utils/auth'
import type { ApiResponse } from '@/types/api'

/** 是否正在刷新 Token */
let isRefreshing = false
/** 等待刷新 Token 的请求队列 */
let pendingRequests: Array<(token: string) => void> = []

/** 创建 Axios 实例 */
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

/** 请求拦截器 */
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

/** 响应拦截器 */
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data

    // 如果是文件下载，直接返回
    if (response.config.responseType === 'blob') {
      return response
    }

    // 业务错误（后端 code=0 表示成功，非0表示失败）
    if (res.code !== 0) {
      ElMessage.error(res.msg || '请求失败')

      // Token 过期
      if (res.code === 401) {
        handleTokenExpired(response.config)
      }

      return Promise.reject(new Error(res.msg || '请求失败'))
    }

    return response
  },
  (error) => {
    const { response } = error

    if (response) {
      switch (response.status) {
        case 401:
          handleTokenExpired(error.config)
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(response.data?.msg || '请求失败')
      }
    } else {
      ElMessage.error('网络连接异常')
    }

    return Promise.reject(error)
  }
)

/** 处理 Token 过期 */
function handleTokenExpired(config: AxiosRequestConfig) {
  if (!isRefreshing) {
    isRefreshing = true
    const refreshToken = getRefreshToken()

    if (!refreshToken) {
      handleLogout()
      return
    }

    // 调用刷新 Token 接口
    axios
      .post<ApiResponse<{ accessToken: string }>>(
        `${import.meta.env.VITE_APP_BASE_API}/auth/refresh`,
        { refreshToken }
      )
      .then((res) => {
        if (res.data.code === 0) {
          const newToken = res.data.data.accessToken
          setToken(newToken)

          // 重试队列中的请求
          pendingRequests.forEach((cb) => cb(newToken))
          pendingRequests = []

          // 重试当前请求
          if (config.headers) {
            config.headers.Authorization = `Bearer ${newToken}`
          }
          return service(config)
        } else {
          handleLogout()
        }
      })
      .catch(() => {
        handleLogout()
      })
      .finally(() => {
        isRefreshing = false
      })
  } else {
    // 将请求加入队列，等待 Token 刷新
    return new Promise((resolve) => {
      pendingRequests.push((token: string) => {
        if (config.headers) {
          config.headers.Authorization = `Bearer ${token}`
        }
        resolve(service(config))
      })
    })
  }
}

/** 处理登出 */
function handleLogout() {
  clearTokens()
  window.location.href = '/login'
}

export default service
