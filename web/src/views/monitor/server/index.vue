<template>
  <div class="monitor-container">
    <!-- 页面标题和操作栏 -->
    <div class="page-header section-spacing">
      <div class="page-header-content">
        <h1 class="page-title">{{ t('monitor.title') }}</h1>
        <div class="page-actions">
          <el-switch
            v-model="autoRefresh"
            :active-text="t('monitor.autoRefresh')"
            @change="handleAutoRefreshChange"
          />
          <el-button type="primary" :icon="Refresh" @click="fetchData">
            {{ t('monitor.refresh') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && !serverInfo" class="flex h-64 items-center justify-center">
      <el-icon class="is-loading" :size="40" color="#4096ff">
        <Loading />
      </el-icon>
    </div>

    <template v-else-if="serverInfo">
      <!-- 统计卡片 -->
      <div class="stat-grid section-spacing">
        <div
          v-for="item in statCards"
          :key="item.key"
          class="stat-item"
        >
          <div class="stat-item-inner">
            <el-statistic :value="item.value">
              <template #title>
                <div class="stat-title">
                  <el-icon :size="16" :color="item.color">
                    <component :is="item.icon" />
                  </el-icon>
                  <span>{{ item.label }}</span>
                </div>
              </template>
              <template #suffix>
                <span class="stat-suffix">{{ item.suffix }}</span>
              </template>
            </el-statistic>
          </div>
          <div class="stat-footer">
            <el-progress
              :percentage="Number(item.rate.toFixed(1))"
              :color="getProgressColor(item.rate)"
              :stroke-width="6"
              :show-text="false"
            />
          </div>
        </div>
      </div>

      <!-- 中间区域：CPU + 内存 -->
      <div class="content-grid section-spacing">
        <!-- CPU 信息 -->
        <div class="content-section">
          <div class="section-header">
            <el-icon :size="18" class="text-blue-500"><Cpu /></el-icon>
            <span>{{ t('monitor.cpu') }}</span>
          </div>
          <div class="mb-4 flex items-center justify-between">
            <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.cpuCores') }}</span>
            <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.cpu.cores }}</span>
          </div>
          <div class="mb-4 flex items-center justify-between">
            <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.cpuUsage') }}</span>
            <span class="font-medium" :class="getUsageColor(serverInfo.cpu.usedRate)">
              {{ serverInfo.cpu.usedRate.toFixed(2) }}%
            </span>
          </div>
          <el-progress
            :percentage="Number(serverInfo.cpu.usedRate.toFixed(1))"
            :color="getProgressColor(serverInfo.cpu.usedRate)"
            :stroke-width="20"
            :text-inside="true"
          />
          <!-- 各核心使用率 -->
          <div v-if="serverInfo.cpu.usage.length > 0" class="mt-4">
            <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">{{ t('monitor.coreUsage') }}</div>
            <div class="core-grid">
              <div
                v-for="(usage, index) in serverInfo.cpu.usage"
                :key="index"
                class="core-item"
              >
                <div class="core-label">Core {{ index }}</div>
                <div class="core-value" :class="getUsageColor(usage)">
                  {{ usage.toFixed(1) }}%
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 内存信息 -->
        <div class="content-section">
          <div class="section-header">
            <el-icon :size="18" class="text-green-500"><Coin /></el-icon>
            <span>{{ t('monitor.memory') }}</span>
          </div>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.memoryTotal') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ formatMB(serverInfo.memory.total) }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.memoryUsed') }}</span>
              <span class="font-medium" :class="getUsageColor(serverInfo.memory.usedRate)">
                {{ formatMB(serverInfo.memory.used) }}
              </span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.memoryFree') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ formatMB(serverInfo.memory.free) }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.memoryUsage') }}</span>
              <span class="font-medium" :class="getUsageColor(serverInfo.memory.usedRate)">
                {{ serverInfo.memory.usedRate.toFixed(2) }}%
              </span>
            </div>
          </div>
          <el-progress
            class="mt-4"
            :percentage="Number(serverInfo.memory.usedRate.toFixed(1))"
            :color="getProgressColor(serverInfo.memory.usedRate)"
            :stroke-width="20"
            :text-inside="true"
          />
        </div>
      </div>

      <!-- 磁盘信息 -->
      <div class="content-section section-spacing">
        <div class="section-header">
          <el-icon :size="18" class="text-orange-500"><Box /></el-icon>
          <span>{{ t('monitor.disk') }}</span>
        </div>
        <el-table :data="serverInfo.disk" stripe>
          <el-table-column :label="t('monitor.diskMountPoint')" prop="mountPoint" min-width="120" />
          <el-table-column :label="t('monitor.diskTotal')" min-width="100">
            <template #default="{ row }">{{ formatGB(row.total) }}</template>
          </el-table-column>
          <el-table-column :label="t('monitor.diskUsed')" min-width="100">
            <template #default="{ row }">
              <span :class="getUsageColor(row.usedRate)">{{ formatGB(row.used) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('monitor.diskFree')" min-width="100">
            <template #default="{ row }">{{ formatGB(row.free) }}</template>
          </el-table-column>
          <el-table-column :label="t('monitor.diskUsage')" min-width="200">
            <template #default="{ row }">
              <el-progress
                :percentage="Number(row.usedRate.toFixed(1))"
                :color="getProgressColor(row.usedRate)"
                :stroke-width="16"
                :text-inside="true"
              />
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 底部区域：Go 运行时 + 服务状态 -->
      <div class="content-grid section-spacing">
        <!-- Go 运行时信息 -->
        <div class="content-section">
          <div class="section-header">
            <el-icon :size="18" class="text-cyan-500"><Monitor /></el-icon>
            <span>{{ t('monitor.goRuntime') }}</span>
          </div>
          <div class="runtime-grid">
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.goVersion') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.goVersion }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.os') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.os }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.arch') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.arch }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.goroutines') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.goroutines }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.heapAlloc') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.heapAlloc }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.heapSys') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.heapSys }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.heapIdle') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.heapIdle }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.heapInuse') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.heapInuse }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.numGC') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.numGC }}</span>
            </div>
            <div class="runtime-item">
              <span class="runtime-label">{{ t('monitor.lastGC') }}</span>
              <span class="runtime-value">{{ serverInfo.goRuntime.lastGC }}</span>
            </div>
          </div>
        </div>

        <!-- 服务状态 -->
        <div class="content-section">
          <div class="section-header">
            <el-icon :size="18"><Connection /></el-icon>
            <span>{{ t('monitor.serviceStatus') }}</span>
          </div>
          <div class="service-list">
            <!-- 数据库状态 -->
            <div class="service-item">
              <div class="service-header">
                <div class="service-info">
                  <el-icon :size="20" class="text-blue-500"><Coin /></el-icon>
                  <span class="service-name">{{ t('monitor.db') }}</span>
                </div>
                <el-tag :type="serverInfo.db.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ serverInfo.db.status === 'online' ? t('monitor.online') : t('monitor.offline') }}
                </el-tag>
              </div>
              <p class="service-desc">{{ serverInfo.db.message }}</p>
            </div>
            <!-- Redis 状态 -->
            <div class="service-item">
              <div class="service-header">
                <div class="service-info">
                  <el-icon :size="20" class="text-red-500"><Connection /></el-icon>
                  <span class="service-name">{{ t('monitor.redis') }}</span>
                </div>
                <el-tag :type="serverInfo.redis.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ serverInfo.redis.status === 'online' ? t('monitor.online') : t('monitor.offline') }}
                </el-tag>
              </div>
              <p class="service-desc">{{ serverInfo.redis.message }}</p>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, Loading, Cpu, Coin, Box, Monitor, Connection } from '@element-plus/icons-vue'
import { getServerMonitor } from '@/api/monitor'
import type { ServerMonitorInfo } from '@/api/monitor'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const loading = ref(false)
const serverInfo = ref<ServerMonitorInfo | null>(null)
const autoRefresh = ref(false)
let refreshTimer: ReturnType<typeof setInterval> | null = null

/** 统计卡片配置 */
const statCards = computed(() => {
  if (!serverInfo.value) return []
  return [
    {
      key: 'cpu',
      icon: markRaw(Cpu),
      color: '#4096ff',
      label: t('monitor.cpuUsage'),
      value: Number(serverInfo.value.cpu.usedRate.toFixed(1)),
      suffix: '%',
      rate: serverInfo.value.cpu.usedRate,
    },
    {
      key: 'memory',
      icon: markRaw(Coin),
      color: '#52c41a',
      label: t('monitor.memoryUsage'),
      value: Number(serverInfo.value.memory.usedRate.toFixed(1)),
      suffix: '%',
      rate: serverInfo.value.memory.usedRate,
    },
    {
      key: 'disk',
      icon: markRaw(Box),
      color: '#faad14',
      label: t('monitor.diskUsage'),
      value: serverInfo.value.disk.length > 0 ? Number(serverInfo.value.disk[0].usedRate.toFixed(1)) : 0,
      suffix: '%',
      rate: serverInfo.value.disk.length > 0 ? serverInfo.value.disk[0].usedRate : 0,
    },
    {
      key: 'goroutines',
      icon: markRaw(Monitor),
      color: '#13c2c2',
      label: t('monitor.goroutines'),
      value: serverInfo.value.goRuntime.goroutines,
      suffix: '',
      rate: Math.min(serverInfo.value.goRuntime.goroutines / 10, 100), // 假设1000为满
    },
  ]
})

/** 获取服务器监控数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data } = await getServerMonitor()
    if (data.code === 0) {
      serverInfo.value = data.data
    } else {
      ElMessage.error(data.msg || t('monitor.fetchFailed'))
    }
  } catch (error: any) {
    ElMessage.error(error.message || t('monitor.fetchFailed'))
  } finally {
    loading.value = false
  }
}

/** 自动刷新切换 */
function handleAutoRefreshChange(val: string | number | boolean) {
  if (val) {
    refreshTimer = setInterval(fetchData, 10000) // 每10秒刷新
  } else {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }
}

/** 格式化 MB */
function formatMB(mb: number): string {
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb + ' MB'
}

/** 格式化 GB */
function formatGB(gb: number): string {
  return gb + ' GB'
}

/** 获取使用率颜色 */
function getUsageColor(rate: number): string {
  if (rate >= 90) return 'text-red-500'
  if (rate >= 70) return 'text-orange-500'
  return 'text-green-500'
}

/** 获取进度条颜色 */
function getProgressColor(rate: number): string {
  if (rate >= 90) return '#f56c6c'
  if (rate >= 70) return '#e6a23c'
  return '#67c23a'
}

onMounted(() => {
  fetchData()
})

onBeforeUnmount(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.monitor-container {
  width: 100%;
}

/* 统一间距 */
.section-spacing {
  margin-bottom: 24px;
}

/* 页面标题 */
.page-header {
  margin-bottom: 24px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin: 0;
}

.page-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 统计网格 */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1200px) {
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
}

.stat-item {
  background: var(--el-bg-color);
  border-radius: 10px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
  overflow: hidden;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.06), 0 2px 4px -1px rgba(0, 0, 0, 0.04);
}

.stat-item-inner {
  padding: 20px 20px 16px;
}

.stat-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.stat-suffix {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.stat-footer {
  padding: 10px 20px;
  border-top: 1px solid var(--el-border-color-extra-light);
}

/* 内容网格 */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 992px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

.content-section {
  background: var(--el-bg-color);
  border-radius: 10px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
  padding: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-extra-light);
}

/* 核心使用率网格 */
.core-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

@media (max-width: 768px) {
  .core-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.core-item {
  background: var(--el-fill-color-light);
  border-radius: 8px;
  padding: 12px;
  text-align: center;
}

.core-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.core-value {
  font-size: 14px;
  font-weight: 600;
}

/* 运行时信息网格 */
.runtime-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

@media (max-width: 768px) {
  .runtime-grid {
    grid-template-columns: 1fr;
  }
}

.runtime-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.runtime-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.runtime-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

/* 服务状态列表 */
.service-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.service-item {
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.service-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.service-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.service-desc {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0;
}
</style>
