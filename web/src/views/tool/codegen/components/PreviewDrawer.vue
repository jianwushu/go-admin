<template>
  <el-drawer
    v-model="visible"
    :title="t('codegen.codePreview')"
    size="80%"
    :before-close="handleClose"
    class="preview-drawer"
  >
    <template #default>
      <div v-if="loading" class="flex h-full items-center justify-center">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      </div>
      <CodePreview v-else :preview-data="previewData" />
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Loading } from '@element-plus/icons-vue'
import CodePreview from './CodePreview.vue'
import { previewCode } from '@/api/codegen'
import type { CodegenConfig } from '@/api/codegen'

const { t } = useI18n()

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  config: CodegenConfig | null
}>()

const loading = ref(false)
const previewData = ref<Record<string, string>>({})

watch(visible, async (val) => {
  if (val && props.config) {
    await loadPreview()
  }
})

watch(() => props.config, async () => {
  if (visible.value && props.config) {
    await loadPreview()
  }
})

async function loadPreview() {
  if (!props.config) return
  loading.value = true
  try {
    const res = await previewCode({
      tableName: props.config.tableName,
      tableComment: props.config.tableComment,
      className: props.config.className,
      businessName: props.config.businessName,
      functionName: props.config.functionName,
      moduleName: props.config.moduleName,
      packageName: props.config.packageName,
      author: props.config.author,
      fields: Array.isArray(props.config.fields) ? props.config.fields : [],
    })
    const files = res.data.data?.files || []
    const previewMap: Record<string, string> = {}
    for (const file of files) {
      previewMap[file.filePath] = file.content
    }
    previewData.value = previewMap
  } catch (error) {
    console.error('代码预览失败:', error)
    previewData.value = {}
  } finally {
    loading.value = false
  }
}

function handleClose() {
  visible.value = false
  previewData.value = {}
}
</script>

<style scoped>
:deep(.el-drawer__body) {
  padding: 0;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}
</style>
