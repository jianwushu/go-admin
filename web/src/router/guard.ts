import router from '@/router'
import { useUserStore } from '@/store/modules/user'
import { usePermissionStore } from '@/store/modules/permission'
import { getToken } from '@/utils/auth'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

/** 白名单路由（不需要登录） */
const whiteList = ['/login', '/404', '/403']

/** 是否已添加 404 兜底路由 */
let hasAdded404Route = false

router.beforeEach(async (to, _from, next) => {
  NProgress.start()

  const token = getToken()

  if (token) {
    if (to.path === '/login') {
      // 已登录，跳转首页
      next({ path: '/' })
      NProgress.done()
    } else {
      const userStore = useUserStore()
      const permissionStore = usePermissionStore()

      if (userStore.userInfo) {
        next()
      } else {
        try {
          // 获取用户信息
          await userStore.fetchUserInfo()
          // 生成动态路由
          const accessRoutes = await permissionStore.generateRoutes()
          // 添加动态路由
          accessRoutes.forEach((route) => {
            router.addRoute(route)
          })
          // 添加 404 兜底路由（必须在动态路由之后，且只添加一次）
          if (!hasAdded404Route) {
            router.addRoute({
              path: '/:pathMatch(.*)*',
              redirect: '/404',
            })
            hasAdded404Route = true
          }
          // 使用 replace 确保导航正确
          next({ ...to, replace: true })
        } catch (err) {
          console.error('[router-guard] 加载用户信息或动态路由失败:', err)
          userStore.logout()
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
