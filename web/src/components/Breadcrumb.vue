<template>
  <el-breadcrumb separator="/" class="breadcrumb">
    <el-breadcrumb-item
      v-for="(item, index) in breadcrumbs"
      :key="item.path"
    >
      <span
        v-if="index === breadcrumbs.length - 1"
        class="breadcrumb-text no-redirect"
      >
        {{ getBreadcrumbTitle(item) }}
      </span>
      <router-link v-else class="breadcrumb-text" :to="(item.redirect as string) || item.path">
        {{ getBreadcrumbTitle(item) }}
      </router-link>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, type RouteLocationMatched } from 'vue-router'
import { useI18n } from 'vue-i18n'

defineOptions({ name: 'Breadcrumb' })

const route = useRoute()
const { t } = useI18n()

const breadcrumbs = ref<RouteLocationMatched[]>([])

/** 获取面包屑标题（处理国际化） */
function getBreadcrumbTitle(item: RouteLocationMatched): string {
  const i18nKey = item.meta?.i18nKey as string
  const title = item.meta?.title as string

  console.log(item,i18nKey,1,title)
  if (!title) return ''
  // 优先使用 i18nKey 翻译
  if (i18nKey) {
    const translated = t(i18nKey)
    if (translated !== i18nKey) return translated
  }
  // 回退：尝试通过路由 name 翻译
  const name = typeof item.name === 'string' ? item.name : ''
  const translated = name ? t(`menu.${name}`, title) : title
  return translated !== `menu.${name}` ? translated : title
}

/** 获取面包屑数据 */
function getBreadcrumbs() {
  const matched = route.matched.filter((item) => item.meta?.title)
  // 首页始终显示
  const first = matched[0]
  if (!isDashboard(first)) {
    matched.unshift({
      path: '/dashboard',
      name: 'Dashboard',
      meta: { title: 'Dashboard', i18nKey: 'menu.dashboard' },
    } as RouteLocationMatched)
  }
  breadcrumbs.value = matched
}

/** 判断是否为首页 */
function isDashboard(route: RouteLocationMatched): boolean {
  const name = route.name as string
  return name?.trim().toLocaleLowerCase() === 'Dashboard'.toLocaleLowerCase()
}

watch(
  () => route.path,
  () => getBreadcrumbs(),
  { immediate: true }
)
</script>

<style scoped>
.breadcrumb {
  line-height: var(--navbar-height);
}

:deep(.el-breadcrumb__inner) {
  color: var(--el-text-color-placeholder);
}

:deep(.el-breadcrumb__inner a) {
  color: var(--el-text-color-secondary);
  font-weight: normal;
}

:deep(.el-breadcrumb__inner a:hover) {
  color: var(--el-color-primary);
}

.breadcrumb-text.no-redirect {
  color: var(--el-text-color-placeholder);
  cursor: text;
}
</style>
