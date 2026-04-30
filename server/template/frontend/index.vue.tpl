<template>
  <div class="p-4">
    <!-- 搜索区域 -->
    <el-card shadow="never" class="mb-4">
      <el-form :model="queryParams" inline>
{{- range .Fields}}
{{- if .IsQuery}}
        <el-form-item label="{{.Label}}">
{{- if eq .HtmlType "select"}}
          <el-select v-model="queryParams.{{firstLower .GoField}}" placeholder="请选择{{.Label}}" clearable>
            <el-option label="全部" value="" />
          </el-select>
{{- else if eq .HtmlType "date"}}
          <el-date-picker v-model="queryParams.{{firstLower .GoField}}" type="date" placeholder="选择{{.Label}}" />
{{- else}}
          <el-input v-model="queryParams.{{firstLower .GoField}}" placeholder="请输入{{.Label}}" clearable />
{{- end}}
        </el-form-item>
{{- end}}
{{- end}}
        <el-form-item>
          <el-button type="primary" @click="handleQuery">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="mb-4">
      <div class="flex justify-between items-center">
        <div>
          <el-button type="primary" v-permission="['{{.ModuleName}}:{{.BusinessName}}:add']" @click="handleAdd">
            <el-icon><Plus /></el-icon> 新增
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column type="index" label="#" width="50" align="center" />
{{- range .Fields}}
{{- if .IsList}}
        <el-table-column prop="{{firstLower .GoField}}" label="{{.Label}}" min-width="120" show-overflow-tooltip />
{{- end}}
{{- end}}
        <el-table-column prop="createdAt" label="创建时间" width="170" align="center" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" v-permission="['{{.ModuleName}}:{{.BusinessName}}:edit']" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" v-permission="['{{.ModuleName}}:{{.BusinessName}}:remove']" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleQuery"
          @current-change="handleQuery"
        />
      </div>
    </el-card>

    <!-- 表单弹窗 -->
    <{{.ClassName}}Form ref="formRef" @success="handleQuery" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { get{{.ClassName}}List, delete{{.ClassName}} } from '@/api/{{.ModuleName}}/{{.BusinessName}}'
import {{.ClassName}}Form from './components/{{.ClassName}}Form.vue'

const { t } = useI18n()
const loading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const formRef = ref()

const queryParams = reactive({
  page: 1,
  pageSize: 10,
{{- range .Fields}}
{{- if .IsQuery}}
  {{firstLower .GoField}}: {{- if eq .GoType "string"}}''{{- else if or (eq .GoType "int") (eq .GoType "int64")}}0{{- else}}''{{- end}},
{{- end}}
{{- end}}
})

// 查询列表
const handleQuery = async () => {
  loading.value = true
  try {
    const res = await get{{.ClassName}}List(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

// 重置查询
const handleReset = () => {
{{- range .Fields}}
{{- if .IsQuery}}
  queryParams.{{firstLower .GoField}} = {{- if eq .GoType "string"}}''{{- else if or (eq .GoType "int") (eq .GoType "int64")}}0{{- else}}''{{- end}}
{{- end}}
{{- end}}
  queryParams.page = 1
  handleQuery()
}

// 新增
const handleAdd = () => {
  formRef.value?.open()
}

// 编辑
const handleEdit = (row: any) => {
  formRef.value?.open(row.id)
}

// 删除
const handleDelete = (row: any) => {
  ElMessageBox.confirm('确认删除该记录？', '提示', {
    type: 'warning',
  }).then(async () => {
    await delete{{.ClassName}}(row.id)
    ElMessage.success('删除成功')
    handleQuery()
  })
}

onMounted(() => {
  handleQuery()
})
</script>
