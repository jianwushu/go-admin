<template>
  <div class="codegen p-4">
    <!-- 页面标题 -->
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('codegen.title') }}
      </h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon class="mr-1"><Plus /></el-icon>
        {{ t('codegen.createConfig') }}
      </el-button>
    </div>

    <!-- 搜索栏 -->
    <el-card shadow="never" class="mb-4">
      <el-input
        v-model="searchKeyword"
        :placeholder="t('codegen.searchPlaceholder')"
        clearable
        style="width: 320px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </el-card>

    <!-- 配置列表 -->
    <el-card shadow="never">
      <el-table
        :data="filteredConfigs"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="tableName" :label="t('codegen.tableName')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="tableComment" :label="t('codegen.tableComment')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="className" :label="t('codegen.className')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="businessName" :label="t('codegen.businessName')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="functionName" :label="t('codegen.functionName')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="moduleName" :label="t('codegen.moduleName')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="author" :label="t('codegen.author')" min-width="100" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="t('codegen.createdAt')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="t('common.operation')" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="info" link size="small" @click="handlePreview(row)">
              <el-icon class="mr-1"><View /></el-icon>
              {{ t('codegen.preview') }}
            </el-button>
            <el-button type="success" link size="small" @click="handleDownload(row)">
              <el-icon class="mr-1"><Download /></el-icon>
              {{ t('codegen.download') }}
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

    </el-card>

    <!-- 预览抽屉 -->
    <PreviewDrawer v-model="drawerVisible" :config="previewConfig" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Edit, View, Download, Delete, DocumentCopy } from '@element-plus/icons-vue'
import PreviewDrawer from './components/PreviewDrawer.vue'
import { getAllConfigs, deleteConfig, generateCode } from '@/api/codegen'
import type { CodegenConfig } from '@/api/codegen'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const configs = ref<CodegenConfig[]>([])
const searchKeyword = ref('')
const drawerVisible = ref(false)
const previewConfig = ref<CodegenConfig | null>(null)

/** 搜索过滤 */
const filteredConfigs = computed(() => {
  const keyword = searchKeyword.value.toLowerCase().trim()
  if (!keyword) return configs.value
  return configs.value.filter(
    (c) =>
      c.tableName.toLowerCase().includes(keyword) ||
      c.className.toLowerCase().includes(keyword) ||
      c.businessName.toLowerCase().includes(keyword) ||
      (c.functionName && c.functionName.toLowerCase().includes(keyword)),
  )
})

/** 加载配置列表 */
async function loadConfigs() {
  loading.value = true
  try {
    const res = await getAllConfigs()
    configs.value = res.data.data || []
  } catch (error) {
    console.error('获取配置列表失败:', error)
  } finally {
    loading.value = false
  }
}

/** 新建配置 */
function handleCreate() {
  router.push('/tool/codegen/create')
}

/** 编辑配置 */
function handleEdit(config: CodegenConfig) {
  router.push(`/tool/codegen/edit/${config.id}`)
}

/** 预览代码 */
function handlePreview(config: CodegenConfig) {
  previewConfig.value = config
  drawerVisible.value = true
}

/** 下载代码 */
async function handleDownload(config: CodegenConfig) {
  try {
    const fields = Array.isArray(config.fields) ? config.fields : []
    const response = await generateCode({
      tableName: config.tableName,
      tableComment: config.tableComment,
      className: config.className,
      businessName: config.businessName,
      functionName: config.functionName,
      moduleName: config.moduleName,
      packageName: config.packageName,
      author: config.author,
      fields,
    })

    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${config.businessName}_${new Date().toISOString().slice(0, 10)}.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success(t('codegen.generateSuccess'))
  } catch (error) {
    console.error('代码生成失败:', error)
  }
}

/** 删除配置 */
async function handleDelete(config: CodegenConfig) {
  try {
    await ElMessageBox.confirm(t('codegen.deleteConfirm'), t('common.tip'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
    })
    await deleteConfig(config.id!)
    ElMessage.success(t('codegen.deleteSuccess'))
    await loadConfigs()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  loadConfigs()
})
</script>
