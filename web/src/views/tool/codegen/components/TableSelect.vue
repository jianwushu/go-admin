<template>
  <div class="table-select">
    <el-select
      v-model="selectedTable"
      :placeholder="t('codegen.tableSelectPlaceholder')"
      filterable
      style="width: 100%"
      @change="handleTableChange"
    >
      <el-option
        v-for="table in tables"
        :key="table.tableName"
        :label="`${table.tableName} - ${table.tableComment || ''}`"
        :value="table.tableName"
      />
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getTableList } from '@/api/codegen'
import type { TableInfo } from '@/api/codegen'

const { t } = useI18n()

const emit = defineEmits<{
  (e: 'select', table: TableInfo): void
}>()

const selectedTable = ref('')
const tables = ref<TableInfo[]>([])
const loading = ref(false)

/** 加载表列表 */
async function loadTables() {
  loading.value = true
  try {
    const res = await getTableList()
    tables.value = res.data.data || []
  } catch (error) {
    console.error('获取表列表失败:', error)
  } finally {
    loading.value = false
  }
}

/** 选择表 */
function handleTableChange(tableName: string) {
  const table = tables.value.find(t => t.tableName === tableName)
  if (table) {
    emit('select', table)
  }
}

/** 清空选择 */
function clearSelection() {
  selectedTable.value = ''
}

onMounted(() => {
  loadTables()
})

defineExpose({ clearSelection })
</script>
