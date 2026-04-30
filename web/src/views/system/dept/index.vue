<template>
  <div class="p-4">
    <!-- 操作栏 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ t('dept.title') }}</span>
          <div>
            <el-button type="primary" v-permission="'system:dept:add'" @click="handleAdd()">
              <el-icon class="mr-1"><Plus /></el-icon>
              {{ t('common.add') }}
            </el-button>
            <el-button @click="toggleExpand">
              <el-icon class="mr-1"><Sort /></el-icon>
              {{ isExpand ? t('common.collapse') : t('common.expand') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 树形表格 -->
      <el-table
        v-if="refreshTable"
        v-loading="loading"
        :data="tableData"
        row-key="id"
        :default-expand-all="isExpand"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        border
        style="width: 100%"
      >
        <el-table-column prop="name" :label="t('dept.name')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="sort" :label="t('dept.sort')" width="80" align="center" />
        <el-table-column prop="leader" :label="t('dept.leader')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="phone" :label="t('dept.phone')" min-width="130" show-overflow-tooltip />
        <el-table-column prop="email" :label="t('dept.email')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="t('dept.status')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">{{ t('dept.enabled') }}</el-tag>
            <el-tag v-else type="danger">{{ t('dept.disabled') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.operation')" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link v-permission="'system:dept:add'" @click="handleAdd(row)">
              <el-icon class="mr-1"><Plus /></el-icon>
              {{ t('common.add') }}
            </el-button>
            <el-button type="primary" link v-permission="'system:dept:edit'" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="danger" link v-permission="'system:dept:remove'" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 部门表单弹窗 -->
    <DeptForm
      v-model:visible="formVisible"
      :data="currentDept"
      :is-edit="isEdit"
      :dept-tree="tableData"
      @success="fetchData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Sort } from '@element-plus/icons-vue'
import { getDeptTree, deleteDept } from '@/api/system'
import type { DeptInfo } from '@/types/api'
import DeptForm from './components/DeptForm.vue'

defineOptions({ name: 'DeptManagement' })

const { t } = useI18n()

// 表格数据
const loading = ref(false)
const tableData = ref<DeptInfo[]>([])
const isExpand = ref(true)
const refreshTable = ref(true)

// 表单相关
const formVisible = ref(false)
const isEdit = ref(false)
const currentDept = ref<DeptInfo | null>(null)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getDeptTree()
    tableData.value = res.data || []
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

/** 切换展开/折叠 */
async function toggleExpand() {
  refreshTable.value = false
  isExpand.value = !isExpand.value
  await nextTick()
  refreshTable.value = true
}

/** 新增 */
function handleAdd(row?: DeptInfo) {
  currentDept.value = row ? { ...row } as DeptInfo : null
  isEdit.value = false
  formVisible.value = true
}

/** 编辑 */
function handleEdit(row: DeptInfo) {
  currentDept.value = { ...row }
  isEdit.value = true
  formVisible.value = true
}

/** 删除 */
async function handleDelete(row: DeptInfo) {
  try {
    await ElMessageBox.confirm(t('common.confirmDelete'), t('common.tip'), {
      type: 'warning',
    })
    await deleteDept(row.id)
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
