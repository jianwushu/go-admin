import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPublicSystemConfig } from '@/api/system'

/** 系统配置 Store - 管理全局系统配置（标题、版权等） */
export const useConfigStore = defineStore('config', () => {
  /** 系统标题 */
  const systemTitle = ref<string>(localStorage.getItem('system_title') || 'Go-Admin')

  /** 版权信息 */
  const systemCopyright = ref<string>(localStorage.getItem('system_copyright') || '')

  /** 系统 Logo */
  const systemLogo = ref<string>(localStorage.getItem('system_logo') || '')

  /** 是否已加载 */
  const loaded = ref(false)

  /** 加载系统配置（使用公开接口，无需认证） */
  async function loadConfig() {
    if (loaded.value) return
    try {
      const { data: res } = await getPublicSystemConfig(['system_title', 'system_copyright', 'system_logo'])
      if (res.code === 0 && res.data) {
        if (res.data.system_title) {
          systemTitle.value = res.data.system_title
          localStorage.setItem('system_title', res.data.system_title)
          // 同步更新页面标题
          document.title = res.data.system_title
        }
        if (res.data.system_copyright) {
          systemCopyright.value = res.data.system_copyright
          localStorage.setItem('system_copyright', res.data.system_copyright)
        }
        if (res.data.system_logo) {
          systemLogo.value = res.data.system_logo
          localStorage.setItem('system_logo', res.data.system_logo)
        }
      }
      loaded.value = true
    } catch {
      // 加载失败时使用缓存值，静默处理
      loaded.value = true
    }
  }

  /** 强制刷新配置 */
  async function refreshConfig() {
    loaded.value = false
    await loadConfig()
  }

  return {
    systemTitle,
    systemCopyright,
    systemLogo,
    loaded,
    loadConfig,
    refreshConfig,
  }
})
