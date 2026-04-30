<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? t('common.edit') + t('user.title') : t('common.add') + t('user.title')"
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
      <el-form-item :label="t('user.username')" prop="username">
        <el-input
          v-model="formData.username"
          :placeholder="t('user.usernamePlaceholder')"
          :disabled="isEdit"
          maxlength="64"
        />
      </el-form-item>

      <el-form-item v-if="!isEdit" :label="t('user.password')" prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          :placeholder="t('user.passwordPlaceholder')"
          show-password
          maxlength="128"
        />
      </el-form-item>

      <el-form-item :label="t('user.nickname')" prop="nickname">
        <el-input
          v-model="formData.nickname"
          :placeholder="t('user.nicknamePlaceholder')"
          maxlength="64"
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('user.email')" prop="email">
            <el-input
              v-model="formData.email"
              :placeholder="t('user.emailPlaceholder')"
              maxlength="128"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('user.phone')" prop="phone">
            <el-input
              v-model="formData.phone"
              :placeholder="t('user.phonePlaceholder')"
              maxlength="20"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('user.status')" prop="status">
            <el-radio-group v-model="formData.status">
              <el-radio :value="1">{{ t('user.enabled') }}</el-radio>
              <el-radio :value="0">{{ t('user.disabled') }}</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('user.dept')" prop="deptId">
            <el-input-number
              v-model="formData.deptId"
              :min="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item :label="t('user.roles')" prop="roleIds">
        <el-select
          v-model="formData.roleIds"
          multiple
          :placeholder="t('user.selectRoles')"
          style="width: 100%"
        >
          <el-option
            v-for="role in roleOptions"
            :key="role.id"
            :label="role.name"
            :value="role.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('user.remark')" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :placeholder="t('user.remarkPlaceholder')"
          :rows="3"
          maxlength="512"
        />
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
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createUser, updateUser, getAllRoles } from '@/api/system'
import type { UserItem, RoleInfo } from '@/types/api'

const props = defineProps<{
  visible: boolean
  data: UserItem | null
  isEdit: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const roleOptions = ref<RoleInfo[]>([])

/** 表单数据 */
const formData = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  status: 1,
  deptId: 0,
  avatar: '',
  remark: '',
  roleIds: [] as number[],
})

/** 表单校验规则 */
const formRules: FormRules = {
  username: [
    { required: true, message: () => t('user.usernameRequired'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: () => t('user.passwordRequired'), trigger: 'blur' },
    { min: 6, message: () => t('user.passwordMin'), trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: () => t('user.nicknameRequired'), trigger: 'blur' },
  ],
  email: [
    { type: 'email', message: () => t('user.emailInvalid'), trigger: 'blur' },
  ],
}

/** 获取角色选项 */
async function fetchRoleOptions() {
  try {
    const { data: res } = await getAllRoles()
    roleOptions.value = res.data || []
  } catch {
    // 错误已在 request 拦截器中处理
  }
}

/** 监听弹窗打开，初始化表单 */
watch(() => props.visible, (val) => {
  if (val) {
    fetchRoleOptions()
    if (props.isEdit && props.data) {
      // 编辑模式：填充数据
      formData.id = props.data.id
      formData.username = props.data.username
      formData.password = ''
      formData.nickname = props.data.nickname
      formData.email = props.data.email
      formData.phone = props.data.phone
      formData.status = props.data.status
      formData.deptId = props.data.deptId
      formData.avatar = props.data.avatar
      formData.remark = props.data.remark
      formData.roleIds = props.data.roles?.map(r => r.id) || []
    } else {
      // 新增模式：重置表单
      formData.id = 0
      formData.username = ''
      formData.password = ''
      formData.nickname = ''
      formData.email = ''
      formData.phone = ''
      formData.status = 1
      formData.deptId = 0
      formData.avatar = ''
      formData.remark = ''
      formData.roleIds = []
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
      await updateUser({
        id: formData.id,
        nickname: formData.nickname,
        email: formData.email,
        phone: formData.phone,
        status: formData.status,
        deptId: formData.deptId,
        avatar: formData.avatar,
        remark: formData.remark,
        roleIds: formData.roleIds,
        username: formData.username,
      })
    } else {
      await createUser({
        username: formData.username,
        password: formData.password,
        nickname: formData.nickname,
        email: formData.email,
        phone: formData.phone,
        status: formData.status,
        deptId: formData.deptId,
        avatar: formData.avatar,
        remark: formData.remark,
        roleIds: formData.roleIds,
      })
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
