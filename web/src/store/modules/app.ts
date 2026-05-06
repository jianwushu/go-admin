import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { getThemeMode, setThemeMode, type ThemeMode, getPrimaryColor, applyPrimaryColor } from '@/utils/theme'
import type { LayoutMode } from '@/types/store'

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

  /** 布局模式 */
  const layoutMode = ref<LayoutMode>(
    (localStorage.getItem('layoutMode') as LayoutMode) || 'sidebar'
  )

  /** 主题色 */
  const primaryColor = ref<string>(getPrimaryColor())

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

  /** 设置布局模式 */
  function setLayoutMode(mode: LayoutMode) {
    layoutMode.value = mode
    localStorage.setItem('layoutMode', mode)
    // 移动端自动降级为侧边栏模式
    if (device.value === 'mobile' && mode !== 'sidebar') {
      layoutMode.value = 'sidebar'
      localStorage.setItem('layoutMode', 'sidebar')
    }
  }

  /** 设置主题色 */
  function setColor(color: string) {
    primaryColor.value = color
    localStorage.setItem('primaryColor', color)
    applyPrimaryColor(color)
  }

  return {
    sidebar,
    device,
    theme,
    language,
    size,
    layoutMode,
    primaryColor,
    toggleSidebar,
    closeSidebar,
    toggleDevice,
    toggleTheme,
    setTheme,
    setLanguage,
    setSize,
    setLayoutMode,
    setColor,
  }
})
