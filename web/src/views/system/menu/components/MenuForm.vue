<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? t('common.edit') + t('menu.title') : t('common.add') + t('menu.title')"
    width="650px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="110px"
    >
      <el-form-item :label="t('menu.parentMenu')" prop="parentId">
        <el-tree-select
            v-model="formData.parentId"
            :data="menuTreeOptions"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            :placeholder="t('menu.parentMenuPlaceholder')"
            check-strictly
            clearable
            style="width: 100%"
          />
      </el-form-item>

      <el-form-item :label="t('menu.type')" prop="type">
        <el-radio-group v-model="formData.type">
          <el-radio :value="0">{{ t('menu.typeDir') }}</el-radio>
          <el-radio :value="1">{{ t('menu.typeMenu') }}</el-radio>
          <el-radio :value="2">{{ t('menu.typeButton') }}</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('menu.name')" prop="name">
            <el-input
              v-model="formData.name"
              :placeholder="t('menu.namePlaceholder')"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('menu.i18nKey')" prop="i18nKey">
            <el-input
              v-model="formData.i18nKey"
              :placeholder="t('menu.i18nKeyPlaceholder')"
              maxlength="128"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('menu.sort')" prop="sort">
            <el-input-number
              v-model="formData.sort"
              :min="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row v-if="formData.type !== 2" :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('menu.icon')" prop="icon">
            <el-input
              v-model="formData.icon"
              :placeholder="t('menu.iconPlaceholder')"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('menu.path')" prop="path">
            <el-input
              v-model="formData.path"
              :placeholder="t('menu.pathPlaceholder')"
              maxlength="128"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item v-if="formData.type === 1" :label="t('menu.component')" prop="component">
        <el-input
          v-model="formData.component"
          :placeholder="t('menu.componentPlaceholder')"
          maxlength="128"
        />
      </el-form-item>

      <el-form-item v-if="formData.type !== 0" :label="t('menu.perms')" prop="perms">
        <el-input
          v-model="formData.perms"
          :placeholder="t('menu.permsPlaceholder')"
          maxlength="128"
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('menu.visible')" prop="visible">
            <el-radio-group v-model="formData.visible">
              <el-radio :value="1">{{ t('menu.visibleShow') }}</el-radio>
              <el-radio :value="0">{{ t('menu.visibleHide') }}</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('menu.status')" prop="status">
            <el-radio-group v-model="formData.status">
              <el-radio :value="1">{{ t('menu.enabled') }}</el-radio>
              <el-radio :value="0">{{ t('menu.disabled') }}</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>
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
import { createMenu, updateMenu } from '@/api/system'
import type { MenuItem, MenuFormData } from '@/types/api'

const props = defineProps<{
  visible: boolean
  data: MenuItem | null
  isEdit: boolean
  menuTree: MenuItem[]
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const submitLoading = ref(false)

/** 菜单树选项（用于上级菜单选择，添加一个"顶级菜单"选项） */
const menuTreeOptions = computed(() => {
  return [
    { id: 0, name: t('menu.topMenu'), children: props.menuTree },
  ]
})

/** 表单数据 */
const formData = reactive<MenuFormData>({
  id: 0,
  parentId: 0,
  name: '',
  i18nKey: '',
  path: '',
  component: '',
  icon: '',
  type: 0,
  sort: 0,
  visible: 1,
  status: 1,
  perms: '',
})

/** 表单校验规则 */
const formRules: FormRules = {
  name: [
    { required: true, message: () => t('menu.nameRequired'), trigger: 'blur' },
  ],
  type: [
    { required: true, message: () => t('menu.typeRequired'), trigger: 'change' },
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
      formData.i18nKey = props.data.i18nKey || ''
      formData.path = props.data.path
      formData.component = props.data.component
      formData.icon = props.data.icon
      formData.type = props.data.type
      formData.sort = props.data.sort
      formData.visible = props.data.visible
      formData.status = props.data.status
      formData.perms = props.data.perms
    } else {
      // 新增模式：重置表单
      formData.id = 0
      formData.parentId = props.data?.id || 0
      formData.name = ''
      formData.i18nKey = ''
      formData.path = ''
      formData.component = ''
      formData.icon = ''
      formData.type = 0
      formData.sort = 0
      formData.visible = 1
      formData.status = 1
      formData.perms = ''
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
      await updateMenu(formData)
    } else {
      await createMenu(formData)
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
