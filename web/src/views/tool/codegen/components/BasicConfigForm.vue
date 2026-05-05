<template>
  <el-form
    ref="formRef"
    :model="formData"
    :rules="rules"
    label-width="120px"
    class="max-w-3xl"
  >
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="t('codegen.tableName')" prop="tableName">
          <el-input v-model="formData.tableName" :disabled="true" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="t('codegen.tableComment')" prop="tableComment">
          <el-input v-model="formData.tableComment" :placeholder="t('codegen.tableCommentPlaceholder')" />
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="t('codegen.className')" prop="className">
          <el-input v-model="formData.className" :placeholder="t('codegen.classNamePlaceholder')" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="t('codegen.businessName')" prop="businessName">
          <el-input v-model="formData.businessName" :placeholder="t('codegen.businessNamePlaceholder')" />
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="t('codegen.functionName')" prop="functionName">
          <el-input v-model="formData.functionName" :placeholder="t('codegen.functionNamePlaceholder')" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="t('codegen.moduleName')" prop="moduleName">
          <el-input v-model="formData.moduleName" :placeholder="t('codegen.moduleNamePlaceholder')" />
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="t('codegen.packageName')" prop="packageName">
          <el-input v-model="formData.packageName" :placeholder="t('codegen.packageNamePlaceholder')" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="t('codegen.author')" prop="author">
          <el-input v-model="formData.author" :placeholder="t('codegen.authorPlaceholder')" />
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { FormInstance, FormRules } from 'element-plus'
import type { ColumnConfig } from '@/api/codegen'

const { t } = useI18n()

const props = defineProps<{
  tableName: string
  tableComment: string
  className: string
  businessName: string
  functionName: string
  moduleName: string
  packageName: string
  author: string
  fields: ColumnConfig[]
}>()

const emit = defineEmits<{
  (e: 'update:tableComment', val: string): void
  (e: 'update:className', val: string): void
  (e: 'update:businessName', val: string): void
  (e: 'update:functionName', val: string): void
  (e: 'update:moduleName', val: string): void
  (e: 'update:packageName', val: string): void
  (e: 'update:author', val: string): void
}>()

const formRef = ref<FormInstance>()

const formData = reactive({
  tableName: props.tableName,
  tableComment: props.tableComment,
  className: props.className,
  businessName: props.businessName,
  functionName: props.functionName,
  moduleName: props.moduleName,
  packageName: props.packageName,
  author: props.author,
})

const rules: FormRules = {
  className: [{ required: true, message: t('codegen.classNameRequired'), trigger: 'blur' }],
  businessName: [{ required: true, message: t('codegen.businessNameRequired'), trigger: 'blur' }],
  moduleName: [{ required: true, message: t('codegen.moduleNameRequired'), trigger: 'blur' }],
  packageName: [{ required: true, message: t('codegen.packageNameRequired'), trigger: 'blur' }],
}

// 同步外部 props 变化到内部 formData
watch(() => props.tableName, (val) => { formData.tableName = val })
watch(() => props.tableComment, (val) => { formData.tableComment = val })
watch(() => props.className, (val) => { formData.className = val })
watch(() => props.businessName, (val) => { formData.businessName = val })
watch(() => props.functionName, (val) => { formData.functionName = val })
watch(() => props.moduleName, (val) => { formData.moduleName = val })
watch(() => props.packageName, (val) => { formData.packageName = val })
watch(() => props.author, (val) => { formData.author = val })

// 同步内部 formData 变化到外部
watch(() => formData.tableComment, (val) => { emit('update:tableComment', val) })
watch(() => formData.className, (val) => { emit('update:className', val) })
watch(() => formData.businessName, (val) => { emit('update:businessName', val) })
watch(() => formData.functionName, (val) => { emit('update:functionName', val) })
watch(() => formData.moduleName, (val) => { emit('update:moduleName', val) })
watch(() => formData.packageName, (val) => { emit('update:packageName', val) })
watch(() => formData.author, (val) => { emit('update:author', val) })

/** 验证表单 */
async function validate(): Promise<boolean> {
  try {
    await formRef.value?.validate()
    return true
  } catch {
    return false
  }
}

/** 重置表单 */
function resetFields() {
  formRef.value?.resetFields()
}

defineExpose({ validate, resetFields })
</script>
