import { ref, onMounted, onUnmounted } from 'vue'

/** 断点定义 */
const breakpoints = {
  xs: 0,
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
  xxl: 1600,
} as const

export type Breakpoint = keyof typeof breakpoints
export type DeviceType = 'desktop' | 'tablet' | 'mobile'

/**
 * 响应式断点检测 Hook
 * 根据窗口宽度判断当前设备类型和激活的断点
 */
export function useBreakpoint() {
  const width = ref(window.innerWidth)
  const activeBreakpoint = ref<Breakpoint>('xl')
  const device = ref<DeviceType>('desktop')

  function update() {
    width.value = window.innerWidth
    const w = width.value

    if (w < breakpoints.md) {
      activeBreakpoint.value = w < breakpoints.sm ? 'xs' : 'sm'
      device.value = 'mobile'
    } else if (w < breakpoints.lg) {
      activeBreakpoint.value = 'md'
      device.value = 'tablet'
    } else if (w < breakpoints.xl) {
      activeBreakpoint.value = 'lg'
      device.value = 'desktop'
    } else if (w < breakpoints.xxl) {
      activeBreakpoint.value = 'xl'
      device.value = 'desktop'
    } else {
      activeBreakpoint.value = 'xxl'
      device.value = 'desktop'
    }
  }

  function isMobile() {
    return device.value === 'mobile'
  }

  function isTablet() {
    return device.value === 'tablet'
  }

  function isDesktop() {
    return device.value === 'desktop'
  }

  /** 判断当前宽度是否小于指定断点 */
  function smallerThan(breakpoint: Breakpoint): boolean {
    return width.value < breakpoints[breakpoint]
  }

  /** 判断当前宽度是否大于等于指定断点 */
  function greaterOrEqual(breakpoint: Breakpoint): boolean {
    return width.value >= breakpoints[breakpoint]
  }

  onMounted(() => {
    update()
    window.addEventListener('resize', update)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', update)
  })

  return {
    width,
    activeBreakpoint,
    device,
    isMobile,
    isTablet,
    isDesktop,
    smallerThan,
    greaterOrEqual,
    breakpoints,
  }
}
