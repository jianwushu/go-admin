<template>
  <div class="codegen-edit p-4">
    <!-- 页面标题 -->
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('codegen.editConfig') }}
      </h1>
      <el-button @click="router.push('/tool/codegen')">
        <el-icon class="mr-1"><Back /></el-icon>
        {{ t('codegen.backToList') }}
      </el-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="pageLoading" class="flex h-64 items-center justify-center">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <template v-else>
      <!-- 步骤条 -->
      <el-card shadow="never" class="mb-4">
        <el-steps :active="activeStep" finish-status="success" align-center>
          <el-step :title="t('codegen.basicConfig')" />
          <el-step :title="t('codegen.fieldConfig')" />
          <el-step :title="t('codegen.codePreview')" />
        </el-steps>
      </el-card>

      <!-- 步骤1：基础配置 -->
      <el-card v-show="activeStep === 0" shadow="never" class="mb-4">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-medium">{{ t('codegen.basicConfig') }}</span>
            <el-tag type="info" size="small">
              {{ t('codegen.tableName') }}: {{ configForm.tableName }}
            </el-tag>
          </div>
        </template>
        <BasicConfigForm
          ref="basicConfigRef"
          :table-name="configForm.tableName"
          :table-comment="configForm.tableComment"
          :class-name="configForm.className"
          :business-name="configForm.businessName"
          :function-name="configForm.functionName"
          :module-name="configForm.moduleName"
          :package-name="configForm.packageName"
          :author="configForm.author"
          :fields="configForm.fields"
          @update:table-comment="configForm.tableComment = $event"
          @update:class-name="configForm.className = $event"
          @update:business-name="configForm.businessName = $event"
          @update:function-name="configForm.functionName = $event"
          @update:module-name="configForm.moduleName = $event"
          @update:package-name="configForm.packageName = $event"
          @update:author="configForm.author = $event"
        />
      </el-card>

      <!-- 步骤2：字段配置 -->
      <el-card v-show="activeStep === 1" shadow="never" class="mb-4">
        <template #header>
          <span class="font-medium">{{ t('codegen.fieldConfig') }}</span>
        </template>
        <FieldConfig :fields="configForm.fields" />
      </el-card>

      <!-- 步骤3：代码预览 -->
      <el-card v-show="activeStep === 2" shadow="never" class="mb-4 code-preview-card">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-medium">{{ t('codegen.codePreview') }}</span>
            <el-button type="primary" :loading="previewLoading" @click="handlePreview">
              <el-icon class="mr-1"><View /></el-icon>
              {{ t('codegen.preview') }}
            </el-button>
          </div>
        </template>
        <div class="code-preview-wrapper">
          <CodePreview :preview-data="previewData" />
        </div>
      </el-card>

      <!-- 底部操作栏 -->
      <el-card shadow="never">
        <div class="flex items-center justify-between">
          <div>
            <el-button v-if="activeStep > 0" @click="activeStep--">
              <el-icon class="mr-1"><ArrowLeft /></el-icon>
              {{ t('common.back') }}
            </el-button>
          </div>
          <div class="flex items-center gap-3">
            <el-button @click="handleReset">
              <el-icon class="mr-1"><RefreshLeft /></el-icon>
              {{ t('codegen.resetConfig') }}
            </el-button>
            <el-button type="success" @click="handleSaveConfig" :loading="saveLoading">
              <el-icon class="mr-1"><FolderChecked /></el-icon>
              {{ t('codegen.saveConfig') }}
            </el-button>
            <el-button v-if="activeStep < 2" type="primary" @click="handleNext">
              {{ t('common.confirm') }}
              <el-icon class="ml-1"><ArrowRight /></el-icon>
            </el-button>
            <el-button v-if="activeStep === 2" type="primary" :loading="generateLoading" @click="handleGenerate">
              <el-icon class="mr-1"><Download /></el-icon>
              {{ t('codegen.download') }}
            </el-button>
          </div>
        </div>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Back, View, ArrowLeft, ArrowRight, RefreshLeft, FolderChecked, Download, Loading } from '@element-plus/icons-vue'
import BasicConfigForm from './components/BasicConfigForm.vue'
import FieldConfig from './components/FieldConfig.vue'
import CodePreview from './components/CodePreview.vue'
import {
  previewCode,
  generateCode,
  saveConfig,
  getConfigById,
} from '@/api/codegen'
import type { ColumnConfig } from '@/api/codegen'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// ==================== 状态 ====================
const activeStep = ref(0)
const pageLoading = ref(true)
const basicConfigRef = ref<InstanceType<typeof BasicConfigForm>>()
const previewLoading = ref(false)
const saveLoading = ref(false)
const generateLoading = ref(false)

const previewData = ref<Record<string, string>>({})
const originalTableName = ref('')

/** 配置表单 */
const configForm = reactive({
  tableName: '',
  tableComment: '',
  className: '',
  businessName: '',
  functionName: '',
  moduleName: '',
  packageName: 'go-admin',
  author: '',
  fields: [] as ColumnConfig[],
})

// ==================== 初始化 ====================

onMounted(async () => {
  const id = Number(route.params.id)
  if (!id) {
    ElMessage.error('无效的配置ID')
    router.push('/tool/codegen')
    return
  }

  try {
    const res = await getConfigById(id)
    const config = res.data.data
    if (config) {
      configForm.tableName = config.tableName
      configForm.tableComment = config.tableComment || ''
      configForm.className = config.className
      configForm.businessName = config.businessName
      configForm.functionName = config.functionName || ''
      configForm.moduleName = config.moduleName
      configForm.packageName = config.packageName
      configForm.author = config.author || ''
      originalTableName.value = config.tableName

      if (config.fields && typeof config.fields === 'string') {
        configForm.fields = JSON.parse(config.fields)
      } else if (Array.isArray(config.fields)) {
        configForm.fields = config.fields
      }
    } else {
      ElMessage.error('配置不存在')
      router.push('/tool/codegen')
    }
  } catch (error) {
    console.error('加载配置失败:', error)
    ElMessage.error('加载配置失败')
    router.push('/tool/codegen')
  } finally {
    pageLoading.value = false
  }
})

// ==================== 事件处理 ====================

/** 下一步 */
async function handleNext() {
  if (activeStep.value === 0) {
    const valid = await basicConfigRef.value?.validate()
    if (!valid) return
    activeStep.value++
  } else {
    activeStep.value++
  }
}

/** 预览代码 */
async function handlePreview() {
  previewLoading.value = true
  try {
    const res = await previewCode({
      tableName: configForm.tableName,
      tableComment: configForm.tableComment,
      className: configForm.className,
      businessName: configForm.businessName,
      functionName: configForm.functionName,
      moduleName: configForm.moduleName,
      packageName: configForm.packageName,
      author: configForm.author,
      fields: configForm.fields,
    })
    const files = res.data.data?.files || []
    const previewMap: Record<string, string> = {}
    for (const file of files) {
      previewMap[file.filePath] = file.content
    }
    previewData.value = previewMap
  } catch (error) {
    console.error('代码预览失败:', error)
  } finally {
    previewLoading.value = false
  }
}

/** 生成并下载代码 */
async function handleGenerate() {
  if (Object.keys(previewData.value).length === 0) {
    ElMessage.warning(t('codegen.noPreview'))
    return
  }
  generateLoading.value = true
  try {
    const response = await generateCode({
      tableName: configForm.tableName,
      tableComment: configForm.tableComment,
      className: configForm.className,
      businessName: configForm.businessName,
      functionName: configForm.functionName,
      moduleName: configForm.moduleName,
      packageName: configForm.packageName,
      author: configForm.author,
      fields: configForm.fields,
    })

    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${configForm.businessName}_${new Date().toISOString().slice(0, 10)}.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success(t('codegen.generateSuccess'))
  } catch (error) {
    console.error('代码生成失败:', error)
  } finally {
    generateLoading.value = false
  }
}

/** 保存配置 */
async function handleSaveConfig() {
  if (!configForm.tableName) {
    ElMessage.warning(t('codegen.tableNameRequired'))
    return
  }
  saveLoading.value = true
  try {
    await saveConfig({
      tableName: configForm.tableName,
      tableComment: configForm.tableComment,
      className: configForm.className,
      businessName: configForm.businessName,
      functionName: configForm.functionName,
      moduleName: configForm.moduleName,
      packageName: configForm.packageName,
      author: configForm.author,
      fields: configForm.fields,
    })
    ElMessage.success(t('codegen.saveSuccess'))
  } catch (error) {
    console.error('保存配置失败:', error)
  } finally {
    saveLoading.value = false
  }
}

/** 重置配置（恢复到初始加载的状态） */
function handleReset() {
  ElMessageBox.confirm(t('codegen.resetConfig') + '?', t('common.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(async () => {
    const id = Number(route.params.id)
    try {
      const res = await getConfigById(id)
      const config = res.data.data
      if (config) {
        configForm.tableComment = config.tableComment || ''
        configForm.className = config.className
        configForm.businessName = config.businessName
        configForm.functionName = config.functionName || ''
        configForm.moduleName = config.moduleName
        configForm.packageName = config.packageName
        configForm.author = config.author || ''
        if (config.fields && typeof config.fields === 'string') {
          configForm.fields = JSON.parse(config.fields)
        } else if (Array.isArray(config.fields)) {
          configForm.fields = config.fields
        }
      }
      previewData.value = {}
      activeStep.value = 0
    } catch (error) {
      console.error('重置配置失败:', error)
    }
  }).catch(() => {})
}
</script>

<style scoped>
.code-preview-card :deep(.el-card__body) {
  padding: 0;
}

.code-preview-wrapper {
  height: calc(100vh - 320px);
  min-height: 500px;
}
</style>
