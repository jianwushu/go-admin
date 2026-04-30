/**
 * 主题工具 - 管理白天/黑夜模式切换
 * 集成 Element-Plus 暗黑主题，支持清新简约风格
 */

export type ThemeMode = 'light' | 'dark'

const THEME_KEY = 'theme'

/**
 * 获取当前主题模式
 * 默认为白天模式
 */
export function getThemeMode(): ThemeMode {
  const stored = localStorage.getItem(THEME_KEY)
  if (stored === 'dark' || stored === 'light') {
    return stored
  }
  // 默认白天模式
  return 'light'
}

/**
 * 设置主题模式
 */
export function setThemeMode(mode: ThemeMode): void {
  localStorage.setItem(THEME_KEY, mode)
  applyTheme(mode)
}

/**
 * 切换主题模式
 */
export function toggleThemeMode(): ThemeMode {
  const current = getThemeMode()
  const next = current === 'light' ? 'dark' : 'light'
  setThemeMode(next)
  return next
}

/**
 * 应用主题到 DOM
 */
export function applyTheme(mode: ThemeMode): void {
  const html = document.documentElement
  
  if (mode === 'dark') {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
}

/**
 * 初始化主题
 * 在应用启动时调用，从 localStorage 读取并应用主题
 */
export function initTheme(): void {
  const mode = getThemeMode()
  applyTheme(mode)
}
