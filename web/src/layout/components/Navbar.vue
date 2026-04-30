<template>
  <div class="navbar">
    <!-- 左侧：折叠按钮 + 面包屑 -->
    <div class="navbar-left">
      <div class="hamburger" @click="toggleSidebar">
        <el-icon :size="20">
          <Fold v-if="appStore.sidebar.opened" />
          <Expand v-else />
        </el-icon>
      </div>
      <Breadcrumb />
    </div>

    <!-- 右侧：功能区 -->
    <div class="navbar-right">
      <!-- 全屏按钮 -->
      <div class="navbar-action" :title="t('navbar.fullscreen')" @click="toggleFullscreen">
        <el-icon :size="18">
          <FullScreen />
        </el-icon>
      </div>

      <!-- 主题切换 -->
      <el-tooltip :content="appStore.theme === 'light' ? t('navbar.darkMode') : t('navbar.lightMode')" placement="bottom">
        <div class="navbar-action" @click="toggleTheme">
          <el-icon :size="18">
            <Moon v-if="appStore.theme === 'light'" />
            <Sunny v-else />
          </el-icon>
        </div>
      </el-tooltip>

      <!-- 语言切换 -->
      <el-dropdown trigger="click" @command="handleSetLanguage">
        <div class="navbar-action" :title="t('navbar.language')">
          <el-icon :size="18">
            <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="18" height="18">
              <path
                fill="currentColor"
                d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm0 820c-205.4 0-372-166.6-372-372s166.6-372 372-372 372 166.6 372 372-166.6 372-372 372z"
              />
              <path
                fill="currentColor"
                d="M686.8 448.7H580.3c-11.6 0-21.4 7.3-24.8 17.8l-30.2 93.4h-79.4l-30.2-93.4c-3.4-10.5-13.2-17.8-24.8-17.8H335.2c-15.3 0-25.7 15.8-19.7 29.6l52.3 120.2c3.4 7.8 3.4 16.7 0 24.5l-52.3 120.2c-6 13.8 4.4 29.6 19.7 29.6h56.3c11.6 0 21.4-7.3 24.8-17.8l26.6-82.2h90.2l26.6 82.2c3.4 10.5 13.2 17.8 24.8 17.8h56.3c15.3 0 25.7-15.8 19.7-29.6l-52.3-120.2c-3.4-7.8-3.4-16.7 0-24.5l52.3-120.2c6-13.8-4.4-29.6-19.7-29.6zM544 660l-18.4-56.8c-2.1-6.5-2.1-13.5 0-20l18.4-56.8h47.2l18.4 56.8c2.1 6.5 2.1 13.5 0 20L591.2 660H544zM768 352H640c-17.7 0-32 14.3-32 32s14.3 32 32 32h44.4L614.6 544h-57.2l-45.8-128h-59.2l-45.8 128h-57.2L339.6 416H384c17.7 0 32-14.3 32-32s-14.3-32-32-32H256c-17.7 0-32 14.3-32 32s14.3 32 32 32h22.8l-59.8 168c-1.2 3.4-1.9 7-1.9 10.7 0 17.7 14.3 32 32 32h70.4c17.7 0 32-14.3 32-32 0-3.7-.7-7.3-1.9-10.7L396 416h38.4l43.6 128h57.2l43.6-128H628l-59.8 168c-1.2 3.4-1.9 7-1.9 10.7 0 17.7 14.3 32 32 32h70.4c17.7 0 32-14.3 32-32 0-3.7-.7-7.3-1.9-10.7L647.6 416H768c17.7 0 32-14.3 32-32s-14.3-32-32-32z"
              />
            </svg>
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item :disabled="appStore.language === 'zh-CN'" command="zh-CN">
              中文
            </el-dropdown-item>
            <el-dropdown-item :disabled="appStore.language === 'en'" command="en">
              English
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- 用户信息 -->
      <el-dropdown trigger="click" @command="handleUserCommand">
        <div class="navbar-user">
          <el-avatar :size="28" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="navbar-username">{{ userStore.userInfo?.nickname || 'Admin' }}</span>
          <el-icon :size="14"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">{{ t('navbar.profile') }}</el-dropdown-item>
            <el-dropdown-item divided command="logout">{{ t('navbar.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import { Fold, Expand, FullScreen, Moon, Sunny, ArrowDown } from '@element-plus/icons-vue'
import { useAppStore } from '@/store/modules/app'
import { useUserStore } from '@/store/modules/user'
import { usePermissionStore } from '@/store/modules/permission'
import Breadcrumb from '@/components/Breadcrumb.vue'

defineOptions({ name: 'Navbar' })

const router = useRouter()
const { t, locale } = useI18n()
const appStore = useAppStore()
const userStore = useUserStore()
const permissionStore = usePermissionStore()

/** 切换侧边栏 */
function toggleSidebar() {
  appStore.toggleSidebar()
}

/** 切换全屏 */
function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

/** 切换主题 */
function toggleTheme() {
  appStore.toggleTheme()
}

/** 设置语言 */
function handleSetLanguage(lang: string) {
  locale.value = lang
  appStore.setLanguage(lang)
}

/** 用户菜单命令 */
function handleUserCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      handleLogout()
      break
  }
}

/** 退出登录 */
async function handleLogout() {
  try {
    await ElMessageBox.confirm(t('navbar.logoutConfirm'), t('common.tip'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
    })
    // 调用后端登出接口
    const { logout } = await import('@/api/auth')
    await logout().catch(() => {})
    // 清除本地状态
    userStore.logout()
    permissionStore.resetPermission()
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--navbar-height);
  padding: 0 16px;
  background: inherit;
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hamburger {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s;
}

.hamburger:hover {
  background: rgba(0, 0, 0, 0.06);
}

html.dark .hamburger:hover {
  background: rgba(255, 255, 255, 0.08);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.navbar-action {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
  color: var(--el-text-color-regular);
}

.navbar-action:hover {
  background: rgba(0, 0, 0, 0.06);
  color: var(--el-color-primary);
}

html.dark .navbar-action:hover {
  background: rgba(255, 255, 255, 0.08);
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.navbar-user:hover {
  background: rgba(0, 0, 0, 0.06);
}

html.dark .navbar-user:hover {
  background: rgba(255, 255, 255, 0.08);
}

.navbar-username {
  font-size: 14px;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--el-text-color-regular);
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .navbar {
    padding: 0 12px;
  }
  
  .navbar-username {
    display: none;
  }
  
  .navbar-action {
    width: 32px;
    height: 32px;
  }
}

/* 平板端适配 */
@media screen and (min-width: 769px) and (max-width: 992px) {
  .navbar-username {
    max-width: 80px;
  }
}
</style>
