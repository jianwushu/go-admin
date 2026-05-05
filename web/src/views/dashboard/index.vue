<template>
  <div class="dashboard-container">
    <!-- 欢迎横幅 -->
    <div class="welcome-banner">
      <div class="welcome-content">
        <div class="welcome-text">
          <h1 class="welcome-title">{{ t('dashboard.welcome') }}</h1>
          <p class="welcome-desc">{{ t('dashboard.welcomeDesc') }}</p>
        </div>
        <div class="welcome-time">
          <el-text size="large" class="welcome-time-text">{{ currentTime }}</el-text>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stat-grid section-spacing">
      <div
        v-for="item in statCards"
        :key="item.key"
        class="stat-item"
      >
        <div class="stat-item-inner">
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
      </div>
    </div>

    <!-- 中间区域：快捷入口 + 最近活动 -->
    <div class="content-grid section-spacing">
      <!-- 快捷入口 -->
      <div class="content-section">
        <div class="section-header">
          <el-icon :size="18"><Grid /></el-icon>
          <span>{{ t('dashboard.quickEntry') }}</span>
        </div>
        <div class="entry-grid">
          <router-link
            v-for="entry in quickEntries"
            :key="entry.path"
            :to="entry.path"
            class="quick-entry-item"
          >
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
        </div>
      </div>

      <!-- 最近活动 -->
      <div class="content-section">
        <div class="section-header">
          <el-icon :size="18"><Clock /></el-icon>
          <span>{{ t('dashboard.recentActivity') }}</span>
        </div>
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
      </div>
    </div>

    <!-- 底部区域：系统信息 -->
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, markRaw } from 'vue'
import { useI18n } from 'vue-i18n'
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
  todayLogin:0
})

// 统计卡片配置
const statCards = [
  {
    key: 'userCount',
    icon: markRaw(User),
    color: '#4096ff',
    bgColor: 'rgba(64, 150, 255, 0.06)',
  },
  {
    key: 'roleCount',
    icon: markRaw(UserFilled),
    color: '#52c41a',
    bgColor: 'rgba(82, 196, 26, 0.06)',
  },
  {
    key: 'menuCount',
    icon: markRaw(Menu),
    color: '#faad14',
    bgColor: 'rgba(250, 173, 20, 0.06)',
  },
  {
    key: 'deptCount',
    icon: markRaw(OfficeBuilding),
    color: '#ff4d4f',
    bgColor: 'rgba(255, 77, 79, 0.06)',
  },
  {
    key: 'todayLogin',
    icon: markRaw(Connection),
    color: '#722ed1',
    bgColor: 'rgba(114, 46, 209, 0.06)',
  },
] as const

// 快捷入口
const quickEntries = [
  {
    path: '/system/user',
    icon: markRaw(User),
    label: 'userManagement',
    iconColor: '#4096ff',
    bgColor: 'rgba(64, 150, 255, 0.08)',
  },
  {
    path: '/system/role',
    icon: markRaw(UserFilled),
    label: 'roleManagement',
    iconColor: '#52c41a',
    bgColor: 'rgba(82, 196, 26, 0.08)',
  },
  {
    path: '/system/menu',
    icon: markRaw(Menu),
    label: 'menuManagement',
    iconColor: '#faad14',
    bgColor: 'rgba(250, 173, 20, 0.08)',
  },
  {
    path: '/system/dept',
    icon: markRaw(OfficeBuilding),
    label: 'deptManagement',
    iconColor: '#ff4d4f',
    bgColor: 'rgba(255, 77, 79, 0.08)',
  },
  {
    path: '/monitor/server',
    icon: markRaw(Cpu),
    label: 'systemMonitor',
    iconColor: '#13c2c2',
    bgColor: 'rgba(19, 194, 194, 0.08)',
  },
  {
    path: '/codegen',
    icon: markRaw(Platform),
    label: 'codeGenerator',
    iconColor: '#722ed1',
    bgColor: 'rgba(114, 46, 209, 0.08)',
  },
]

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    content: '用户 admin 登录了系统',
    time: '2024-01-15 10:30:00',
    color: '#4096ff',
  },
  {
    id: 2,
    content: '角色 管理员 被修改',
    time: '2024-01-15 09:15:00',
    color: '#faad14',
  },
  {
    id: 3,
    content: '新用户 testuser 被创建',
    time: '2024-01-14 16:45:00',
    color: '#52c41a',
  },
  {
    id: 4,
    content: '菜单 系统管理 被更新',
    time: '2024-01-14 14:20:00',
    color: '#13c2c2',
  },
])

// 系统信息
const systemInfoItems: Array<{
  label: string
  icon: ReturnType<typeof markRaw>
  value: string
  tagType: 'success' | 'primary' | 'warning' | 'danger' | 'info'
}> = [
  { label: 'goVersion', icon: markRaw(Document), value: 'Go 1.21', tagType: 'success' },
  { label: 'framework', icon: markRaw(Setting), value: 'Gin', tagType: 'primary' },
  { label: 'database', icon: markRaw(Monitor), value: 'MySQL 8.0', tagType: 'warning' },
  { label: 'cache', icon: markRaw(Connection), value: 'Redis 7.0', tagType: 'danger' },
  { label: 'os', icon: markRaw(Cpu), value: 'Linux', tagType: 'info' },
  { label: 'arch', icon: markRaw(Platform), value: 'amd64', tagType: 'info' },
]

// 获取统计数据
async function fetchStats() {
  try {
    const { data: res } = await getDashboardStats()
    if (res.data) {
      stats.value = res.data
    }
  } catch {
    // 错误已在 request 拦截器中处理
  }
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
  fetchStats()
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.dashboard-container {
  max-width: 1400px;
}

/* 统一间距 */
.section-spacing {
  margin-bottom: 24px;
}

/* 欢迎横幅 - 白天模式 */
.welcome-banner {
  background: linear-gradient(135deg, #4096ff 0%, #1677ff 100%);
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 24px;
  position: relative;
  overflow: hidden;
}

.welcome-banner::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, transparent 70%);
  border-radius: 50%;
}

.welcome-banner::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: 10%;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
}

/* 欢迎横幅 - 暗黑模式 */
html.dark .welcome-banner {
  background: linear-gradient(135deg, #173066 0%, #0f1d3d 100%);
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
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.welcome-desc {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.95);
  margin: 0;
  text-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
}

.welcome-time {
  text-align: right;
}

.welcome-time-text {
  color: rgba(255, 255, 255, 0.9) !important;
  font-size: 14px;
  text-shadow: 0 1px 1px rgba(0, 0, 0, 0.1);
}

/* 统计网格 */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}

@media (max-width: 1200px) {
  .stat-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.stat-item {
  background: var(--el-bg-color);
  border-radius: 10px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s ease;
  cursor: pointer;
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

.stat-footer {
  padding: 10px 20px;
  border-top: 1px solid var(--el-border-color-extra-light);
  text-align: center;
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

/* 快捷入口 */
.entry-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

@media (max-width: 768px) {
  .entry-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

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

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
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

  .welcome-banner {
    padding: 24px 20px;
  }
}
</style>
