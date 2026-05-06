<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="login-bg">
      <div class="login-bg-circle login-bg-circle-1"></div>
      <div class="login-bg-circle login-bg-circle-2"></div>
      <div class="login-bg-circle login-bg-circle-3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <h1 class="login-title">{{ configStore.systemTitle }}</h1>
        <p class="login-subtitle">{{ t('login.title') }}</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        size="large"
        @keyup.enter="handleLogin"
      >
        <!-- 用户名 -->
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            :placeholder="t('login.usernamePlaceholder')"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <!-- 密码 -->
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            :placeholder="t('login.passwordPlaceholder')"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <!-- 登录按钮 -->
        <el-form-item>
          <el-button
            type="primary"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? t('login.logging') : t('login.login') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 主题切换按钮 -->
    <div class="login-theme-toggle" @click="toggleTheme">
      <el-icon :size="20">
        <Moon v-if="appStore.theme === 'light'" />
        <Sunny v-else />
      </el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock, Moon, Sunny } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { useAppStore } from '@/store/modules/app'
import { useConfigStore } from '@/store/modules/config'
import { login } from '@/api/auth'

defineOptions({ name: 'Login' })

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const userStore = useUserStore()
const appStore = useAppStore()
const configStore = useConfigStore()

/** 页面加载时获取系统配置 */
onMounted(() => {
  configStore.loadConfig()
})

const formRef = ref<FormInstance>()
const loading = ref(false)

/** 登录表单 */
const loginForm = reactive({
  username: 'admin',
  password: 'admin123',
})

/** 表单校验规则 */
const loginRules: FormRules = {
  username: [
    { required: true, message: () => t('login.usernameRequired'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: () => t('login.passwordRequired'), trigger: 'blur' },
  ],
}

/** 切换主题 */
function toggleTheme() {
  appStore.toggleTheme()
}

/** 登录 */
async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const { data: res } = await login({
      username: loginForm.username,
      password: loginForm.password,
    })

    // 存储 Token
    userStore.setToken(res.data.accessToken, res.data.refreshToken)

    ElMessage.success(t('login.loginSuccess'))

    // 跳转到目标页面或首页
    const redirect = (route.query.redirect as string) || '/'

    console.log(redirect)
    router.push(redirect)
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100vh;
  background: linear-gradient(135deg, #e0f2fe 0%, #dbeafe 50%, #e0e7ff 100%);
  position: relative;
  overflow: hidden;
  transition: background 0.3s ease;
}

html.dark .login-container {
  background: linear-gradient(135deg, #0c1929 0%, #0f172a 50%, #1e1b4b 100%);
}

/* 背景装饰圆 */
.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.login-bg-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

html.dark .login-bg-circle {
  background: rgba(255, 255, 255, 0.03);
}

.login-bg-circle-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  right: -100px;
}

.login-bg-circle-2 {
  width: 300px;
  height: 300px;
  bottom: -80px;
  left: -80px;
}

.login-bg-circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 10%;
}

/* 登录卡片 */
.login-card {
  width: 420px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(10px);
  z-index: 1;
  transition: background 0.3s ease, box-shadow 0.3s ease;
}

html.dark .login-card {
  background: rgba(31, 31, 31, 0.95);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

/* 头部 */
.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: #1f1f1f;
  margin-bottom: 8px;
  letter-spacing: 2px;
  transition: color 0.3s ease;
}

html.dark .login-title {
  color: #ffffffd9;
}

.login-subtitle {
  font-size: 14px;
  color: #8c8c8c;
  transition: color 0.3s ease;
}

html.dark .login-subtitle {
  color: #ffffffa6;
}

/* 表单 */
.login-form {
  width: 100%;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-color-primary) inset;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 44px;
  border-radius: 8px;
  font-size: 16px;
  letter-spacing: 4px;
}

/* 主题切换按钮 */
.login-theme-toggle {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(4px);
}

.login-theme-toggle:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(30deg);
}

html.dark .login-theme-toggle {
  background: rgba(255, 255, 255, 0.1);
}

html.dark .login-theme-toggle:hover {
  background: rgba(255, 255, 255, 0.15);
}
</style>
