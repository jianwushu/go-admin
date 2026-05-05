<template>
  <div class="p-4">
    <el-row :gutter="20">
      <!-- 左侧：用户卡片 -->
      <el-col :span="8" :xs="24">
        <el-card shadow="never">
          <div class="profile-card">
            <div class="profile-avatar">
              <el-avatar :size="100" :src="profileData.avatar">
                {{ profileData.nickname?.charAt(0) || 'U' }}
              </el-avatar>
            </div>
            <h3 class="profile-name">{{ profileData.nickname || profileData.username }}</h3>
            <p class="profile-dept">{{ profileData.deptName || t('profile.noDept') }}</p>
            <el-divider />
            <div class="profile-info-list">
              <div class="profile-info-item">
                <el-icon><User /></el-icon>
                <span>{{ t('profile.username') }}：{{ profileData.username }}</span>
              </div>
              <div class="profile-info-item">
                <el-icon><Message /></el-icon>
                <span>{{ t('profile.email') }}：{{ profileData.email || '-' }}</span>
              </div>
              <div class="profile-info-item">
                <el-icon><Phone /></el-icon>
                <span>{{ t('profile.phone') }}：{{ profileData.phone || '-' }}</span>
              </div>
              <div class="profile-info-item">
                <el-icon><Timer /></el-icon>
                <span>{{ t('profile.createdAt') }}：{{ profileData.createdAt }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：编辑区域 -->
      <el-col :span="16" :xs="24">
        <el-card shadow="never">
          <el-tabs v-model="activeTab">
            <!-- 基本信息 -->
            <el-tab-pane :label="t('profile.basicInfo')" name="basic">
              <el-form
                ref="profileFormRef"
                :model="profileForm"
                :rules="profileRules"
                label-width="100px"
                style="max-width: 500px; margin-top: 20px"
              >
                <el-form-item :label="t('profile.nickname')" prop="nickname">
                  <el-input v-model="profileForm.nickname" :placeholder="t('profile.inputNickname')" />
                </el-form-item>
                <el-form-item :label="t('profile.email')" prop="email">
                  <el-input v-model="profileForm.email" :placeholder="t('profile.inputEmail')" />
                </el-form-item>
                <el-form-item :label="t('profile.phone')" prop="phone">
                  <el-input v-model="profileForm.phone" :placeholder="t('profile.inputPhone')" />
                </el-form-item>
                <el-form-item :label="t('profile.avatar')" prop="avatar">
                  <el-input v-model="profileForm.avatar" placeholder="请输入头像URL" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="profileSaving" @click="handleSaveProfile">
                    {{ t('common.save') }}
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <!-- 修改密码 -->
            <el-tab-pane :label="t('profile.changePassword')" name="password">
              <el-form
                ref="passwordFormRef"
                :model="passwordForm"
                :rules="passwordRules"
                label-width="120px"
                style="max-width: 500px; margin-top: 20px"
              >
                <el-form-item :label="t('profile.oldPassword')" prop="oldPassword">
                  <el-input
                    v-model="passwordForm.oldPassword"
                    type="password"
                    show-password
                    :placeholder="t('profile.inputOldPassword')"
                  />
                </el-form-item>
                <el-form-item :label="t('profile.newPassword')" prop="newPassword">
                  <el-input
                    v-model="passwordForm.newPassword"
                    type="password"
                    show-password
                    :placeholder="t('profile.inputNewPassword')"
                  />
                </el-form-item>
                <el-form-item :label="t('profile.confirmPassword')" prop="confirmPassword">
                  <el-input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    show-password
                    :placeholder="t('profile.inputConfirmPassword')"
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="passwordSaving" @click="handleChangePassword">
                    {{ t('common.save') }}
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Message, Phone, Timer } from '@element-plus/icons-vue'
import { getUserProfile, updateUserProfile, changePassword } from '@/api/auth'
import type { UserProfile } from '@/api/auth'
import { useUserStore } from '@/store/modules/user'

defineOptions({ name: 'Profile' })

const { t } = useI18n()
const userStore = useUserStore()

const activeTab = ref('basic')
const profileSaving = ref(false)
const passwordSaving = ref(false)

// 个人资料数据
const profileData = ref<UserProfile>({
  id: 0,
  username: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  deptId: 0,
  deptName: '',
  remark: '',
  createdAt: '',
})

// 个人资料表单
const profileFormRef = ref<FormInstance>()
const profileForm = reactive({
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
})

// 密码表单
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

// 个人资料校验规则
const profileRules: FormRules = {
  nickname: [
    { max: 64, message: '昵称长度不能超过64个字符', trigger: 'blur' },
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  phone: [
    { max: 20, message: '手机号长度不能超过20个字符', trigger: 'blur' },
  ],
}

// 密码校验规则
const passwordRules: FormRules = {
  oldPassword: [
    { required: true, message: () => t('profile.inputOldPassword'), trigger: 'blur' },
    { min: 6, max: 128, message: '密码长度在6到128个字符', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: () => t('profile.inputNewPassword'), trigger: 'blur' },
    { min: 6, max: 128, message: '密码长度在6到128个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: () => t('profile.inputConfirmPassword'), trigger: 'blur' },
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error(t('profile.passwordNotMatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

/** 加载个人资料 */
async function loadProfile() {
  try {
    const { data: res } = await getUserProfile()
    if (res.code === 0) {
      profileData.value = res.data
      // 同步到表单
      profileForm.nickname = res.data.nickname || ''
      profileForm.email = res.data.email || ''
      profileForm.phone = res.data.phone || ''
      profileForm.avatar = res.data.avatar || ''
    }
  } catch {
    // ignore
  }
}

/** 保存个人资料 */
async function handleSaveProfile() {
  const valid = await profileFormRef.value?.validate().catch(() => false)
  if (!valid) return

  profileSaving.value = true
  try {
    const { data: res } = await updateUserProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone,
      avatar: profileForm.avatar,
    })
    if (res.code === 0) {
      ElMessage.success(t('profile.updateSuccess'))
      // 刷新个人资料
      await loadProfile()
      // 刷新 store 中的用户信息
      await userStore.fetchUserInfo()
    }
  } catch {
    // ignore
  } finally {
    profileSaving.value = false
  }
}

/** 修改密码 */
async function handleChangePassword() {
  const valid = await passwordFormRef.value?.validate().catch(() => false)
  if (!valid) return

  passwordSaving.value = true
  try {
    const { data: res } = await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    if (res.code === 0) {
      ElMessage.success(t('profile.passwordChangeSuccess'))
      // 清空密码表单
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
      passwordFormRef.value?.resetFields()
    }
  } catch {
    // ignore
  } finally {
    passwordSaving.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-card {
  text-align: center;
  padding: 20px 0;
}

.profile-avatar {
  margin-bottom: 16px;
}

.profile-name {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 8px;
  color: var(--el-text-color-primary);
}

.profile-dept {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin: 0;
}

.profile-info-list {
  text-align: left;
  padding: 0 20px;
}

.profile-info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 0;
  font-size: 14px;
  color: var(--el-text-color-regular);
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.profile-info-item:last-child {
  border-bottom: none;
}

.profile-info-item .el-icon {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}
</style>
