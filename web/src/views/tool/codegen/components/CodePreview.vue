<template>
  <div class="code-preview">
    <div v-if="Object.keys(previewData).length === 0" class="flex h-64 items-center justify-center text-gray-400">
      <div class="text-center">
        <el-icon :size="48" class="mb-4"><Document /></el-icon>
        <p>{{ t('codegen.noPreview') }}</p>
        <p class="text-sm">{{ t('codegen.previewTip') }}</p>
      </div>
    </div>
    <div v-else class="preview-layout">
      <!-- 左侧文件树 -->
      <div class="file-tree">
        <div class="tree-header">
          <el-icon class="mr-1"><FolderOpened /></el-icon>
          <span>{{ t('codegen.fileList') }}</span>
        </div>
        <el-tree
          :data="treeData"
          :props="{ label: 'label', children: 'children' }"
          default-expand-all
          highlight-current
          :current-node-key="activeFilePath"
          node-key="path"
          @node-click="handleNodeClick"
        >
          <template #default="{ node, data }">
            <div class="tree-node">
              <el-icon v-if="data.isDir" class="mr-1 text-yellow-500"><Folder /></el-icon>
              <el-icon v-else class="mr-1 text-blue-400"><Document /></el-icon>
              <span>{{ node.label }}</span>
            </div>
          </template>
        </el-tree>
      </div>

      <!-- 右侧代码预览 -->
      <div class="code-panel">
        <div v-if="activeFilePath" class="code-container">
          <div class="code-header flex items-center justify-between bg-gray-50 px-4 py-2 dark:bg-gray-700">
            <span class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ activeFilePath }}</span>
            <el-button type="primary" link size="small" @click="copyCode(previewData[activeFilePath])">
              <el-icon class="mr-1"><CopyDocument /></el-icon>
              {{ t('common.copy') }}
            </el-button>
          </div>
          <pre class="code-block"><code v-html="highlightedCode(previewData[activeFilePath], activeFilePath)"></code></pre>
        </div>
        <div v-else class="flex h-full items-center justify-center text-gray-400">
          <div class="text-center">
            <el-icon :size="36" class="mb-2"><Document /></el-icon>
            <p class="text-sm">{{ t('codegen.selectFileToPreview') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Document, FolderOpened, Folder, CopyDocument } from '@element-plus/icons-vue'
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

const activeFilePath = ref('')

/** 构建树形结构数据 */
interface TreeNode {
  label: string
  path?: string
  isDir: boolean
  children?: TreeNode[]
}

const treeData = computed<TreeNode[]>(() => {
  const root: Record<string, any> = {}

  for (const filePath of Object.keys(props.previewData)) {
    const parts = filePath.split('/')
    let current = root

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      if (!current[part]) {
        current[part] = {
          __isDir: i < parts.length - 1,
          __path: i === parts.length - 1 ? filePath : undefined,
        }
      }
      if (i < parts.length - 1) {
        current[part].__isDir = true
        if (!current[part].__children) {
          current[part].__children = {}
        }
        current = current[part].__children
      }
    }
  }

  function buildTree(obj: Record<string, any>, parentPath = ''): TreeNode[] {
    const nodes: TreeNode[] = []
    const entries = Object.entries(obj)

    // 目录排前面，文件排后面
    const dirs = entries.filter(([, v]) => v.__isDir)
    const files = entries.filter(([, v]) => !v.__isDir)

    for (const [name, val] of dirs) {
      nodes.push({
        label: name,
        isDir: true,
        children: buildTree(val.__children || {}, `${parentPath}/${name}`),
      })
    }

    for (const [name, val] of files) {
      nodes.push({
        label: name,
        path: val.__path,
        isDir: false,
      })
    }

    return nodes
  }

  return buildTree(root)
})

// 当 previewData 变化时，自动选中第一个文件或重置无效路径
watch(() => props.previewData, (data) => {
  const keys = Object.keys(data)
  if (keys.length > 0) {
    if (!activeFilePath.value || !data[activeFilePath.value]) {
      activeFilePath.value = keys[0]
    }
  } else {
    activeFilePath.value = ''
  }
}, { immediate: true })

/** 点击文件节点 */
function handleNodeClick(data: TreeNode) {
  if (!data.isDir && data.path) {
    activeFilePath.value = data.path
  }
}

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
.preview-layout {
  display: flex;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  overflow: hidden;
  min-height: 500px;
}

.file-tree {
  width: 260px;
  min-width: 200px;
  border-right: 1px solid var(--el-border-color-lighter);
  background-color: var(--el-bg-color);
  overflow-y: auto;
}

.tree-header {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color-lighter);
  background-color: var(--el-fill-color-light);
}

.tree-node {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.code-panel {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.code-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.code-header {
  flex-shrink: 0;
}

.code-block {
  margin: 0;
  padding: 16px;
  overflow: auto;
  font-size: 13px;
  line-height: 1.6;
  font-family: 'Monaco', 'Menlo', 'Consolas', 'Courier New', monospace;
  background-color: #fafafa;
  flex: 1;
  max-height: 500px;
}

:root.dark .code-block {
  background-color: #1e1e1e;
  color: #d4d4d4;
}

:deep(.el-tree-node__content) {
  height: 32px;
}

:deep(.el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content) {
  background-color: var(--el-color-primary-light-9);
}
</style>
