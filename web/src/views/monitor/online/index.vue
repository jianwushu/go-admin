<template>
  <div style="padding: 16px;">
    <!-- 查询栏 -->
    <el-card shadow="never" style="margin-bottom: 16px;">
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
            <el-icon style="margin-right: 4px;"><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon style="margin-right: 4px;"><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" body-style="padding: 20px;">
      <!-- 表格工具栏 -->
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
        <span style="font-weight: 500; font-size: 18px;">{{ t('monitor.onlineUserTitle') }}</span>
        <div style="display: flex; align-items: center; gap: 8px;">
          <el-button @click="fetchData">
            <el-icon style="margin-right: 4px;"><Refresh /></el-icon>
            {{ t('monitor.refresh') }}
          </el-button>
        </div>
      </div>

      <!-- 表格主体 -->
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
    </el-card>
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
