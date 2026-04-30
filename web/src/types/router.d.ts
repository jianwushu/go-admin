import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /** 菜单标题 */
    title?: string
    /** 菜单图标 */
    icon?: string
    /** 是否隐藏 */
    hidden?: boolean
    /** 是否缓存 */
    keepAlive?: boolean
    /** 权限标识 */
    permission?: string
    /** 面包屑 */
    breadcrumb?: boolean
    /** 激活菜单 */
    activeMenu?: string
    /** 标签页是否固定 */
    affix?: boolean
  }
}
