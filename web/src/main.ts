import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import pinia from './store'
import i18n from './i18n'
import { setupPermissionDirective } from './directives/permission'
import { initTheme, initPrimaryColor } from './utils/theme'

// Element-Plus 样式
import 'element-plus/dist/index.css'
// Element-Plus 暗黑主题 CSS
import 'element-plus/theme-chalk/dark/css-vars.css'

// 全局样式
import '@/assets/styles/index.css'

// 初始化主题（必须在 DOM 渲染前执行）
initTheme()
// 初始化主题色
initPrimaryColor()

// 路由守卫
import './router/guard'

const app = createApp(App)

app.use(router)
app.use(pinia)
app.use(i18n)

// 注册全局指令
setupPermissionDirective(app)

app.mount('#app')
