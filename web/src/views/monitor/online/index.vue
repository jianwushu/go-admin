<template>
  <div class="online-user-container">
    <!-- 页面标题 -->
    <div class="page-header section-spacing">
      <div class="page-header-content">
        <h1 class="page-title">{{ t('monitor.onlineUserTitle') }}</h1>
      </div>
    </div>

    <!-- 查询栏 -->
    <div class="content-section section-spacing">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('monitor.username')">
          <el-input
            v-model="queryParams.username"
            :placeholder="t('monitor.usernamePlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon class="mr-1"><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon class="mr-1"><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="content-section">
      <!-- 工具栏 -->
      <div class="table-toolbar">
        <div></div>
        <el-button @click="fetchData">
          <el-icon class="mr-1"><Refresh /></el-icon>
          {{ t('monitor.refresh') }}
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%; margin-bottom: 16px;"
      >
        <el-table-column prop="userId" label="ID" width="80" align="center" />
        <el-table-column prop="username" :label="t('monitor.username')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="nickname" :label="t('monitor.nickname')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="deptName" :label="t('monitor.deptName')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="ip" :label="t('monitor.ip')" width="140" show-overflow-tooltip />
        <el-table-column :label="t('monitor.loginTime')" width="170" align="center">
          <template #default="{ row }">
            {{ row.loginTime ? formatTimestamp(row.loginTime) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="onlineDuration" :label="t('monitor.onlineDuration')" width="120" align="center">
          <template #default="{ row }">
            <el-tag type="success" size="small">{{ row.onlineDuration || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.operation')" width="120" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              link
              v-permission="'monitor:online:forceLogout'"
              @click="handleForceLogout(row)"
            >
              {{ t('monitor.forceLogout') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 空数据提示 -->
      <el-empty
        v-if="!loading && tableData.length === 0"
        :description="t('monitor.noOnlineUser')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { getOnlineUsers, forceLogoutUser } from '@/api/monitor'
import type { OnlineUser } from '@/api/monitor'

defineOptions({ name: 'OnlineUser' })

const { t } = useI18n()

// 查询参数
const queryParams = reactive({
  username: '',
})

// 表格数据
const loading = ref(false)
const tableData = ref<OnlineUser[]>([])

// 自动刷新定时器
let refreshTimer: ReturnType<typeof setInterval> | null = null

/** 格式化时间戳 */
function formatTimestamp(timestamp: number): string {
  const date = new Date(timestamp * 1000)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getOnlineUsers(queryParams.username || undefined)
    tableData.value = res.data || []
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

/** 搜索 */
function handleSearch() {
  fetchData()
}

/** 重置 */
function handleReset() {
  queryParams.username = ''
  fetchData()
}

/** 强制下线 */
async function handleForceLogout(row: OnlineUser) {
  try {
    await ElMessageBox.confirm(t('monitor.forceLogoutConfirm'), t('common.tip'), {
      type: 'warning',
    })
    await forceLogoutUser(row.userId)
    ElMessage.success(t('monitor.forceLogoutSuccess'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

/** 启动自动刷新 */
function startAutoRefresh() {
  refreshTimer = setInterval(() => {
    fetchData()
  }, 30000) // 每30秒刷新一次
}

/** 停止自动刷新 */
function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  fetchData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.online-user-container {
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

/* 内容区域 */
.content-section {
  background: var(--el-bg-color);
  border-radius: 10px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
  padding: 20px;
}

/* 查询栏表单垂直居中 */
.content-section :deep(.el-form--inline) {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

/* 表格工具栏 */
.table-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
