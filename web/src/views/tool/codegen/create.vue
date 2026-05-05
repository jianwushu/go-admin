<template>
  <div class="codegen-create p-4">
    <!-- 页面标题 -->
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('codegen.createConfig') }}
      </h1>
      <el-button @click="router.push('/tool/codegen')">
        <el-icon class="mr-1"><Back /></el-icon>
        {{ t('codegen.backToList') }}
      </el-button>
    </div>

    <!-- 步骤条 -->
    <el-card shadow="never" class="mb-4">
      <el-steps :active="activeStep" finish-status="success" align-center>
        <el-step :title="t('codegen.tableSelect')" :description="t('codegen.selectTableFirst')" />
        <el-step :title="t('codegen.basicConfig')" />
        <el-step :title="t('codegen.fieldConfig')" />
        <el-step :title="t('codegen.codePreview')" />
      </el-steps>
    </el-card>

    <!-- 步骤1：选择表 -->
    <el-card v-show="activeStep === 0" shadow="never" class="mb-4">
      <template #header>
        <span class="font-medium">{{ t('codegen.tableSelect') }}</span>
      </template>
      <TableSelect ref="tableSelectRef" @select="handleTableSelect" />
    </el-card>

    <!-- 步骤2：基础配置 -->
    <el-card v-show="activeStep === 1" shadow="never" class="mb-4">
      <template #header>
        <span class="font-medium">{{ t('codegen.basicConfig') }}</span>
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

    <!-- 步骤3：字段配置 -->
    <el-card v-show="activeStep === 2" shadow="never" class="mb-4">
      <template #header>
        <span class="font-medium">{{ t('codegen.fieldConfig') }}</span>
      </template>
      <FieldConfig :fields="configForm.fields" />
    </el-card>

    <!-- 步骤4：代码预览 -->
    <el-card v-show="activeStep === 3" shadow="never" class="mb-4 code-preview-card">
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
          <el-button v-if="activeStep < 3" type="primary" @click="handleNext">
            {{ t('common.confirm') }}
            <el-icon class="ml-1"><ArrowRight /></el-icon>
          </el-button>
          <el-button v-if="activeStep === 3" type="primary" :loading="generateLoading" @click="handleGenerate">
            <el-icon class="mr-1"><Download /></el-icon>
            {{ t('codegen.download') }}
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Back, View, ArrowLeft, ArrowRight, RefreshLeft, FolderChecked, Download } from '@element-plus/icons-vue'
import TableSelect from './components/TableSelect.vue'
import BasicConfigForm from './components/BasicConfigForm.vue'
import FieldConfig from './components/FieldConfig.vue'
import CodePreview from './components/CodePreview.vue'
import {
  getColumnList,
  previewCode,
  generateCode,
  saveConfig,
} from '@/api/codegen'
import type { TableInfo, ColumnConfig } from '@/api/codegen'

const { t } = useI18n()
const router = useRouter()

// ==================== 状态 ====================
const activeStep = ref(0)
const tableSelectRef = ref<InstanceType<typeof TableSelect>>()
const basicConfigRef = ref<InstanceType<typeof BasicConfigForm>>()
const previewLoading = ref(false)
const saveLoading = ref(false)
const generateLoading = ref(false)

const previewData = ref<Record<string, string>>({})

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

// ==================== 工具函数 ====================

/** 将下划线命名转为大驼峰 */
function toPascalCase(str: string): string {
  return str
    .split('_')
    .map(s => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase())
    .join('')
}

/** 将下划线命名转为小驼峰 */
function toCamelCase(str: string): string {
  const pascal = toPascalCase(str)
  return pascal.charAt(0).toLowerCase() + pascal.slice(1)
}

/** 根据数据库类型推断 Go 类型 */
function inferGoType(columnType: string, isPk: boolean): string {
  const type = columnType.toLowerCase()
  if (type.includes('bigint') || (isPk && type.includes('integer'))) return 'int64'
  if (type.includes('int') || type.includes('tinyint') || type.includes('smallint') || type.includes('mediumint')) return 'int'
  if (type.includes('float') || type.includes('double') || type.includes('decimal') || type.includes('numeric')) return 'float64'
  if (type.includes('bool')) return 'bool'
  if (type.includes('datetime') || type.includes('timestamp') || type.includes('date') || type.includes('time')) return 'time.Time'
  if (type.includes('blob') || type.includes('binary') || type.includes('varbinary')) return '[]byte'
  return 'string'
}

/** 根据数据库类型推断 TS 类型 */
function inferTsType(columnType: string): string {
  const type = columnType.toLowerCase()
  if (type.includes('int') || type.includes('float') || type.includes('double') || type.includes('decimal') || type.includes('numeric')) return 'number'
  if (type.includes('bool')) return 'boolean'
  if (type.includes('datetime') || type.includes('timestamp') || type.includes('date')) return 'string'
  return 'string'
}

/** 根据字段名推断表单类型 */
function inferHtmlType(columnName: string, columnType: string): string {
  const name = columnName.toLowerCase()
  const type = columnType.toLowerCase()
  if (name.includes('status') || name.includes('type') || name.includes('sex') || name.includes('gender')) return 'select'
  if (name.includes('content') || name.includes('description') || name.includes('remark') || name.includes('note') || name.includes('memo')) return 'textarea'
  if (name.includes('time') || name.includes('date') || name.includes('_at')) {
    if (type.includes('datetime') || type.includes('timestamp')) return 'datetime'
    return 'date'
  }
  if (name.includes('image') || name.includes('img') || name.includes('avatar') || name.includes('file') || name.includes('photo')) return 'upload'
  return 'input'
}

/** 根据字段名推断查询方式 */
function inferQueryType(columnName: string, columnType: string): string {
  const type = columnType.toLowerCase()
  if (type.includes('datetime') || type.includes('timestamp') || type.includes('date')) return 'between'
  if (type.includes('varchar') || type.includes('text') || type.includes('char')) return 'like'
  return 'eq'
}

/** 根据字段名推断中文标签 */
function inferLabel(columnName: string): string {
  const labelMap: Record<string, string> = {
    id: 'ID', name: '名称', title: '标题', status: '状态', type: '类型',
    sort: '排序', remark: '备注', description: '描述', content: '内容',
    created_at: '创建时间', updated_at: '更新时间', deleted_at: '删除时间',
    create_time: '创建时间', update_time: '更新时间',
    user_id: '用户ID', username: '用户名', password: '密码',
    email: '邮箱', phone: '手机号', avatar: '头像', nickname: '昵称',
    dept_id: '部门ID', role_id: '角色ID', menu_id: '菜单ID', parent_id: '父级ID',
  }
  if (labelMap[columnName]) return labelMap[columnName]
  return columnName.split('_').map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ')
}

// ==================== 事件处理 ====================

/** 选择表 */
async function handleTableSelect(table: TableInfo) {
  configForm.tableName = table.tableName
  configForm.tableComment = table.tableComment || ''
  configForm.className = toPascalCase(table.tableName)
  configForm.businessName = toCamelCase(table.tableName)
  configForm.functionName = table.tableComment || table.tableName

  try {
    const res = await getColumnList(table.tableName)
    const columns = res.data.data || []
    configForm.fields = columns.map((col, index) => ({
      columnName: col.columnName,
      columnType: col.columnType,
      goType: inferGoType(col.columnType, col.isPk),
      goField: toPascalCase(col.columnName),
      tsType: inferTsType(col.columnType),
      label: col.columnComment || inferLabel(col.columnName),
      htmlType: inferHtmlType(col.columnName, col.columnType),
      queryType: inferQueryType(col.columnName, col.columnType),
      isList: true,
      isQuery: !col.isPk && col.columnName !== 'created_at' && col.columnName !== 'updated_at' && col.columnName !== 'deleted_at' && col.columnName !== 'create_time' && col.columnName !== 'update_time',
      isRequired: !col.isNull && !col.isPk && col.columnDefault === '',
      isEdit: !col.isPk && col.columnName !== 'created_at' && col.columnName !== 'updated_at' && col.columnName !== 'create_time' && col.columnName !== 'update_time',
      dictType: '',
      comment: col.columnComment || '',
      isPk: col.isPk,
      isNull: col.isNull,
      maxLength: col.maxLength || 0,
      sort: index,
    }))
    activeStep.value = 1
  } catch (error) {
    console.error('获取列信息失败:', error)
  }
}

/** 下一步 */
async function handleNext() {
  if (activeStep.value === 1) {
    const valid = await basicConfigRef.value?.validate()
    if (!valid) return
    activeStep.value++
  } else if (activeStep.value === 0) {
    if (!configForm.tableName) {
      ElMessage.warning(t('codegen.tableNameRequired'))
      return
    }
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

/** 重置配置 */
function handleReset() {
  ElMessageBox.confirm(t('codegen.resetConfig') + '?', t('common.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(() => {
    configForm.tableName = ''
    configForm.tableComment = ''
    configForm.className = ''
    configForm.businessName = ''
    configForm.functionName = ''
    configForm.moduleName = ''
    configForm.packageName = 'go-admin'
    configForm.author = ''
    configForm.fields = []
    previewData.value = {}
    activeStep.value = 0
    tableSelectRef.value?.clearSelection()
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
