<template>
  <section class="app-main">
    <router-view v-slot="{ Component }">
      <transition name="fade-transform" mode="out-in">
        <keep-alive :include="cachedViews">
          <component :is="Component" :key="route.path" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTagsViewStore } from '@/store/modules/tagsView'

defineOptions({ name: 'AppMain' })

const route = useRoute()
const tagsViewStore = useTagsViewStore()

const cachedViews = computed(() => tagsViewStore.cachedViews)
</script>

<style scoped>
.app-main {
  flex: 1;
  overflow: auto;
  padding: 16px;
  min-height: calc(100vh - var(--navbar-height) - var(--tagsview-height));
  background-color: var(--content-bg);
  transition: background-color 0.3s ease;
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
