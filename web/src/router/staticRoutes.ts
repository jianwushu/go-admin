import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/layout/index.vue'

/** 静态路由（无需权限） */
export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      hidden: true,
    },
  },
  {
    path: '/redirect',
    component: Layout,
    meta: {
      hidden: true,
    },
    children: [
      {
        path: '/redirect/:path(.*)',
        name: 'Redirect',
        component: () => import('@/views/redirect/index.vue'),
        meta: {
          title: 'Redirect',
          hidden: true,
        },
      },
    ],
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: 'Dashboard',
          i18nKey: 'menu.dashboard',
          icon: 'House',
          affix: true,
        },
      },
    ],
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404',
      hidden: true,
    },
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: {
      title: '403',
      hidden: true,
    },
  },
  {
    path: '/profile',
    component: Layout,
    meta: {
      hidden: true,
    },
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          title: '个人中心',
          i18nKey: 'menu.profile',
          hidden: true,
        },
      },
    ],
  },
  {
    path: '/tool',
    component: Layout,
    meta: {
      hidden: true,
    },
    children: [
      {
        path: 'codegen/create',
        name: 'CodegenCreate',
        component: () => import('@/views/tool/codegen/create.vue'),
        meta: {
          title: '新建配置',
          i18nKey: 'codegen.createConfig',
          hidden: true,
        },
      },
      {
        path: 'codegen/edit/:id',
        name: 'CodegenEdit',
        component: () => import('@/views/tool/codegen/edit.vue'),
        meta: {
          title: '编辑配置',
          i18nKey: 'codegen.editConfig',
          hidden: true,
        },
      },
    ],
  },
]
