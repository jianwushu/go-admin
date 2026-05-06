<template>
  <el-drawer
    v-model="visible"
    :title="t('settings.title')"
    size="320px"
    direction="rtl"
  >
    <!-- 主题颜色 -->
    <div class="settings-section">
      <h4 class="settings-section-title">{{ t('settings.themeColor') }}</h4>
      <div class="color-grid">
        <div
          v-for="color in PRESET_COLORS"
          :key="color.value"
          class="color-item"
          :class="{ active: appStore.primaryColor === color.value }"
          :style="{ backgroundColor: color.value }"
          :title="color.name"
          @click="handleSetColor(color.value)"
        >
          <el-icon v-if="appStore.primaryColor === color.value" class="color-check">
            <Check />
          </el-icon>
        </div>
      </div>
    </div>

    <!-- 布局模式 -->
    <div class="settings-section">
      <h4 class="settings-section-title">{{ t('settings.layoutMode') }}</h4>
      <div class="layout-grid">
        <div
          v-for="mode in layoutModes"
          :key="mode.value"
          class="layout-item"
          :class="{ active: appStore.layoutMode === mode.value }"
          @click="handleSetLayoutMode(mode.value)"
        >
          <div class="layout-preview" :class="`layout-${mode.value}`">
            <div class="layout-preview-header"></div>
            <div class="layout-preview-body">
              <div v-if="mode.value !== 'top'" class="layout-preview-sidebar"></div>
              <div class="layout-preview-content"></div>
            </div>
          </div>
          <span class="layout-label">{{ mode.label }}</span>
        </div>
      </div>
    </div>

    <!-- 重置按钮 -->
    <template #footer>
      <el-button @click="handleReset">{{ t('settings.reset') }}</el-button>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check } from '@element-plus/icons-vue'
import { useAppStore } from '@/store/modules/app'
import { PRESET_COLORS } from '@/utils/theme'
import type { LayoutMode } from '@/types/store'

defineOptions({ name: 'SettingsDrawer' })

const { t } = useI18n()
const appStore = useAppStore()

const visible = defineModel<boolean>({ default: false })

/** 布局模式列表 */
const layoutModes = computed(() => [
  { value: 'sidebar' as LayoutMode, label: t('settings.sidebarMode') },
  { value: 'top' as LayoutMode, label: t('settings.topMode') },
  { value: 'mixed' as LayoutMode, label: t('settings.mixedMode') },
])

/** 设置主题色 */
function handleSetColor(color: string) {
  appStore.setColor(color)
}

/** 设置布局模式 */
function handleSetLayoutMode(mode: LayoutMode) {
  appStore.setLayoutMode(mode)
}

/** 重置为默认配置 */
function handleReset() {
  appStore.setColor('#4096ff')
  appStore.setLayoutMode('sidebar')
}
</script>

<style scoped>
.settings-section {
  margin-bottom: 24px;
}

.settings-section-title {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

/* 主题色网格 */
.color-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.color-item {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, box-shadow 0.2s;
  position: relative;
}

.color-item:hover {
  transform: scale(1.1);
}

.color-item.active {
  box-shadow: 0 0 0 2px var(--el-bg-color), 0 0 0 4px currentColor;
}

.color-check {
  color: #fff;
  font-size: 16px;
}

/* 布局模式网格 */
.layout-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.layout-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 2px solid var(--el-border-color-light);
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.layout-item:hover {
  border-color: var(--el-color-primary-light-5);
}

.layout-item.active {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

/* 布局预览缩略图 */
.layout-preview {
  width: 64px;
  height: 48px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.layout-preview-header {
  height: 10px;
  background: var(--el-color-primary-light-7);
}

.layout-preview-body {
  display: flex;
  height: calc(100% - 10px);
}

.layout-preview-sidebar {
  width: 16px;
  background: var(--el-color-primary-light-5);
}

.layout-preview-content {
  flex: 1;
  background: var(--el-fill-color-light);
}

.layout-label {
  font-size: 13px;
  color: var(--el-text-color-regular);
}

.layout-item.active .layout-label {
  color: var(--el-color-primary);
  font-weight: 500;
}
</style>
