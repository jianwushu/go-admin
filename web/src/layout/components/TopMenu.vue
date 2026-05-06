<template>
  <el-menu
    :default-active="activeMenu"
    mode="horizontal"
    :ellipsis="false"
    background-color="transparent"
    text-color="var(--el-text-color-regular)"
    active-text-color="var(--el-color-primary)"
    :router="true"
    class="top-menu"
  >
    <template v-for="route in menuRoutes" :key="route.path">
      <!-- 有多个子菜单的目录 -->
      <el-sub-menu
        v-if="hasChildren(route)"
        :index="resolvePath(route.path)"
      >
        <template #title>
          <el-icon v-if="getIcon(route.meta?.icon)">
            <component :is="getIcon(route.meta?.icon)" />
          </el-icon>
          <span>{{ getMenuTitle(route) }}</span>
        </template>
        <template v-for="child in getVisibleChildren(route)" :key="child.path">
          <!-- 二级菜单还有子菜单 -->
          <el-sub-menu
            v-if="hasChildren(child)"
            :index="resolveChildPath(route.path, child.path)"
          >
            <template #title>
              <el-icon v-if="getIcon(child.meta?.icon)">
                <component :is="getIcon(child.meta?.icon)" />
              </el-icon>
              <span>{{ getMenuTitle(child) }}</span>
            </template>
            <el-menu-item
              v-for="grandchild in getVisibleChildren(child)"
              :key="grandchild.path"
              :index="resolveGrandchildPath(route.path, child.path, grandchild.path)"
            >
              <el-icon v-if="getIcon(grandchild.meta?.icon)">
                <component :is="getIcon(grandchild.meta?.icon)" />
              </el-icon>
              <span>{{ getMenuTitle(grandchild) }}</span>
            </el-menu-item>
          </el-sub-menu>
          <!-- 二级菜单无子菜单 -->
          <el-menu-item
            v-else
            :index="resolveChildPath(route.path, child.path)"
          >
            <el-icon v-if="getIcon(child.meta?.icon)">
              <component :is="getIcon(child.meta?.icon)" />
            </el-icon>
            <span>{{ getMenuTitle(child) }}</span>
          </el-menu-item>
        </template>
      </el-sub-menu>

      <!-- 无子菜单或只有一个子菜单 -->
      <el-menu-item
        v-else
        :index="getSingleChildPath(route)"
      >
        <el-icon v-if="getIcon(getSingleChild(route)?.meta?.icon || route.meta?.icon)">
          <component :is="getIcon(getSingleChild(route)?.meta?.icon || route.meta?.icon)" />
        </el-icon>
        <span>{{ getMenuTitle(getSingleChild(route) || route) }}</span>
      </el-menu-item>
    </template>
  </el-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import type { RouteRecordRaw } from 'vue-router'
import { usePermissionStore } from '@/store/modules/permission'
import { getIconComponent } from '@/utils/icon'

defineOptions({ name: 'TopMenu' })

const route = useRoute()
const { t } = useI18n()
const permissionStore = usePermissionStore()

/** 当前激活菜单 */
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

/** 菜单路由列表（过滤掉隐藏的路由） */
const menuRoutes = computed(() => {
  const routes = permissionStore.routes.length > 0 ? permissionStore.routes : []
  return routes.filter((route) => !route.meta?.hidden)
})

/** 获取菜单标题 */
function getMenuTitle(route: RouteRecordRaw): string {
  const i18nKey = route.meta?.i18nKey as string
  if (i18nKey) {
    const translated = t(i18nKey)
    if (translated !== i18nKey) return translated
  }
  return (route.meta?.title as string) || ''
}

/** 获取图标组件 */
function getIcon(iconName?: string) {
  return iconName ? getIconComponent(iconName) : undefined
}

/** 获取可见子菜单 */
function getVisibleChildren(route: RouteRecordRaw): RouteRecordRaw[] {
  return (route.children || []).filter((child) => !child.meta?.hidden)
}

/** 判断是否有子菜单 */
function hasChildren(route: RouteRecordRaw): boolean {
  return getVisibleChildren(route).length > 1
}

/** 获取单个子菜单（只有一个子菜单时使用） */
function getSingleChild(route: RouteRecordRaw): RouteRecordRaw | undefined {
  const children = getVisibleChildren(route)
  if (children.length === 1) {
    return children[0]
  }
  return undefined
}

/** 获取单个子菜单的路径 */
function getSingleChildPath(route: RouteRecordRaw): string {
  const child = getSingleChild(route)
  if (child) {
    return resolveChildPath(route.path, child.path)
  }
  return resolvePath(route.path)
}

/** 解析路径 */
function resolvePath(routePath: string): string {
  if (routePath.startsWith('/')) {
    return routePath
  }
  return `/${routePath}`
}

/** 解析子菜单路径 */
function resolveChildPath(parentPath: string, childPath: string): string {
  if (childPath.startsWith('/')) {
    return childPath
  }
  const base = parentPath.endsWith('/') ? parentPath : `${parentPath}/`
  return `${base}${childPath}`
}

/** 解析三级菜单路径 */
function resolveGrandchildPath(parentPath: string, childPath: string, grandchildPath: string): string {
  if (grandchildPath.startsWith('/')) {
    return grandchildPath
  }
  const childFullPath = resolveChildPath(parentPath, childPath)
  const base = childFullPath.endsWith('/') ? childFullPath : `${childFullPath}/`
  return `${base}${grandchildPath}`
}
</script>

<style scoped>
.top-menu {
  border-bottom: none !important;
  height: 100%;
}

.top-menu :deep(.el-menu-item),
.top-menu :deep(.el-sub-menu__title) {
  height: var(--navbar-height);
  line-height: var(--navbar-height);
}

.top-menu :deep(.el-sub-menu__title) {
  border-bottom: none !important;
}

.top-menu :deep(.el-menu-item.is-active) {
  border-bottom: 2px solid var(--el-color-primary) !important;
}
</style>
