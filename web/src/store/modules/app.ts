import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { getThemeMode, setThemeMode, type ThemeMode } from '@/utils/theme'

export const useAppStore = defineStore('app', () => {
  const sidebar = reactive({
    opened: localStorage.getItem('sidebarStatus') !== 'closed',
    withoutAnimation: false,
  })

  const device = ref<'desktop' | 'mobile' | 'tablet'>('desktop')
  const theme = ref<ThemeMode>(getThemeMode())
  const language = ref<string>(localStorage.getItem('language') || 'zh-CN')
  const size = ref<'large' | 'default' | 'small'>(
    (localStorage.getItem('size') as 'large' | 'default' | 'small') || 'default'
  )

  /** 切换侧边栏 */
  function toggleSidebar() {
    sidebar.opened = !sidebar.opened
    sidebar.withoutAnimation = false
    localStorage.setItem('sidebarStatus', sidebar.opened ? 'opened' : 'closed')
  }

  /** 关闭侧边栏 */
  function closeSidebar(withoutAnimation: boolean) {
    sidebar.opened = false
    sidebar.withoutAnimation = withoutAnimation
    localStorage.setItem('sidebarStatus', 'closed')
  }

  /** 切换设备类型 */
  function toggleDevice(val: 'desktop' | 'mobile' | 'tablet') {
    device.value = val
  }

  /** 切换主题（白天/黑夜） */
  function toggleTheme() {
    const nextMode: ThemeMode = theme.value === 'light' ? 'dark' : 'light'
    theme.value = nextMode
    setThemeMode(nextMode)
  }

  /** 设置主题 */
  function setTheme(mode: ThemeMode) {
    theme.value = mode
    setThemeMode(mode)
  }

  /** 设置语言 */
  function setLanguage(lang: string) {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  /** 设置尺寸 */
  function setSize(s: 'large' | 'default' | 'small') {
    size.value = s
    localStorage.setItem('size', s)
  }

  return {
    sidebar,
    device,
    theme,
    language,
    size,
    toggleSidebar,
    closeSidebar,
    toggleDevice,
    toggleTheme,
    setTheme,
    setLanguage,
    setSize,
  }
})
