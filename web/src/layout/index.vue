<template>
  <div class="layout-container" :class="containerClass">
    <!-- 移动端遮罩层 -->
    <transition name="fade">
      <div
        v-if="appStore.device === 'mobile' && appStore.sidebar.opened"
        class="sidebar-mask"
        @click="handleCloseSidebar"
      />
    </transition>

    <!-- 侧边栏 -->
    <Sidebar />

    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <div class="fixed-header">
        <Navbar />
        <TagsView />
      </div>

      <!-- 内容区 -->
      <AppMain />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useAppStore } from '@/store/modules/app'
import { useBreakpoint } from '@/hooks/useBreakpoint'
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
import TagsView from './components/TagsView.vue'
import AppMain from './components/AppMain.vue'

defineOptions({ name: 'Layout' })

const appStore = useAppStore()
const { device, isMobile, isTablet } = useBreakpoint()

/** 容器样式类 */
const containerClass = computed(() => ({
  'sidebar-collapsed': !appStore.sidebar.opened,
  'mobile': appStore.device === 'mobile',
  'tablet': appStore.device === 'tablet',
  'desktop': appStore.device === 'desktop',
}))

/** 监听设备类型变化，自动更新 appStore */
watch(
  device,
  (val) => {
    appStore.toggleDevice(val)
    // 移动端自动收起侧边栏
    if (val === 'mobile') {
      appStore.closeSidebar(false)
    }
    // 平板端自动收起侧边栏（图标模式）
    if (val === 'tablet') {
      appStore.closeSidebar(false)
    }
  },
  { immediate: true }
)

/** 关闭侧边栏（移动端遮罩点击） */
function handleCloseSidebar() {
  appStore.closeSidebar(false)
}
</script>

<style scoped>
.layout-container {
  display: flex;
  width: 100%;
  height: 100%;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  transition: margin-left 0.3s;
  margin-left: var(--sidebar-width);
}

.sidebar-collapsed .main-container {
  margin-left: var(--sidebar-collapsed-width);
}

.fixed-header {
  position: sticky;
  top: 0;
  z-index: 9;
  background: var(--el-bg-color);
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.02);
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

html.dark .fixed-header {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

/* 移动端遮罩层 */
.sidebar-mask {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 999;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(2px);
  transition: opacity 0.3s ease;
}

/* 遮罩层淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 移动端布局适配 */
@media screen and (max-width: 768px) {
  .main-container {
    margin-left: 0 !important;
  }
}

/* 平板端布局适配：跟随侧边栏折叠状态 */
@media screen and (min-width: 769px) and (max-width: 992px) {
  .sidebar-collapsed .main-container {
    margin-left: var(--sidebar-collapsed-width) !important;
  }
}
</style>
