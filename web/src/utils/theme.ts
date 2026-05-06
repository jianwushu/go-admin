/**
 * 主题工具 - 管理白天/黑夜模式切换 + 主题颜色
 * 集成 Element-Plus 暗黑主题，支持清新简约风格
 */

export type ThemeMode = 'light' | 'dark'

const THEME_KEY = 'theme'
const PRIMARY_COLOR_KEY = 'primaryColor'

/** 预设主题色 */
export const PRESET_COLORS = [
  { name: '极光蓝', value: '#4096ff' },
  { name: '翡翠绿', value: '#52c41a' },
  { name: '星空紫', value: '#722ed1' },
  { name: '落日橙', value: '#fa8c16' },
  { name: '玫瑰红', value: '#f5222d' },
  { name: '天际青', value: '#13c2c2' },
  { name: '暗夜黑', value: '#434343' },
]

/** 默认主题色 */
const DEFAULT_PRIMARY_COLOR = '#4096ff'

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

// ==================== 主题颜色 ====================

/**
 * 获取当前主题色
 */
export function getPrimaryColor(): string {
  return localStorage.getItem(PRIMARY_COLOR_KEY) || DEFAULT_PRIMARY_COLOR
}

/**
 * 将十六进制颜色转换为 RGB
 */
function hexToRgb(hex: string): { r: number; g: number; b: number } {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  if (!result) {
    return { r: 64, g: 150, b: 255 } // 默认蓝色
  }
  return {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16),
  }
}

/**
 * 混合颜色（与白色或黑色混合生成派生色）
 * @param color 原始颜色
 * @param ratio 混合比例 (0-1)，越大越接近混合目标
 * @param target 混合目标颜色，默认白色
 */
function mixColor(
  color: string,
  ratio: number,
  target: { r: number; g: number; b: number } = { r: 255, g: 255, b: 255 }
): string {
  const { r, g, b } = hexToRgb(color)
  const mixed = {
    r: Math.round(r + (target.r - r) * ratio),
    g: Math.round(g + (target.g - g) * ratio),
    b: Math.round(b + (target.b - b) * ratio),
  }
  return `#${mixed.r.toString(16).padStart(2, '0')}${mixed.g.toString(16).padStart(2, '0')}${mixed.b.toString(16).padStart(2, '0')}`
}

/**
 * 应用主题色到 DOM
 * 动态修改 CSS 变量，与 Element Plus 的 css-vars 机制兼容
 */
export function applyPrimaryColor(color: string): void {
  const root = document.documentElement

  // 主色
  root.style.setProperty('--el-color-primary', color)

  // light 派生色（与白色混合，比例递增）
  root.style.setProperty('--el-color-primary-light-3', mixColor(color, 0.3))
  root.style.setProperty('--el-color-primary-light-5', mixColor(color, 0.5))
  root.style.setProperty('--el-color-primary-light-7', mixColor(color, 0.7))
  root.style.setProperty('--el-color-primary-light-8', mixColor(color, 0.8))
  root.style.setProperty('--el-color-primary-light-9', mixColor(color, 0.9))

  // dark 派生色（与黑色混合）
  root.style.setProperty('--el-color-primary-dark-2', mixColor(color, 0.2, { r: 0, g: 0, b: 0 }))

  // 侧边栏相关变量
  const { r, g, b } = hexToRgb(color)
  root.style.setProperty('--sidebar-active-text', color)
  root.style.setProperty('--sidebar-hover-bg', `rgba(${r}, ${g}, ${b}, 0.06)`)
  root.style.setProperty('--sidebar-active-bg', `rgba(${r}, ${g}, ${b}, 0.1)`)
}

/**
 * 初始化主题色
 * 在应用启动时调用，从 localStorage 读取并应用主题色
 */
export function initPrimaryColor(): void {
  const color = getPrimaryColor()
  applyPrimaryColor(color)
}
