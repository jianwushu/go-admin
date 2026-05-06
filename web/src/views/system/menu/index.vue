<template>
  <div class="menu-container">
    <!-- 页面标题和操作栏 -->
    <div class="page-header section-spacing">
      <div class="page-header-content">
        <h1 class="page-title">{{ t('menu.title') }}</h1>
        <div class="page-actions">
          <el-button type="primary" v-permission="'system:menu:add'" @click="handleAdd()">
            <el-icon class="mr-1"><Plus /></el-icon>
            {{ t('common.add') }}
          </el-button>
          <el-button @click="toggleExpand">
            <el-icon class="mr-1"><Sort /></el-icon>
            {{ isExpand ? t('common.collapse') : t('common.expand') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 树形表格 -->
    <div class="content-section">
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
        <el-table-column prop="name" :label="t('menu.name')" min-width="180" show-overflow-tooltip />
        <el-table-column prop="i18nKey" :label="t('menu.i18nKey')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="icon" :label="t('menu.icon')" width="80" align="center">
          <template #default="{ row }">
            <el-icon v-if="getIcon(row.icon)"><component :is="getIcon(row.icon)" /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="t('menu.type')" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.type === 0" type="primary">{{ t('menu.typeDir') }}</el-tag>
            <el-tag v-else-if="row.type === 1" type="success">{{ t('menu.typeMenu') }}</el-tag>
            <el-tag v-else-if="row.type === 2" type="warning">{{ t('menu.typeButton') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" :label="t('menu.sort')" width="80" align="center" />
        <el-table-column prop="perms" :label="t('menu.perms')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="path" :label="t('menu.path')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="component" :label="t('menu.component')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="t('menu.visible')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.visible === 1" type="success">{{ t('menu.visibleShow') }}</el-tag>
            <el-tag v-else type="info">{{ t('menu.visibleHide') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('menu.status')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">{{ t('menu.enabled') }}</el-tag>
            <el-tag v-else type="danger">{{ t('menu.disabled') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.operation')" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.type !== 2" type="primary" link v-permission="'system:menu:add'" @click="handleAdd(row)">
              <el-icon class="mr-1"><Plus /></el-icon>
              {{ t('common.add') }}
            </el-button>
            <el-button type="primary" link v-permission="'system:menu:edit'" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="danger" link v-permission="'system:menu:remove'" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 菜单表单弹窗 -->
    <MenuForm
      v-model:visible="formVisible"
      :data="currentMenu"
      :is-edit="isEdit"
      :menu-tree="tableData"
      @success="fetchData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Sort } from '@element-plus/icons-vue'
import { getMenuTree, deleteMenu } from '@/api/system'
import { getIconComponent } from '@/utils/icon'
import type { MenuItem } from '@/types/api'
import MenuForm from './components/MenuForm.vue'

defineOptions({ name: 'MenuManagement' })

const { t } = useI18n()

/** 获取图标组件 */
function getIcon(iconName?: string) {
  return iconName ? getIconComponent(iconName) : undefined
}

// 表格数据
const loading = ref(false)
const tableData = ref<MenuItem[]>([])
const isExpand = ref(true)
const refreshTable = ref(true)

// 表单相关
const formVisible = ref(false)
const isEdit = ref(false)
const currentMenu = ref<MenuItem | null>(null)

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getMenuTree()
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
function handleAdd(row?: MenuItem) {
  currentMenu.value = row ? { ...row } as MenuItem : null
  isEdit.value = false
  formVisible.value = true
}

/** 编辑 */
function handleEdit(row: MenuItem) {
  currentMenu.value = { ...row }
  isEdit.value = true
  formVisible.value = true
}

/** 删除 */
async function handleDelete(row: MenuItem) {
  try {
    await ElMessageBox.confirm(t('common.confirmDelete'), t('common.tip'), {
      type: 'warning',
    })
    await deleteMenu(row.id)
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

<style scoped>
.menu-container {
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

.page-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 内容区域 */
.content-section {
  background: var(--el-bg-color);
  border-radius: 10px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
  padding: 20px;
}
</style>
