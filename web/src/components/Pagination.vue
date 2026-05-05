<template>
  <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="pageSizes"
      :total="total"
      :background="background"
      :layout="layout"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'Pagination' })

const props = withDefaults(defineProps<{
  /** 当前页码 */
  page: number
  /** 每页条数 */
  pageSize: number
  /** 总条数 */
  total: number
  /** 每页条数选项 */
  pageSizes?: number[]
  /** 是否显示背景色 */
  background?: boolean
  /** 布局配置 */
  layout?: string
}>(), {
  pageSizes: () => [10, 20, 50, 100],
  background: true,
  layout: 'total, sizes, prev, pager, next, jumper',
})

const emit = defineEmits<{
  'update:page': [value: number]
  'update:pageSize': [value: number]
  'change': []
}>()

const currentPage = computed({
  get: () => props.page,
  set: (val) => emit('update:page', val),
})

const pageSize = computed({
  get: () => props.pageSize,
  set: (val) => emit('update:pageSize', val),
})

function handleSizeChange() {
  emit('update:page', 1)
  emit('change')
}

function handleCurrentChange() {
  emit('change')
}
</script>
