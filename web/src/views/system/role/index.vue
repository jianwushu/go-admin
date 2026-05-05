<template>
  <div style="padding: 16px;">
    <!-- 查询栏 -->
    <el-card shadow="never" style="margin-bottom: 16px;">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('role.name')">
          <el-input
            v-model="queryParams.name"
            :placeholder="t('role.namePlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('role.code')">
          <el-input
            v-model="queryParams.code"
            :placeholder="t('role.codePlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('role.status')">
          <el-select v-model="queryParams.status" :placeholder="t('common.status')" clearable style="width: 120px">
            <el-option :label="t('role.enabled')" :value="1" />
            <el-option :label="t('role.disabled')" :value="0" />
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
        <span style="font-weight: 500; font-size: 18px;">{{ t('role.title') }}</span>
        <div style="display: flex; align-items: center; gap: 8px;">
          <el-button type="primary" v-permission="'system:role:add'" @click="handleAdd">
            <el-icon style="margin-right: 4px;"><Plus /></el-icon>
            {{ t('common.add') }}
          </el-button>
          <el-button type="danger" v-permission="'system:role:remove'" :disabled="selectedIds.length === 0" @click="handleBatchDelete">
            <el-icon style="margin-right: 4px;"><Delete /></el-icon>
            {{ t('common.batchDelete') }}
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
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column prop="name" :label="t('role.name')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="code" :label="t('role.code')" min-width="120" show-overflow-tooltip />
        <el-table-column :label="t('role.dataScope')" min-width="140" align="center">
          <template #default="{ row }">
            {{ getDataScopeLabel(row.dataScope) }}
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="t('role.sort')" width="80" align="center" />
        <el-table-column :label="t('role.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              v-permission="'system:role:edit'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="t('role.remark')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="t('common.createTime')" width="170" align="center" />
        <el-table-column :label="t('common.operation')" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link v-permission="'system:role:edit'" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="danger" link v-permission="'system:role:remove'" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
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
    </el-card>

    <!-- 角色表单弹窗 -->
    <RoleForm
      v-model:visible="formVisible"
      :data="currentRole"
      :is-edit="isEdit"
      @success="fetchData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getRoleList, deleteRole, changeRoleStatus } from '@/api/system'
import type { RoleInfo, RoleListParams } from '@/types/api'
import RoleForm from './components/RoleForm.vue'
import Pagination from '@/components/Pagination.vue'

defineOptions({ name: 'RoleManagement' })

const { t } = useI18n()

// 查询参数
const queryParams = reactive<RoleListParams>({
  page: 1,
  pageSize: 10,
  name: undefined,
  code: undefined,
  status: undefined,
})

// 表格数据
const loading = ref(false)
const tableData = ref<RoleInfo[]>([])
const total = ref(0)

// 多选相关
const selectedIds = ref<number[]>([])

// 表单相关
const formVisible = ref(false)
const isEdit = ref(false)
const currentRole = ref<RoleInfo | null>(null)

/** 获取数据权限标签 */
function getDataScopeLabel(scope: number): string {
  const map: Record<number, string> = {
    1: t('role.dataScopeAll'),
    2: t('role.dataScopeDept'),
    3: t('role.dataScopeDeptAndChild'),
    4: t('role.dataScopeSelf'),
    5: t('role.dataScopeCustom'),
  }
  return map[scope] || ''
}

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getRoleList(queryParams)
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
  queryParams.code = undefined
  queryParams.status = undefined
  queryParams.page = 1
  fetchData()
}

/** 多选变化 */
function handleSelectionChange(rows: RoleInfo[]) {
  selectedIds.value = rows.map(row => row.id)
}

/** 新增 */
function handleAdd() {
  currentRole.value = null
  isEdit.value = false
  formVisible.value = true
}

/** 编辑 */
function handleEdit(row: RoleInfo) {
  currentRole.value = { ...row }
  isEdit.value = true
  formVisible.value = true
}

/** 删除 */
async function handleDelete(row: RoleInfo) {
  try {
    await ElMessageBox.confirm(t('common.confirmDelete'), t('common.tip'), {
      type: 'warning',
    })
    await deleteRole(row.id)
    ElMessage.success(t('common.success'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

/** 批量删除 */
async function handleBatchDelete() {
  if (selectedIds.value.length === 0) {
    ElMessage.warning(t('role.selectAtLeastOne'))
    return
  }
  try {
    await ElMessageBox.confirm(t('role.confirmBatchDelete'), t('common.tip'), {
      type: 'warning',
    })
    await deleteRole(selectedIds.value)
    ElMessage.success(t('role.batchDeleteSuccess'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

/** 状态变更 */
async function handleStatusChange(row: RoleInfo) {
  try {
    await changeRoleStatus(row.id, row.status)
    ElMessage.success(t('common.success'))
  } catch {
    // 回滚状态
    row.status = row.status === 1 ? 0 : 1
  }
}

onMounted(() => {
  fetchData()
})
</script>
