<template>
  <div class="dashboard-container">
    <!-- 欢迎信息 -->
    <el-card
      :body-style="{ padding: '32px 40px' }"
      class="welcome-card mb-8"
      shadow="never"
    >
      <div class="welcome-content">
        <div class="welcome-text">
          <h1 class="welcome-title">{{ t('dashboard.welcome') }}</h1>
          <p class="welcome-desc">{{ t('dashboard.welcomeDesc') }}</p>
        </div>
        <div class="welcome-time">
          <el-text type="info" size="large">{{ currentTime }}</el-text>
        </div>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="mb-8">
      <el-col
        v-for="item in statCards"
        :key="item.key"
        :xs="12"
        :sm="12"
        :md="8"
        :lg="4"
        :xl="4"
        class="mb-4 lg:mb-0"
      >
        <el-card
          shadow="hover"
          class="stat-card"
          :body-style="{ padding: '24px' }"
        >
          <div class="stat-card-inner">
            <el-statistic :value="stats[item.key] ?? 0">
              <template #title>
                <div class="stat-title">
                  <el-icon :size="16" :color="item.color">
                    <component :is="item.icon" />
                  </el-icon>
                  <span>{{ t(`dashboard.${item.key}`) }}</span>
                </div>
              </template>
            </el-statistic>
          </div>
          <div class="stat-footer" :style="{ background: item.bgColor }">
            <el-text size="small" :style="{ color: item.color }">
              {{ t('dashboard.viewDetails') }}
            </el-text>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 中间区域：快捷入口 + 最近活动 -->
    <el-row :gutter="20" class="mb-8">
      <el-col :xs="24" :lg="12" class="mb-4 lg:mb-0">
        <el-card shadow="never" class="h-full">
          <template #header>
            <div class="card-header">
              <el-icon :size="18"><Grid /></el-icon>
              <span>{{ t('dashboard.quickEntry') }}</span>
            </div>
          </template>
          <el-row :gutter="16">
            <el-col
              v-for="entry in quickEntries"
              :key="entry.path"
              :span="8"
              class="mb-4"
            >
              <router-link :to="entry.path" class="quick-entry-item">
                <div
                  class="quick-entry-icon"
                  :style="{ background: entry.bgColor }"
                >
                  <el-icon :size="28" :color="entry.iconColor">
                    <component :is="entry.icon" />
                  </el-icon>
                </div>
                <span class="quick-entry-label">
                  {{ t(`dashboard.${entry.label}`) }}
                </span>
              </router-link>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card shadow="never" class="h-full">
          <template #header>
            <div class="card-header">
              <el-icon :size="18"><Clock /></el-icon>
              <span>{{ t('dashboard.recentActivity') }}</span>
            </div>
          </template>
          <el-empty
            v-if="recentActivities.length === 0"
            :description="t('dashboard.noActivity')"
            :image-size="80"
          />
          <el-timeline v-else>
            <el-timeline-item
              v-for="activity in recentActivities"
              :key="activity.id"
              :timestamp="activity.time"
              placement="top"
              :color="activity.color"
            >
              <el-text>{{ activity.content }}</el-text>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>

    <!-- 底部区域：系统信息 -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon :size="18"><InfoFilled /></el-icon>
              <span>{{ t('dashboard.systemInfo') }}</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item
              v-for="info in systemInfoItems"
              :key="info.label"
              :label="t(`dashboard.${info.label}`)"
              :width="200"
            >
              <template #label>
                <div class="desc-label">
                  <el-icon :size="14"><component :is="info.icon" /></el-icon>
                  <span>{{ t(`dashboard.${info.label}`) }}</span>
                </div>
              </template>
              <el-tag :type="info.tagType" effect="plain" size="small">
                {{ info.value }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/store/modules/app'
import {
  User,
  UserFilled,
  Menu,
  OfficeBuilding,
  Monitor,
  Connection,
  Setting,
  Platform,
  Cpu,
  Document,
  Grid,
  Clock,
  InfoFilled,
} from '@element-plus/icons-vue'
import { getDashboardStats, type DashboardStats } from '@/api/dashboard'

defineOptions({ name: 'Dashboard' })

const { t, locale } = useI18n()
const appStore = useAppStore()

// 当前时间
const currentTime = ref('')
let timer: ReturnType<typeof setInterval> | null = null

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 统计数据
const stats = ref<DashboardStats>({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  deptCount: 0,
  todayLogin: 0,
})

// 统计卡片配置
const statCards = [
  {
    key: 'userCount' as const,
    icon: markRaw(User),
    color: '#3b82f6',
    bgColor: 'rgba(59, 130, 246, 0.08)',
  },
  {
    key: 'roleCount' as const,
    icon: markRaw(UserFilled),
    color: '#22c55e',
    bgColor: 'rgba(34, 197, 94, 0.08)',
  },
  {
    key: 'menuCount' as const,
    icon: markRaw(Menu),
    color: '#f59e0b',
    bgColor: 'rgba(245, 158, 11, 0.08)',
  },
  {
    key: 'deptCount' as const,
    icon: markRaw(OfficeBuilding),
    color: '#64748b',
    bgColor: 'rgba(100, 116, 139, 0.08)',
  },
  {
    key: 'todayLogin' as const,
    icon: markRaw(Connection),
    color: '#ef4444',
    bgColor: 'rgba(239, 68, 68, 0.08)',
  },
]

// 快捷入口配置
const quickEntries = [
  {
    path: '/system/user',
    label: 'userManagement',
    icon: markRaw(User),
    bgColor: 'rgba(59, 130, 246, 0.1)',
    iconColor: '#3b82f6',
  },
  {
    path: '/system/role',
    label: 'roleManagement',
    icon: markRaw(UserFilled),
    bgColor: 'rgba(34, 197, 94, 0.1)',
    iconColor: '#22c55e',
  },
  {
    path: '/system/menu',
    label: 'menuManagement',
    icon: markRaw(Menu),
    bgColor: 'rgba(245, 158, 11, 0.1)',
    iconColor: '#f59e0b',
  },
  {
    path: '/system/dept',
    label: 'deptManagement',
    icon: markRaw(OfficeBuilding),
    bgColor: 'rgba(100, 116, 139, 0.1)',
    iconColor: '#64748b',
  },
  {
    path: '/monitor/server',
    label: 'systemMonitor',
    icon: markRaw(Monitor),
    bgColor: 'rgba(239, 68, 68, 0.1)',
    iconColor: '#ef4444',
  },
  {
    path: '/codegen',
    label: 'codeGenerator',
    icon: markRaw(Setting),
    bgColor: 'rgba(139, 92, 246, 0.1)',
    iconColor: '#8b5cf6',
  },
]

// 最近活动
const recentActivities = ref<
  Array<{ id: number; content: string; time: string; color: string }>
>([])

// 系统信息
const systemInfoItems = [
  {
    label: 'systemName',
    value: 'Go-Admin',
    icon: markRaw(Platform),
    tagType: 'primary' as const,
  },
  {
    label: 'systemVersion',
    value: 'v1.0.0',
    icon: markRaw(Document),
    tagType: 'success' as const,
  },
  {
    label: 'goVersion',
    value: 'Go 1.21+',
    icon: markRaw(Cpu),
    tagType: 'info' as const,
  },
  {
    label: 'vueVersion',
    value: 'Vue 3.4+',
    icon: markRaw(Monitor),
    tagType: 'warning' as const,
  },
]

// 获取统计数据
async function fetchStats() {
  try {
    const { data: res } = await getDashboardStats()
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('Failed to fetch dashboard stats:', error)
  }
}

onMounted(() => {
  fetchStats()
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.dashboard-container {
  min-height: calc(100vh - 84px);
  padding: 20px;
}

/* 欢迎卡片 */
.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
}

html.dark .welcome-card {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
}

.welcome-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

.welcome-card::after {
  content: '';
  position: absolute;
  bottom: -30%;
  right: 10%;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 1;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.welcome-desc {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
}

.welcome-time {
  text-align: right;
}

.welcome-time .el-text {
  color: rgba(255, 255, 255, 0.7) !important;
  font-size: 14px;
}

/* 统计卡片 */
.stat-card {
  border-radius: 10px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card-inner {
  padding-bottom: 16px;
}

.stat-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--el-text-color-regular);
}

.stat-footer {
  margin: 0 -24px -24px;
  padding: 12px 24px;
  border-top: 1px solid var(--el-border-color-extra-light);
  border-radius: 0 0 10px 10px;
  text-align: center;
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

/* 快捷入口 */
.quick-entry-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 12px;
  border-radius: 10px;
  text-decoration: none;
  transition: all 0.3s ease;
  cursor: pointer;
}

.quick-entry-item:hover {
  background: var(--el-fill-color-light);
  transform: translateY(-2px);
}

.quick-entry-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  transition: all 0.3s ease;
}

.quick-entry-item:hover .quick-entry-icon {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

html.dark .quick-entry-item:hover .quick-entry-icon {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.quick-entry-label {
  font-size: 13px;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

/* 系统信息描述 */
.desc-label {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 响应式 */
@media (max-width: 768px) {
  .welcome-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .welcome-title {
    font-size: 22px;
  }

  .welcome-time {
    text-align: left;
  }
}
</style>
