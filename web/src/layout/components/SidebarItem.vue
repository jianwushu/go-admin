<template>
  <!-- 隐藏菜单 -->
  <template v-if="!item.meta?.hidden">
    <!-- 只有一个子菜单或没有子菜单 -->
    <template v-if="hasOneShowingChild(item)">
      <el-menu-item :index="resolvePath(onlyOneChild!.path)">
        <el-icon v-if="getIcon(onlyOneChild!.meta?.icon)">
          <component :is="getIcon(onlyOneChild!.meta?.icon)" />
        </el-icon>
        <template #title>
          <span>{{ onlyOneChild!.meta?.title }}</span>
        </template>
      </el-menu-item>
    </template>

    <!-- 多个子菜单 -->
    <el-sub-menu v-else :index="resolvePath(item.path)">
      <template #title>
        <el-icon v-if="getIcon(item.meta?.icon)">
          <component :is="getIcon(item.meta?.icon)" />
        </el-icon>
        <span>{{ item.meta?.title }}</span>
      </template>

      <SidebarItem
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        :base-path="resolvePath(child.path)"
      />
    </el-sub-menu>
  </template>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { isExternal } from '@/utils/validate'
import { getIconComponent } from '@/utils/icon'

defineOptions({ name: 'SidebarItem' })

const props = defineProps<{
  item: RouteRecordRaw
  basePath: string
}>()

const onlyOneChild = ref<RouteRecordRaw>()

/** 获取图标组件 */
function getIcon(iconName?: string) {
  return iconName ? getIconComponent(iconName) : undefined
}

/** 判断是否只有一个需要显示的子菜单 */
function hasOneShowingChild(parent: RouteRecordRaw): boolean {
  const children = parent.children || []
  const showingChildren = children.filter((item) => !item.meta?.hidden)

  // 只有一个子菜单时直接显示
  if (showingChildren.length === 1) {
    onlyOneChild.value = showingChildren[0]
    return true
  }

  // 没有子菜单时显示父级
  if (showingChildren.length === 0) {
    onlyOneChild.value = { ...parent, path: '' }
    return true
  }

  return false
}

/** 解析路径 */
function resolvePath(routePath: string): string {
  if (isExternal(routePath)) {
    return routePath
  }
  if (isExternal(props.basePath)) {
    return props.basePath
  }
  // 拼接基础路径和子路径
  const base = props.basePath.endsWith('/') ? props.basePath : `${props.basePath}/`
  if (routePath && !routePath.startsWith('/')) {
    return `${base}${routePath}`
  }
  return routePath || props.basePath
}
</script>
