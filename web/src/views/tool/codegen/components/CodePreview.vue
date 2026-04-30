<template>
  <div class="code-preview">
    <div v-if="Object.keys(previewData).length === 0" class="flex h-64 items-center justify-center text-gray-400">
      <div class="text-center">
        <el-icon :size="48" class="mb-4"><Document /></el-icon>
        <p>{{ t('codegen.noPreview') }}</p>
        <p class="text-sm">{{ t('codegen.previewTip') }}</p>
      </div>
    </div>
    <template v-else>
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane
          v-for="(code, filename) in previewData"
          :key="filename"
          :label="filename"
          :name="filename"
        >
          <div class="code-container">
            <div class="code-header flex items-center justify-between bg-gray-50 px-4 py-2 dark:bg-gray-700">
              <span class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ filename }}</span>
              <el-button type="primary" link size="small" @click="copyCode(code)">
                <el-icon class="mr-1"><CopyDocument /></el-icon>
                {{ t('common.copy') }}
              </el-button>
            </div>
            <pre class="code-block"><code v-html="highlightedCode(code, filename)"></code></pre>
          </div>
        </el-tab-pane>
      </el-tabs>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import sql from 'highlight.js/lib/languages/sql'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import 'highlight.js/styles/github.css'

// 注册语言
hljs.registerLanguage('go', go)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('xml', xml)

const { t } = useI18n()

const props = defineProps<{
  previewData: Record<string, string>
}>()

const activeTab = ref('')

// 当 previewData 变化时，设置第一个 tab 为活跃
watch(() => props.previewData, (data) => {
  const keys = Object.keys(data)
  if (keys.length > 0 && !activeTab.value) {
    activeTab.value = keys[0]
  }
}, { immediate: true })

/** 根据文件名获取语言 */
function getLanguage(filename: string): string {
  if (filename.endsWith('.go')) return 'go'
  if (filename.endsWith('.sql')) return 'sql'
  if (filename.endsWith('.ts') || filename.endsWith('.tsx')) return 'typescript'
  if (filename.endsWith('.vue') || filename.endsWith('.html') || filename.endsWith('.xml')) return 'xml'
  return 'plaintext'
}

/** 高亮代码 */
function highlightedCode(code: string, filename: string): string {
  const language = getLanguage(filename)
  try {
    if (hljs.getLanguage(language)) {
      return hljs.highlight(code, { language }).value
    }
  } catch {
    // 忽略高亮错误
  }
  return code
}

/** 复制代码 */
async function copyCode(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    ElMessage.success(t('common.copySuccess'))
  } catch {
    // 降级方案
    const textarea = document.createElement('textarea')
    textarea.value = code
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success(t('common.copySuccess'))
  }
}
</script>

<style scoped>
.code-container {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  overflow: hidden;
}

.code-block {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.6;
  font-family: 'Monaco', 'Menlo', 'Consolas', 'Courier New', monospace;
  background-color: #fafafa;
  max-height: 500px;
}

:root.dark .code-block {
  background-color: #1e1e1e;
  color: #d4d4d4;
}

:deep(.el-tabs__content) {
  padding: 0;
}
</style>
