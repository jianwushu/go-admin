<template>
  <div class="job-log-container">
    <!-- 页面标题 -->
    <div class="page-header section-spacing">
      <div class="page-header-content">
        <h1 class="page-title">{{ t('job.logTitle') }}</h1>
      </div>
    </div>

    <!-- 查询栏 -->
    <div class="content-section section-spacing">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('job.logJobName')">
          <el-input
            v-model="jobNameFilter"
            :placeholder="t('job.namePlaceholder')"
            clearable
            style="width: 180px"
            disabled
          />
        </el-form-item>
        <el-form-item :label="t('job.logStatus')">
          <el-select v-model="queryParams.status" :placeholder="t('job.logStatus')" clearable style="width: 120px">
            <el-option :label="t('job.logStatusSuccess')" :value="1" />
            <el-option :label="t('job.logStatusFail')" :value="2" />
          </el-select>
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
        <el-button type="danger" v-permission="'tool:job:delete'" :disabled="!queryParams.jobId" @click="handleCleanLogs">
          <el-icon class="mr-1"><Delete /></el-icon>
          {{ t('job.cleanLog') }}
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%; margin-bottom: 16px;"
      >
        <el-table-column prop="jobName" :label="t('job.logJobName')" min-width="140" show-overflow-tooltip />
        <el-table-column :label="t('job.logStatus')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? t('job.logStatusSuccess') : t('job.logStatusFail') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('job.logDuration')" width="110" align="center">
          <template #default="{ row }">
            <span>{{ row.duration }} ms</span>
          </template>
        </el-table-column>
        <el-table-column prop="result" :label="t('job.logResult')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="errorMsg" :label="t('job.logErrorMsg')" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.errorMsg" class="text-red-500">{{ row.errorMsg }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('job.logTime')" width="170" align="center" />
      </el-table>

      <!-- 分页 -->
      <Pagination
        v-model:page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        @change="fetchData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import { getJobLogList, cleanJobLogs, getJobById } from '@/api/job'
import type { JobLogItem, JobLogListParams } from '@/types/api'
import Pagination from '@/components/Pagination.vue'

defineOptions({ name: 'JobLog' })

const { t } = useI18n()
const route = useRoute()

// 查询参数
const queryParams = reactive<JobLogListParams>({
  page: 1,
  pageSize: 10,
  jobId: undefined,
  status: undefined,
})

// 任务名称（用于显示）
const jobNameFilter = ref('')

// 表格数据
const loading = ref(false)
const tableData = ref<JobLogItem[]>([])
const total = ref(0)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getJobLogList(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

/** 搜索 */
function handleSearch() {
  queryParams.page = 1
  fetchData()
}

/** 重置 */
function handleReset() {
  queryParams.status = undefined
  queryParams.page = 1
  fetchData()
}

/** 清空日志 */
async function handleCleanLogs() {
  if (!queryParams.jobId) return
  try {
    await ElMessageBox.confirm(t('job.cleanLogConfirm'), t('common.tip'), { type: 'warning' })
    await cleanJobLogs(queryParams.jobId)
    ElMessage.success(t('job.cleanLogSuccess'))
    fetchData()
  } catch {
    // 取消
  }
}

/** 初始化：从路由参数获取 jobId */
async function init() {
  const jobId = route.query.jobId
  if (jobId) {
    queryParams.jobId = Number(jobId)
    // 获取任务名称
    try {
      const { data: res } = await getJobById(Number(jobId))
      if (res.code === 0 && res.data) {
        jobNameFilter.value = (res.data as any).name || ''
      }
    } catch {
      // ignore
    }
  }
  fetchData()
}

onMounted(() => {
  init()
})
</script>

<style scope>
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

/* 表格工具栏 */
.table-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
