<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? t('common.edit') : t('common.add')"
    width="600px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item :label="t('job.name')" prop="name">
        <el-input v-model="formData.name" :placeholder="t('job.namePlaceholder')" maxlength="128" />
      </el-form-item>

      <el-form-item :label="t('job.jobType')" prop="jobType">
        <el-radio-group v-model="formData.jobType">
          <el-radio :value="1">{{ t('job.jobTypeFunc') }}</el-radio>
          <el-radio :value="2">{{ t('job.jobTypeHTTP') }}</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item :label="t('job.cronExpr')" prop="cronExpr">
        <el-input v-model="formData.cronExpr" :placeholder="t('job.cronExprPlaceholder')" maxlength="64" />
      </el-form-item>

      <!-- 内置函数选择 -->
      <el-form-item v-if="formData.jobType === 1" :label="t('job.funcName')" prop="funcName">
        <el-select v-model="formData.funcName" :placeholder="t('job.funcNameRequired')" style="width: 100%">
          <el-option
            v-for="(label, key) in builtinFuncs"
            :key="key"
            :label="label"
            :value="key"
          />
        </el-select>
      </el-form-item>

      <!-- HTTP请求配置 -->
      <template v-if="formData.jobType === 2">
        <el-form-item :label="t('job.httpMethod')" prop="httpMethod">
          <el-select v-model="formData.httpMethod" style="width: 120px">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('job.httpUrl')" prop="httpUrl">
          <el-input v-model="formData.httpUrl" :placeholder="t('job.httpUrlPlaceholder')" maxlength="512" />
        </el-form-item>

        <el-form-item :label="t('job.httpBody')">
          <el-input
            v-model="formData.httpBody"
            type="textarea"
            :placeholder="t('job.httpBodyPlaceholder')"
            :rows="4"
          />
        </el-form-item>
      </template>

      <el-form-item :label="t('job.status')">
        <el-radio-group v-model="formData.status">
          <el-radio :value="1">{{ t('job.enabled') }}</el-radio>
          <el-radio :value="2">{{ t('job.disabled') }}</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item :label="t('job.remark')">
        <el-input v-model="formData.remark" :placeholder="t('job.remarkPlaceholder')" maxlength="512" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ t('common.confirm') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { createJob, updateJob } from '@/api/job'
import type { JobItem } from '@/types/api'

const props = defineProps<{
  visible: boolean
  data: JobItem | null
}>()

const emit = defineEmits<{
  'update:visible': [val: boolean]
  success: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const isEdit = computed(() => !!props.data?.id)

// 内置函数列表
const builtinFuncs = computed(() => ({
  clean_operation_log: t('job.builtinFuncs.clean_operation_log'),
  clean_login_log: t('job.builtinFuncs.clean_login_log'),
}))

// 表单数据
const formData = reactive({
  id: 0,
  name: '',
  jobType: 1,
  cronExpr: '',
  funcName: '',
  httpUrl: '',
  httpMethod: 'GET',
  httpBody: '',
  status: 1,
  remark: '',
})

// 表单验证规则
const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: t('job.nameRequired'), trigger: 'blur' }],
  jobType: [{ required: true, message: t('job.jobTypeRequired'), trigger: 'change' }],
  cronExpr: [{ required: true, message: t('job.cronExprRequired'), trigger: 'blur' }],
  funcName: formData.jobType === 1
    ? [{ required: true, message: t('job.funcNameRequired'), trigger: 'change' }]
    : [],
  httpUrl: formData.jobType === 2
    ? [{ required: true, message: t('job.httpUrlRequired'), trigger: 'blur' }]
    : [],
}))

// 监听数据变化，填充表单
watch(() => props.visible, (val) => {
  if (val) {
    if (props.data) {
      Object.assign(formData, {
        id: props.data.id,
        name: props.data.name,
        jobType: props.data.jobType,
        cronExpr: props.data.cronExpr,
        funcName: props.data.funcName,
        httpUrl: props.data.httpUrl,
        httpMethod: props.data.httpMethod || 'GET',
        httpBody: props.data.httpBody,
        status: props.data.status,
        remark: props.data.remark,
      })
    } else {
      Object.assign(formData, {
        id: 0,
        name: '',
        jobType: 1,
        cronExpr: '',
        funcName: '',
        httpUrl: '',
        httpMethod: 'GET',
        httpBody: '',
        status: 1,
        remark: '',
      })
    }
  }
})

/** 提交表单 */
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateJob(formData as any)
      ElMessage.success(t('common.editSuccess'))
    } else {
      await createJob(formData as any)
      ElMessage.success(t('common.addSuccess'))
    }
    emit('update:visible', false)
    emit('success')
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    submitting.value = false
  }
}
</script>
