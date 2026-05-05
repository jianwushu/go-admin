import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MenuItem } from '@/types/api'
import type { RouteRecordRaw } from 'vue-router'
import { staticRoutes } from '@/router/staticRoutes'
import { getUserMenus } from '@/api/auth'
import { useUserStore } from '@/store/modules/user'
import Layout from '@/layout/index.vue'

/** 动态导入组件映射 */
const componentModules = import.meta.glob('@/views/**/*.vue')

export const usePermissionStore = defineStore('permission', () => {
  /** 所有可访问路由（静态 + 动态） */
  const routes = ref<RouteRecordRaw[]>([...staticRoutes])
  /** 动态添加的路由（需要 addRoute） */
  const addRoutes = ref<RouteRecordRaw[]>([])
  /** 菜单数据 */
  const menus = ref<MenuItem[]>([])
  /** 按钮权限列表 */
  const buttons = ref<string[]>([])

  /** 设置动态路由 */
  function setRoutes(newRoutes: RouteRecordRaw[]) {
    addRoutes.value = newRoutes
    routes.value = [...staticRoutes, ...newRoutes]
  }

  /** 设置菜单数据 */
  function setMenus(menuData: MenuItem[]) {
    menus.value = menuData
  }

  /** 设置按钮权限 */
  function setButtons(buttonList: string[]) {
    buttons.value = buttonList
  }

  /** 检查是否有某个按钮权限 */
  function hasPermission(permission: string): boolean {
    return buttons.value.includes('*:*:*') || buttons.value.includes(permission)
  }

  /** 重置权限状态 */
  function resetPermission() {
    routes.value = [...staticRoutes]
    addRoutes.value = []
    menus.value = []
    buttons.value = []
  }

  /**
   * 生成动态路由
   * 从后端获取菜单数据，转换为 vue-router 路由配置
   */
  async function generateRoutes(): Promise<RouteRecordRaw[]> {
    const { data: res } = await getUserMenus()
    const menuData = res.data as MenuItem[]

    // 保存菜单数据
    setMenus(menuData)

    // 提取按钮权限（从菜单中提取 type=2 的按钮权限）
    const perms: string[] = []
    extractPermissions(menuData, perms)

    // 合并用户信息中的权限列表（超级管理员会包含 *:*:*）
    const userStore = useUserStore()
    if (userStore.permissions && userStore.permissions.length > 0) {
      userStore.permissions.forEach((perm) => {
        if (!perms.includes(perm)) {
          perms.push(perm)
        }
      })
    }

    setButtons(perms)

    // 将菜单转换为路由
    const asyncRoutes = convertMenusToRoutes(menuData)
    setRoutes(asyncRoutes)

    return asyncRoutes
  }

  /**
   * 递归提取按钮权限
   */
  function extractPermissions(menus: MenuItem[], perms: string[]) {
    menus.forEach((menu) => {
      if (menu.type === 2 && menu.perms) {
        perms.push(menu.perms)
      }
      if (menu.children && menu.children.length > 0) {
        extractPermissions(menu.children, perms)
      }
    })
  }

  /**
   * 将菜单数据转换为路由配置
   */
  function convertMenusToRoutes(menus: MenuItem[]): RouteRecordRaw[] {
    const routes: RouteRecordRaw[] = []

    menus.forEach((menu) => {
      // 只处理目录和菜单类型
      if (menu.type !== 0 && menu.type !== 1) return

      // 用菜单 ID 生成唯一 name，避免不同菜单同名导致 addRoute 冲突
      const route: RouteRecordRaw = {
        path: menu.path,
        name: `menu_${menu.id}`,
        component: resolveComponent(menu.component),
        meta: {
          title: menu.name,
          i18nKey: menu.i18nKey || '',
          icon: menu.icon,
          hidden: menu.visible === 0,
          sort: menu.sort,
        },
        children: [],
      }

      // 递归处理子菜单
      if (menu.children && menu.children.length > 0) {
        route.children = convertMenusToRoutes(menu.children)
      }

      // 目录类型且有子路由时，重定向到第一个子路由，避免进入空白页
      if (menu.type === 0 && route.children && route.children.length > 0) {
        const firstChild = route.children[0]
        const childPath = firstChild.path.startsWith('/')
          ? firstChild.path
          : `${menu.path.replace(/\/$/, '')}/${firstChild.path}`
        route.redirect = childPath
      }

      routes.push(route)
    })

    return routes
  }

  /**
   * 解析组件路径
   * 目录类型返回 Layout，菜单类型动态导入
   */
  function resolveComponent(componentPath: string): any {
    if (!componentPath) {
      return Layout
    }
    // 转换为实际的组件路径
    const path = `/src/views/${componentPath}.vue`
    if (componentModules[path]) {
      return componentModules[path]
    }
    // 如果找不到组件，返回空壳组件
    console.warn(`组件未找到: ${path}`)
    return () => import('@/views/error/404.vue')
  }

  return {
    routes,
    addRoutes,
    menus,
    buttons,
    setRoutes,
    setMenus,
    setButtons,
    hasPermission,
    resetPermission,
    generateRoutes,
  }
})
