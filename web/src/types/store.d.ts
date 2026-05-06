/** 布局模式 */
export type LayoutMode = 'sidebar' | 'top' | 'mixed'

/** 用户 Store 状态 */
export interface UserState {
  token: string
  refreshToken: string
  userInfo: UserInfo | null
  roles: string[]
  permissions: string[]
}

/** 应用 Store 状态 */
export interface AppState {
  sidebar: {
    opened: boolean
    withoutAnimation: boolean
  }
  device: 'desktop' | 'mobile' | 'tablet'
  theme: 'light' | 'dark'
  language: string
  size: 'large' | 'default' | 'small'
  layoutMode: LayoutMode
  primaryColor: string
}

/** 权限 Store 状态 */
export interface PermissionState {
  routes: RouteRecordRaw[]
  addRoutes: RouteRecordRaw[]
  menus: MenuItem[]
  buttons: string[]
}

/** 标签页 */
export interface TagView {
  path: string
  fullPath: string
  name: string
  title: string
  meta: RouteMeta
  query?: Record<string, string>
  params?: Record<string, string>
}

/** 标签页 Store 状态 */
export interface TagsViewState {
  visitedViews: TagView[]
  cachedViews: string[]
}

import type { RouteRecordRaw, RouteMeta } from 'vue-router'
import type { UserInfo, MenuItem } from './api.d.ts'
