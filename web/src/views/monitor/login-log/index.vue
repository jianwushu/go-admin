<template>
  <div class="p-4">
    <!-- 搜索区域 -->
    <el-card shadow="never" class="mb-4">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('loginLog.username')">
          <el-input
            v-model="queryParams.username"
            :placeholder="t('loginLog.usernamePlaceholder')"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('loginLog.ip')">
          <el-input
            v-model="queryParams.ip"
            :placeholder="t('loginLog.ipPlaceholder')"
            clearable
            style="width: 180px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('loginLog.status')">
          <el-select v-model="queryParams.status" :placeholder="t('loginLog.status')" clearable style="width: 120px">
            <el-option :label="t('loginLog.success')" :value="1" />
            <el-option :label="t('loginLog.fail')" :value="0" />
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
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ t('loginLog.title') }}</span>
          <el-button type="danger" v-permission="'monitor:loginLog:remove'" @click="handleClear">
            <el-icon class="mr-1"><Delete /></el-icon>
            {{ t('loginLog.clear') }}
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
        <el-table-column prop="username" :label="t('loginLog.username')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="ip" :label="t('loginLog.ip')" width="140" show-overflow-tooltip />
        <el-table-column prop="location" :label="t('loginLog.location')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="browser" :label="t('loginLog.browser')" width="120" show-overflow-tooltip />
        <el-table-column prop="os" :label="t('loginLog.os')" width="120" show-overflow-tooltip />
        <el-table-column :label="t('loginLog.status')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? t('loginLog.success') : t('loginLog.fail') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="msg" :label="t('loginLog.msg')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="t('loginLog.createdAt')" width="170" align="center" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import { getLoginLogList, clearLoginLog } from '@/api/loginLog'
import type { LoginLogItem, LoginLogListParams } from '@/api/loginLog'

defineOptions({ name: 'LoginLog' })

const { t } = useI18n()

// 查询参数
const queryParams = reactive<LoginLogListParams>({
  page: 1,
  pageSize: 10,
  username: undefined,
  ip: undefined,
  status: undefined,
})

// 表格数据
const loading = ref(false)
const tableData = ref<LoginLogItem[]>([])
const total = ref(0)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getLoginLogList(queryParams)
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
  queryParams.username = undefined
  queryParams.ip = undefined
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

/** 清空日志 */
async function handleClear() {
  try {
    await ElMessageBox.confirm(t('loginLog.clearConfirm'), t('common.tip'), {
      type: 'warning',
    })
    await clearLoginLog()
    ElMessage.success(t('common.success'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

onMounted(() => {
  fetchData()
})
</script>
