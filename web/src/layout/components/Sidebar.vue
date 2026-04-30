<template>
  <div class="sidebar-container" :class="sidebarClass">
    <!-- Logo -->
    <div class="sidebar-logo">
      <router-link to="/" class="sidebar-logo-link">
        <img v-if="logo" :src="logo" class="sidebar-logo-img" alt="Logo" />
        <h1 v-show="showLogoTitle" class="sidebar-logo-title">Go-Admin</h1>
      </router-link>
    </div>

    <!-- 菜单 -->
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed"
        :collapse-transition="false"
        background-color="var(--sidebar-bg)"
        text-color="var(--sidebar-text)"
        active-text-color="var(--sidebar-active-text)"
        :unique-opened="true"
        :router="true"
        mode="vertical"
      >
        <SidebarItem
          v-for="route in menuRoutes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
        />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/store/modules/app'
import { usePermissionStore } from '@/store/modules/permission'
import SidebarItem from './SidebarItem.vue'

defineOptions({ name: 'Sidebar' })

const route = useRoute()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const logo = ''

/** 是否折叠（移动端或平板端或手动折叠） */
const isCollapsed = computed(() => {
  // 移动端始终折叠（隐藏侧边栏）
  if (appStore.device === 'mobile') return true
  // 平板端始终以图标模式显示
  if (appStore.device === 'tablet') return true
  // PC端根据用户操作
  return !appStore.sidebar.opened
})

/** 是否显示 Logo 标题 */
const showLogoTitle = computed(() => {
  if (appStore.device === 'mobile') return false
  if (appStore.device === 'tablet') return false
  return appStore.sidebar.opened
})

/** 侧边栏容器样式类 */
const sidebarClass = computed(() => ({
  'is-collapse': isCollapsed.value,
  'is-mobile': appStore.device === 'mobile',
  'is-tablet': appStore.device === 'tablet',
  'is-hidden': appStore.device === 'mobile' && !appStore.sidebar.opened,
}))

/** 当前激活菜单 */
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

/** 菜单路由列表 */
const menuRoutes = computed(() => {
  return permissionStore.routes.length > 0 ? permissionStore.routes : []
})
</script>

<style scoped>
.sidebar-container {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 1001;
  width: var(--sidebar-width);
  height: 100%;
  background-color: var(--sidebar-bg);
  transition: width 0.3s, transform 0.3s, background-color 0.3s;
  overflow: hidden;
  border-right: 1px solid var(--el-border-color-light);
}

.sidebar-container.is-collapse {
  width: var(--sidebar-collapsed-width);
}

/* 移动端：侧边栏滑出/滑入 */
.sidebar-container.is-mobile {
  width: var(--sidebar-width);
  transform: translateX(0);
}

.sidebar-container.is-mobile.is-hidden {
  transform: translateX(-100%);
}

/* 平板端：固定图标模式 */
.sidebar-container.is-tablet {
  width: var(--sidebar-collapsed-width);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 50px;
  padding: 0 12px;
  overflow: hidden;
  border-bottom: 1px solid var(--el-border-color-light);
}

.sidebar-logo-link {
  display: flex;
  align-items: center;
  width: 100%;
  height: 100%;
  text-decoration: none;
}

.sidebar-logo-img {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.sidebar-logo-title {
  margin-left: 10px;
  color: var(--sidebar-text);
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-scrollbar) {
  height: calc(100% - 50px);
}

:deep(.scrollbar-wrapper) {
  overflow-x: hidden !important;
}

:deep(.el-menu--collapse .el-sub-menu__title span) {
  display: none;
}

:deep(.el-menu--collapse .el-sub-menu__title .el-sub-menu__icon-arrow) {
  display: none;
}

/* 菜单项悬停效果 */
:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: var(--sidebar-hover-bg) !important;
}

/* 激活菜单项 */
:deep(.el-menu-item.is-active) {
  background-color: var(--sidebar-active-bg) !important;
}
</style>
