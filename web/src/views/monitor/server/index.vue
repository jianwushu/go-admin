<template>
  <div class="p-4">
    <!-- 页面标题和操作栏 -->
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-xl font-bold" style="color: var(--el-text-color-primary)">
        {{ t('monitor.title') }}
      </h1>
      <div class="flex items-center gap-3">
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

    <!-- 加载状态 -->
    <div v-if="loading && !serverInfo" class="flex h-64 items-center justify-center">
      <el-icon class="is-loading" :size="40" color="#4096ff">
        <Loading />
      </el-icon>
    </div>

    <template v-else-if="serverInfo">
      <!-- 第一行：CPU + 内存 -->
      <div class="mb-4 grid grid-cols-1 gap-4 lg:grid-cols-2">
        <!-- CPU 信息 -->
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon class="text-blue-500"><Cpu /></el-icon>
              <span>{{ t('monitor.cpu') }}</span>
            </div>
          </template>
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
            <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
              <div
                v-for="(usage, index) in serverInfo.cpu.usage"
                :key="index"
                class="rounded p-2 text-center"
                style="background: var(--el-fill-color-light)"
              >
                <div class="text-xs text-gray-500 dark:text-gray-400">Core {{ index }}</div>
                <div class="text-sm font-medium" :class="getUsageColor(usage)">
                  {{ usage.toFixed(1) }}%
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 内存信息 -->
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon class="text-green-500"><Coin /></el-icon>
              <span>{{ t('monitor.memory') }}</span>
            </div>
          </template>
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
        </el-card>
      </div>

      <!-- 第二行：磁盘信息 -->
      <el-card shadow="never" class="mb-4">
        <template #header>
          <div class="card-header">
            <el-icon class="text-orange-500"><Box /></el-icon>
            <span>{{ t('monitor.disk') }}</span>
          </div>
        </template>
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
      </el-card>

      <!-- 第三行：Go 运行时 + 服务状态 -->
      <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <!-- Go 运行时信息 -->
        <el-card shadow="never" class="lg:col-span-2">
          <template #header>
            <div class="card-header">
              <el-icon class="text-cyan-500"><Monitor /></el-icon>
              <span>{{ t('monitor.goRuntime') }}</span>
            </div>
          </template>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.goVersion') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.goVersion }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.os') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.os }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.arch') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.arch }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.goroutines') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.goroutines }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.heapAlloc') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.heapAlloc }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.heapSys') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.heapSys }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.heapIdle') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.heapIdle }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.heapInuse') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.heapInuse }} {{ t('monitor.unitMB') }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.numGC') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.numGC }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400">{{ t('monitor.lastGC') }}</span>
              <span class="font-medium" style="color: var(--el-text-color-primary)">{{ serverInfo.goRuntime.lastGC }}</span>
            </div>
          </div>
        </el-card>

        <!-- 服务状态 -->
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ t('monitor.serviceStatus') }}</span>
            </div>
          </template>
          <div class="space-y-4">
            <!-- 数据库状态 -->
            <div class="rounded-lg p-4" style="border: 1px solid var(--el-border-color-light)">
              <div class="mb-2 flex items-center justify-between">
                <span class="font-medium" style="color: var(--el-text-color-primary)">{{ t('monitor.db') }}</span>
                <el-tag :type="serverInfo.db.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ serverInfo.db.status === 'online' ? t('monitor.online') : t('monitor.offline') }}
                </el-tag>
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ serverInfo.db.message }}</p>
            </div>
            <!-- Redis 状态 -->
            <div class="rounded-lg p-4" style="border: 1px solid var(--el-border-color-light)">
              <div class="mb-2 flex items-center justify-between">
                <span class="font-medium" style="color: var(--el-text-color-primary)">{{ t('monitor.redis') }}</span>
                <el-tag :type="serverInfo.redis.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ serverInfo.redis.status === 'online' ? t('monitor.online') : t('monitor.offline') }}
                </el-tag>
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ serverInfo.redis.message }}</p>
            </div>
          </div>
        </el-card>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, Loading, Cpu, Coin, Box, Monitor } from '@element-plus/icons-vue'
import { getServerMonitor } from '@/api/monitor'
import type { ServerMonitorInfo } from '@/api/monitor'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const loading = ref(false)
const serverInfo = ref<ServerMonitorInfo | null>(null)
const autoRefresh = ref(false)
let refreshTimer: ReturnType<typeof setInterval> | null = null

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
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
</style>
