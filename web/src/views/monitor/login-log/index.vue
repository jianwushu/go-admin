<template>
  <div style="padding: 16px;">
    <!-- 查询栏 -->
    <el-card shadow="never" style="margin-bottom: 16px;">
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
        <span style="font-weight: 500; font-size: 18px;">{{ t('loginLog.title') }}</span>
        <div style="display: flex; align-items: center; gap: 8px;">
          <el-button type="danger" v-permission="'monitor:loginLog:remove'" @click="handleClear">
            <el-icon style="margin-right: 4px;"><Delete /></el-icon>
            {{ t('loginLog.clear') }}
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
      <Pagination
        v-model:page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        @change="fetchData"
      />
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
import Pagination from '@/components/Pagination.vue'

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
