<template>
  <div class="job-container">
    <!-- 页面标题 -->
    <div class="page-header section-spacing">
      <div class="page-header-content">
        <h1 class="page-title">{{ t('job.title') }}</h1>
      </div>
    </div>

    <!-- 查询栏 -->
    <div class="content-section section-spacing">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('job.name')">
          <el-input
            v-model="queryParams.name"
            :placeholder="t('job.namePlaceholder')"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('job.jobType')">
          <el-select v-model="queryParams.jobType" :placeholder="t('job.jobType')" clearable style="width: 140px">
            <el-option :label="t('job.jobTypeFunc')" :value="1" />
            <el-option :label="t('job.jobTypeHTTP')" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('job.status')">
          <el-select v-model="queryParams.status" :placeholder="t('job.status')" clearable style="width: 120px">
            <el-option :label="t('job.enabled')" :value="1" />
            <el-option :label="t('job.disabled')" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon class="mr-1"><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
        </el-form-item>
        <el-form-item>
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
        <div>
          <el-button type="primary" v-permission="'tool:job:add'" @click="handleAdd">
            <el-icon class="mr-1"><Plus /></el-icon>
            {{ t('common.add') }}
          </el-button>
          <el-button type="danger" v-permission="'tool:job:delete'" :disabled="!selectedIds.length" @click="handleBatchDelete">
            <el-icon class="mr-1"><Delete /></el-icon>
            {{ t('common.delete') }}
          </el-button>
        </div>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%; margin-bottom: 16px;"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column prop="name" :label="t('job.name')" min-width="140" show-overflow-tooltip />
        <el-table-column :label="t('job.jobType')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.jobType === 1 ? 'success' : 'warning'" size="small">
              {{ row.jobType === 1 ? t('job.jobTypeFunc') : t('job.jobTypeHTTP') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cronExpr" :label="t('job.cronExpr')" width="160" show-overflow-tooltip />
        <el-table-column :label="t('job.funcName')" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.jobType === 1">{{ row.funcName }}</span>
            <span v-else>{{ row.httpMethod }} {{ row.httpUrl }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('job.status')" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              v-permission="'tool:job:edit'"
              @change="handleStatusChange(row, $event)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="t('job.remark')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="t('common.createdAt')" width="170" align="center" />
        <el-table-column :label="t('common.operation')" width="240" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link v-permission="'tool:job:edit'" @click="handleEdit(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button type="success" link v-permission="'tool:job:edit'" @click="handleRunOnce(row)">
              {{ t('job.runOnce') }}
            </el-button>
            <el-button type="info" link v-permission="'tool:job:list'" @click="handleViewLog(row)">
              {{ t('job.viewLog') }}
            </el-button>
            <el-button type="danger" link v-permission="'tool:job:delete'" @click="handleDelete(row)">
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <Pagination
        v-model:page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        @change="fetchData"
      />
    </div>

    <!-- 任务表单弹窗 -->
    <JobForm
      v-model:visible="formVisible"
      :data="currentJob"
      @success="fetchData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Delete } from '@element-plus/icons-vue'
import { getJobList, deleteJob, changeJobStatus, runJobOnce } from '@/api/job'
import type { JobItem, JobListParams } from '@/types/api'
import Pagination from '@/components/Pagination.vue'
import JobForm from './components/JobForm.vue'

defineOptions({ name: 'JobIndex' })

const { t } = useI18n()
const router = useRouter()

// 查询参数
const queryParams = reactive<JobListParams>({
  page: 1,
  pageSize: 10,
  name: undefined,
  jobType: undefined,
  status: undefined,
})

// 表格数据
const loading = ref(false)
const tableData = ref<JobItem[]>([])
const total = ref(0)

// 选中的ID
const selectedIds = ref<number[]>([])

// 表单弹窗
const formVisible = ref(false)
const currentJob = ref<JobItem | null>(null)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getJobList(queryParams)
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
  queryParams.name = undefined
  queryParams.jobType = undefined
  queryParams.status = undefined
  queryParams.page = 1
  fetchData()
}

/** 选择变化 */
function handleSelectionChange(rows: JobItem[]) {
  selectedIds.value = rows.map(r => r.id)
}

/** 新增 */
function handleAdd() {
  currentJob.value = null
  formVisible.value = true
}

/** 编辑 */
function handleEdit(row: JobItem) {
  currentJob.value = { ...row }
  formVisible.value = true
}

/** 删除 */
async function handleDelete(row: JobItem) {
  try {
    await ElMessageBox.confirm(t('job.confirmDelete'), t('common.tip'), { type: 'warning' })
    await deleteJob(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchData()
  } catch {
    // 取消
  }
}

/** 批量删除 */
async function handleBatchDelete() {
  try {
    await ElMessageBox.confirm(t('job.confirmBatchDelete'), t('common.tip'), { type: 'warning' })
    await deleteJob(selectedIds.value)
    ElMessage.success(t('job.batchDeleteSuccess'))
    fetchData()
  } catch {
    // 取消
  }
}

/** 修改状态 */
async function handleStatusChange(row: JobItem, val: unknown) {
  const status = Number(val)
  try {
    await changeJobStatus(row.id, status)
    ElMessage.success(t('common.success'))
  } catch {
    row.status = status === 1 ? 2 : 1 // 回滚
  }
}

/** 手动执行一次 */
async function handleRunOnce(row: JobItem) {
  try {
    await ElMessageBox.confirm(t('job.runOnceConfirm'), t('common.tip'), { type: 'info' })
    await runJobOnce(row.id)
    ElMessage.success(t('job.runOnceSuccess'))
  } catch {
    // 取消
  }
}

/** 查看日志 */
function handleViewLog(row: JobItem) {
  router.push({ path: '/tool/job-log', query: { jobId: row.id } })
}

onMounted(() => {
  fetchData()
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