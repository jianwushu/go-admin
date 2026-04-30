<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? t('common.edit') + t('dept.title') : t('common.add') + t('dept.title')"
    width="600px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item :label="t('dept.parentDept')" prop="parentId">
        <el-tree-select
          v-model="formData.parentId"
          :data="deptTreeOptions"
          :props="{ label: 'name', children: 'children' }"
          :placeholder="t('dept.parentDeptPlaceholder')"
          check-strictly
          clearable
          style="width: 100%"
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('dept.name')" prop="name">
            <el-input
              v-model="formData.name"
              :placeholder="t('dept.namePlaceholder')"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('dept.sort')" prop="sort">
            <el-input-number
              v-model="formData.sort"
              :min="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('dept.leader')" prop="leader">
            <el-input
              v-model="formData.leader"
              :placeholder="t('dept.leaderPlaceholder')"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('dept.phone')" prop="phone">
            <el-input
              v-model="formData.phone"
              :placeholder="t('dept.phonePlaceholder')"
              maxlength="20"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item :label="t('dept.email')" prop="email">
        <el-input
          v-model="formData.email"
          :placeholder="t('dept.emailPlaceholder')"
          maxlength="128"
        />
      </el-form-item>

      <el-form-item :label="t('dept.status')" prop="status">
        <el-radio-group v-model="formData.status">
          <el-radio :value="1">{{ t('dept.enabled') }}</el-radio>
          <el-radio :value="0">{{ t('dept.disabled') }}</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">
        {{ t('common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
        {{ t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createDept, updateDept } from '@/api/system'
import type { DeptInfo, DeptFormData } from '@/types/api'

const props = defineProps<{
  visible: boolean
  data: DeptInfo | null
  isEdit: boolean
  deptTree: DeptInfo[]
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const submitLoading = ref(false)

/** 部门树选项（用于上级部门选择，添加一个"顶级部门"选项） */
const deptTreeOptions = computed(() => {
  return [
    { id: 0, name: t('dept.topDept'), children: props.deptTree },
  ]
})

/** 表单数据 */
const formData = reactive<DeptFormData>({
  id: 0,
  parentId: 0,
  name: '',
  sort: 0,
  status: 1,
  leader: '',
  phone: '',
  email: '',
})

/** 表单校验规则 */
const formRules: FormRules = {
  name: [
    { required: true, message: () => t('dept.nameRequired'), trigger: 'blur' },
  ],
}

/** 监听弹窗打开，初始化表单 */
watch(() => props.visible, (val) => {
  if (val) {
    if (props.isEdit && props.data) {
      // 编辑模式：填充数据
      formData.id = props.data.id
      formData.parentId = props.data.parentId
      formData.name = props.data.name
      formData.sort = props.data.sort
      formData.status = props.data.status
      formData.leader = props.data.leader
      formData.phone = props.data.phone
      formData.email = props.data.email
    } else {
      // 新增模式：重置表单
      formData.id = 0
      formData.parentId = props.data?.id || 0
      formData.name = ''
      formData.sort = 0
      formData.status = 1
      formData.leader = ''
      formData.phone = ''
      formData.email = ''
    }
  }
})

/** 提交表单 */
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    if (props.isEdit) {
      await updateDept(formData)
    } else {
      await createDept(formData)
    }
    ElMessage.success(t('common.success'))
    emit('update:visible', false)
    emit('success')
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    submitLoading.value = false
  }
}
</script>
