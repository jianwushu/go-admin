<template>
  <section class="app-main-wrapper">
    <section class="app-main">
      <router-view v-slot="{ Component }">
        <transition name="fade-transform" mode="out-in">
          <keep-alive :include="cachedViews">
            <component :is="Component" :key="route.path" />
          </keep-alive>
        </transition>
      </router-view>
    </section>
    <!-- 版权标识 -->
    <footer v-if="configStore.systemCopyright" class="app-footer">
      {{ configStore.systemCopyright }}
    </footer>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTagsViewStore } from '@/store/modules/tagsView'
import { useConfigStore } from '@/store/modules/config'

defineOptions({ name: 'AppMain' })

const route = useRoute()
const tagsViewStore = useTagsViewStore()
const configStore = useConfigStore()

const cachedViews = computed(() => tagsViewStore.cachedViews)
</script>

<style scoped>
.app-main-wrapper {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - var(--navbar-height) - var(--tagsview-height));
}

.app-main {
  flex: 1;
  overflow: auto;
  padding: 16px;
  background-color: var(--content-bg);
  transition: background-color 0.3s ease;
}

.app-footer {
  flex-shrink: 0;
  text-align: center;
  padding: 12px 16px;
  font-size: 12px;
  color: var(--el-text-color-secondary, #909399);
  background-color: var(--content-bg);
  border-top: 1px solid var(--el-border-color-lighter, #e4e7ed);
  transition: background-color 0.3s ease, border-color 0.3s ease;
}

/* 路由切换动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
