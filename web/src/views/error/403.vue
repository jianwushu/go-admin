<template>
  <div class="error-page">
    <h1 class="error-code">403</h1>
    <p class="error-desc">{{ t('error.noPermission') }}</p>
    <router-link to="/" class="error-link">{{ t('error.backHome') }}</router-link>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

defineOptions({ name: 'Forbidden' })

const { t } = useI18n()
</script>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: var(--el-bg-color-page);
  transition: background-color 0.3s ease;
}

.error-code {
  font-size: 120px;
  font-weight: 700;
  color: var(--el-text-color-placeholder);
  margin-bottom: 16px;
  line-height: 1;
}

.error-desc {
  font-size: 18px;
  color: var(--el-text-color-secondary);
  margin-bottom: 32px;
}

.error-link {
  display: inline-block;
  padding: 10px 24px;
  background: var(--el-color-primary);
  color: #fff;
  border-radius: 8px;
  text-decoration: none;
  font-size: 14px;
  transition: opacity 0.2s;
}

.error-link:hover {
  opacity: 0.85;
}
</style>
