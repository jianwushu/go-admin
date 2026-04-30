<template>
  <div class="p-4">
    <!-- 搜索区域 -->
    <div class="search-section mb-4">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('operationLog.module')">
          <el-input
            v-model="queryParams.module"
            :placeholder="t('operationLog.modulePlaceholder')"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('operationLog.operator')">
          <el-input
            v-model="queryParams.operator"
            :placeholder="t('operationLog.operatorPlaceholder')"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('operationLog.method')">
          <el-select v-model="queryParams.method" :placeholder="t('operationLog.allMethods')" clearable style="width: 120px">
            <el-option :label="t('operationLog.allMethods')" value="" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('operationLog.status')">
          <el-select v-model="queryParams.status" :placeholder="t('operationLog.status')" clearable style="width: 120px">
            <el-option :label="t('operationLog.success')" :value="1" />
            <el-option :label="t('operationLog.fail')" :value="0" />
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

    <!-- 数据表格 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ t('operationLog.title') }}</span>
          <el-button type="danger" v-permission="'monitor:operationLog:remove'" @click="handleClear">
            <el-icon class="mr-1"><Delete /></el-icon>
            {{ t('operationLog.clear') }}
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="module" :label="t('operationLog.module')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="action" :label="t('operationLog.action')" min-width="80" show-overflow-tooltip />
        <el-table-column :label="t('operationLog.method')" width="100" align="center">
          <template #default="{ row }">
            <el-tag
              :type="getMethodTagType(row.method)"
              size="small"
            >
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" :label="t('operationLog.url')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="operator" :label="t('operationLog.operator')" width="100" show-overflow-tooltip />
        <el-table-column prop="ip" :label="t('operationLog.ip')" width="140" show-overflow-tooltip />
        <el-table-column :label="t('operationLog.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? t('operationLog.success') : t('operationLog.fail') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('operationLog.duration')" width="90" align="center">
          <template #default="{ row }">
            <span :class="getDurationColor(row.duration)">{{ row.duration }} {{ t('operationLog.durationUnit') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('operationLog.createdAt')" width="170" align="center" />
        <el-table-column :label="t('common.operation')" width="80" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">
              {{ t('operationLog.detail') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="t('operationLog.detail')"
      width="700px"
      destroy-on-close
    >
      <template v-if="currentLog">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="t('operationLog.module')">{{ currentLog.module }}</el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.action')">{{ currentLog.action }}</el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.method')">
            <el-tag :type="getMethodTagType(currentLog.method)" size="small">{{ currentLog.method }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.status')">
            <el-tag :type="currentLog.status === 1 ? 'success' : 'danger'" size="small">
              {{ currentLog.status === 1 ? t('operationLog.success') : t('operationLog.fail') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.url')" :span="2">{{ currentLog.url }}</el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.operator')">{{ currentLog.operator }}</el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.ip')">{{ currentLog.ip }}</el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.duration')">
            <span :class="getDurationColor(currentLog.duration)">{{ currentLog.duration }} {{ t('operationLog.durationUnit') }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="t('operationLog.createdAt')">{{ currentLog.createdAt }}</el-descriptions-item>
        </el-descriptions>

        <div class="mt-4">
          <h4 class="mb-2 font-medium text-gray-700 dark:text-gray-300">{{ t('operationLog.requestParam') }}</h4>
          <el-input
            type="textarea"
            :model-value="formatJson(currentLog.requestParam)"
            :rows="4"
            readonly
          />
        </div>

        <div class="mt-4">
          <h4 class="mb-2 font-medium text-gray-700 dark:text-gray-300">{{ t('operationLog.responseData') }}</h4>
          <el-input
            type="textarea"
            :model-value="formatJson(currentLog.responseData)"
            :rows="4"
            readonly
          />
        </div>

        <div v-if="currentLog.errorMsg" class="mt-4">
          <h4 class="mb-2 font-medium text-red-500">{{ t('operationLog.errorMsg') }}</h4>
          <el-input
            type="textarea"
            :model-value="currentLog.errorMsg"
            :rows="3"
            readonly
          />
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import { getOperationLogList, clearOperationLog } from '@/api/operationLog'
import type { OperationLogItem, OperationLogListParams } from '@/api/operationLog'

defineOptions({ name: 'OperationLog' })

const { t } = useI18n()

// 查询参数
const queryParams = reactive<OperationLogListParams>({
  page: 1,
  pageSize: 10,
  module: undefined,
  operator: undefined,
  method: undefined,
  status: undefined,
})

// 表格数据
const loading = ref(false)
const tableData = ref<OperationLogItem[]>([])
const total = ref(0)

// 详情弹窗
const detailVisible = ref(false)
const currentLog = ref<OperationLogItem | null>(null)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getOperationLogList(queryParams)
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
  queryParams.module = undefined
  queryParams.operator = undefined
  queryParams.method = undefined
  queryParams.status = undefined
  queryParams.page = 1
  fetchData()
}

/** 分页大小变化 */
function handleSizeChange() {
  queryParams.page = 1
  fetchData()
}

/** 页码变化 */
function handleCurrentChange() {
  fetchData()
}

/** 查看详情 */
function handleDetail(row: OperationLogItem) {
  currentLog.value = row
  detailVisible.value = true
}

/** 清空日志 */
async function handleClear() {
  try {
    await ElMessageBox.confirm(t('operationLog.clearConfirm'), t('common.tip'), {
      type: 'warning',
    })
    await clearOperationLog()
    ElMessage.success(t('common.success'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

/** 获取方法标签类型 */
function getMethodTagType(method: string): 'success' | 'warning' | 'danger' | 'info' {
  switch (method) {
    case 'POST': return 'success'
    case 'PUT': return 'warning'
    case 'DELETE': return 'danger'
    default: return 'info'
  }
}

/** 获取耗时颜色 */
function getDurationColor(duration: number): string {
  if (duration >= 3000) return 'text-red-500'
  if (duration >= 1000) return 'text-orange-500'
  return 'text-green-500'
}

/** 格式化 JSON */
function formatJson(str: string): string {
  if (!str) return ''
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

onMounted(() => {
  fetchData()
})
</script>
