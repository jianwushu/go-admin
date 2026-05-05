<template>
  <div style="padding: 16px;">
    <!-- 查询栏 -->
    <el-card shadow="never" style="margin-bottom: 16px;">
      <el-form :model="queryParams" inline>
{{- range .Fields}}
{{- if .IsQuery}}
        <el-form-item label="{{.Label}}">
{{- if eq .HtmlType "select"}}
          <el-select v-model="queryParams.{{firstLower .GoField}}" placeholder="请选择{{.Label}}" clearable style="width: 120px">
            <el-option label="全部" value="" />
          </el-select>
{{- else if eq .HtmlType "date"}}
          <el-date-picker v-model="queryParams.{{firstLower .GoField}}" type="date" placeholder="选择{{.Label}}" />
{{- else}}
          <el-input v-model="queryParams.{{firstLower .GoField}}" placeholder="请输入{{.Label}}" clearable style="width: 200px" @keyup.enter="handleQuery" />
{{- end}}
        </el-form-item>
{{- end}}
{{- end}}
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon style="margin-right: 4px;"><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon style="margin-right: 4px;"><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" body-style="padding: 20px;">
      <!-- 表格工具栏 -->
      <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
        <span style="font-weight: 500; font-size: 18px;">{{.BusinessName}}列表</span>
        <div style="display: flex; align-items: center; gap: 8px;">
          <el-button type="primary" v-permission="['{{.ModuleName}}:{{.BusinessName}}:add']" @click="handleAdd">
            <el-icon style="margin-right: 4px;"><Plus /></el-icon>
            新增
          </el-button>
          <el-button type="danger" v-permission="['{{.ModuleName}}:{{.BusinessName}}:remove']" :disabled="selectedIds.length === 0" @click="handleBatchDelete">
            <el-icon style="margin-right: 4px;"><Delete /></el-icon>
            批量删除
          </el-button>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%; margin-bottom: 16px;"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
{{- range .Fields}}
{{- if .IsList}}
        <el-table-column prop="{{firstLower .GoField}}" label="{{.Label}}" min-width="120" show-overflow-tooltip />
{{- end}}
{{- end}}
        <el-table-column prop="createdAt" label="创建时间" width="170" align="center" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link v-permission="['{{.ModuleName}}:{{.BusinessName}}:edit']" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" link v-permission="['{{.ModuleName}}:{{.BusinessName}}:remove']" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <Pagination
        v-model:page="queryParams.page"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        @change="handleQuery"
      />
    </el-card>

    <!-- 表单弹窗 -->
    <{{.ClassName}}Form ref="formRef" @success="handleQuery" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { get{{.ClassName}}List, delete{{.ClassName}} } from '@/api/{{.ModuleName}}/{{.BusinessName}}'
import {{.ClassName}}Form from './components/{{.ClassName}}Form.vue'
import Pagination from '@/components/Pagination.vue'

const { t } = useI18n()
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const formRef = ref()
const selectedIds = ref<number[]>([])

const queryParams = reactive({
  page: 1,
  pageSize: 10,
{{- range .Fields}}
{{- if .IsQuery}}
  {{firstLower .GoField}}: {{- if eq .GoType "string"}}undefined{{- else if or (eq .GoType "int") (eq .GoType "int64")}}undefined{{- else}}undefined{{- end}},
{{- end}}
{{- end}}
})

/** 获取数据 */
const handleQuery = async () => {
  loading.value = true
  try {
    const res = await get{{.ClassName}}List(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

/** 重置查询 */
const handleReset = () => {
{{- range .Fields}}
{{- if .IsQuery}}
  queryParams.{{firstLower .GoField}} = undefined
{{- end}}
{{- end}}
  queryParams.page = 1
  handleQuery()
}

/** 多选变化 */
const handleSelectionChange = (rows: any[]) => {
  selectedIds.value = rows.map(row => row.id)
}

/** 新增 */
const handleAdd = () => {
  formRef.value?.open()
}

/** 编辑 */
const handleEdit = (row: any) => {
  formRef.value?.open(row.id)
}

/** 删除 */
const handleDelete = (row: any) => {
  ElMessageBox.confirm('确认删除该记录？', '提示', {
    type: 'warning',
  }).then(async () => {
    await delete{{.ClassName}}(row.id)
    ElMessage.success('删除成功')
    handleQuery()
  })
}

/** 批量删除 */
const handleBatchDelete = () => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请至少选择一条记录')
    return
  }
  ElMessageBox.confirm('确认删除所选记录？', '提示', {
    type: 'warning',
  }).then(async () => {
    await delete{{.ClassName}}(selectedIds.value)
    ElMessage.success('批量删除成功')
    handleQuery()
  })
}

onMounted(() => {
  handleQuery()
})
</script>
