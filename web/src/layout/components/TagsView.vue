<template>
  <div class="tags-view-container">
    <el-scrollbar class="tags-view-wrapper">
      <router-link
        v-for="tag in visitedViews"
        :key="tag.path"
        :to="tag.name ? { name: tag.name, params: tag.params, query: tag.query } : { path: tag.path, query: tag.query }"
        class="tags-view-item"
        :class="{ active: isActive(tag) }"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ tag.title }}
        <el-icon
          v-if="!tag.meta?.affix"
          class="tags-view-close"
          :size="12"
          @click.prevent.stop="closeTag(tag)"
        >
          <Close />
        </el-icon>
      </router-link>
    </el-scrollbar>

    <!-- 右键菜单 -->
    <ul v-show="contextMenu.visible" class="contextmenu" :style="contextMenuStyle">
      <li @click="refreshSelectedTag(contextMenu.tag!)">{{ t('tagsView.refresh') }}</li>
      <li v-if="!contextMenu.tag?.meta?.affix" @click="closeTag(contextMenu.tag!)">
        {{ t('tagsView.closeCurrent') }}
      </li>
      <li @click="closeOthersTags">{{ t('tagsView.closeOthers') }}</li>
      <li @click="closeAllTags">{{ t('tagsView.closeAll') }}</li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Close } from '@element-plus/icons-vue'
import { useTagsViewStore } from '@/store/modules/tagsView'
import type { TagView } from '@/types/store'

defineOptions({ name: 'TagsView' })

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const tagsViewStore = useTagsViewStore()

const visitedViews = computed(() => tagsViewStore.visitedViews)

const contextMenu = reactive({
  visible: false,
  tag: null as TagView | null,
  left: 0,
  top: 0,
})

const contextMenuStyle = computed(() => ({
  left: `${contextMenu.left}px`,
  top: `${contextMenu.top}px`,
}))

/** 判断是否为当前激活标签 */
function isActive(tag: TagView): boolean {
  return tag.path === route.path
}

/** 添加标签 */
function addTag() {
  const { name, path, query, params, meta } = route
  if (name) {
    tagsViewStore.addVisitedView({
      name: name as string,
      path,
      fullPath: route.fullPath,
      title: (meta?.title as string) || 'no-name',
      meta: meta as any,
      query: query as Record<string, string>,
      params: params as Record<string, string>,
    })
    tagsViewStore.addCachedView({
      name: name as string,
      path,
      fullPath: route.fullPath,
      title: (meta?.title as string) || 'no-name',
      meta: meta as any,
    })
  }
}

/** 关闭标签 */
function closeTag(tag: TagView) {
  tagsViewStore.delVisitedView(tag)
  tagsViewStore.delCachedView(tag)
  if (isActive(tag)) {
    const views = tagsViewStore.visitedViews
    const lastView = views[views.length - 1]
    if (lastView) {
      router.push(lastView.path)
    } else {
      router.push('/')
    }
  }
  contextMenu.visible = false
}

/** 刷新选中的标签页 */
function refreshSelectedTag(tag: TagView) {
  tagsViewStore.delCachedView(tag)
  const { fullPath } = tag
  router.replace({ path: '/redirect' + fullPath }).catch(() => {})
  contextMenu.visible = false
}

/** 关闭其他标签 */
function closeOthersTags() {
  if (contextMenu.tag) {
    router.push(contextMenu.tag.path)
    tagsViewStore.delOthersViews(contextMenu.tag)
  }
  contextMenu.visible = false
}

/** 关闭所有标签 */
function closeAllTags() {
  tagsViewStore.delAllViews()
  router.push('/')
  contextMenu.visible = false
}

/** 打开右键菜单 */
function openMenu(tag: TagView, e: MouseEvent) {
  contextMenu.tag = tag
  contextMenu.left = e.clientX
  contextMenu.top = e.clientY
  contextMenu.visible = true
}

/** 关闭右键菜单 */
function closeMenu() {
  contextMenu.visible = false
}

watch(
  () => route.path,
  () => {
    addTag()
  },
  { immediate: true }
)

watch(
  () => contextMenu.visible,
  (val) => {
    if (val) {
      document.addEventListener('click', closeMenu, { once: true })
    }
  }
)
</script>

<style scoped>
.tags-view-container {
  width: 100%;
  background: inherit;
  border-bottom: 1px solid var(--el-border-color-light);
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.04);
  transition: border-color 0.3s ease;
}

.tags-view-wrapper {
  display: flex;
  align-items: center;
  height: var(--tagsview-height);
  padding: 0 8px;
}

.tags-view-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 26px;
  padding: 0 10px;
  margin: 0 2px;
  font-size: 12px;
  color: var(--el-text-color-regular);
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
  text-decoration: none;
  transition: all 0.2s;
}

.tags-view-item:hover {
  color: var(--el-color-primary);
}

.tags-view-item.active {
  color: #fff;
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

.tags-view-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  transition: all 0.2s;
}

.tags-view-close:hover {
  background-color: rgba(255, 255, 255, 0.3);
  color: #fff;
}

.contextmenu {
  position: fixed;
  z-index: 3000;
  margin: 0;
  padding: 5px 0;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  list-style: none;
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.contextmenu li {
  padding: 7px 16px;
  cursor: pointer;
  transition: background 0.15s;
}

.contextmenu li:hover {
  background: var(--el-fill-color-light);
  color: var(--el-color-primary);
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .tags-view-item {
    padding: 0 8px;
    font-size: 11px;
    height: 24px;
  }
  
  .tags-view-close {
    width: 14px;
    height: 14px;
  }
}
</style>
