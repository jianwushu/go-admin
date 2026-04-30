import { createRouter, createWebHistory } from 'vue-router'
import { staticRoutes } from './staticRoutes'

const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes,
  scrollBehavior: () => ({ top: 0 }),
})

export default router
