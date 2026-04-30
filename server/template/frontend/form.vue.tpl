<template>
  <el-dialog v-model="visible" :title="isEdit ? '编辑{{.FunctionName}}' : '新增{{.FunctionName}}'" width="600px" destroy-on-close>
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
{{- range .Fields}}
{{- if and (not .IsPk) .IsEdit}}
      <el-form-item label="{{.Label}}" prop="{{firstLower .GoField}}">
{{- if eq .HtmlType "textarea"}}
        <el-input v-model="form.{{firstLower .GoField}}" type="textarea" placeholder="请输入{{.Label}}" />
{{- else if eq .HtmlType "select"}}
        <el-select v-model="form.{{firstLower .GoField}}" placeholder="请选择{{.Label}}">
          <el-option label="请选择" :value="0" />
        </el-select>
{{- else if eq .HtmlType "radio"}}
        <el-radio-group v-model="form.{{firstLower .GoField}}">
          <el-radio :label="1">是</el-radio>
          <el-radio :label="0">否</el-radio>
        </el-radio-group>
{{- else if eq .HtmlType "date"}}
        <el-date-picker v-model="form.{{firstLower .GoField}}" type="date" placeholder="选择{{.Label}}" />
{{- else if eq .HtmlType "datetime"}}
        <el-date-picker v-model="form.{{firstLower .GoField}}" type="datetime" placeholder="选择{{.Label}}" />
{{- else if eq .HtmlType "upload"}}
        <el-upload action="#" :auto-upload="false">
          <el-button size="small" type="primary">点击上传</el-button>
        </el-upload>
{{- else}}
        <el-input v-model="form.{{firstLower .GoField}}" placeholder="请输入{{.Label}}" />
{{- end}}
      </el-form-item>
{{- end}}
{{- end}}
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { get{{.ClassName}}ById, create{{.ClassName}}, update{{.ClassName}} } from '@/api/{{.ModuleName}}/{{.BusinessName}}'

const emit = defineEmits(['success'])

const visible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  id: 0,
{{- range .Fields}}
{{- if and (not .IsPk) .IsEdit}}
  {{firstLower .GoField}}: {{- if eq .GoType "string"}}''{{- else if or (eq .GoType "int") (eq .GoType "int64")}}0{{- else if eq .GoType "float64"}}0{{- else if eq .GoType "bool"}}false{{- else}}''{{- end}},
{{- end}}
{{- end}}
})

const rules: FormRules = {
{{- range .Fields}}
{{- if and (not .IsPk) .IsRequired}}
  {{firstLower .GoField}}: [{ required: true, message: '请输入{{.Label}}', trigger: 'blur' }],
{{- end}}
{{- end}}
}

// 打开弹窗
const open = async (id?: number) => {
  visible.value = true
  if (id) {
    isEdit.value = true
    const res = await get{{.ClassName}}ById(id)
    Object.assign(form, res.data)
  } else {
    isEdit.value = false
    resetForm()
  }
}

// 重置表单
const resetForm = () => {
  form.id = 0
{{- range .Fields}}
{{- if and (not .IsPk) .IsEdit}}
  form.{{firstLower .GoField}} = {{- if eq .GoType "string"}}''{{- else if or (eq .GoType "int") (eq .GoType "int64")}}0{{- else if eq .GoType "float64"}}0{{- else if eq .GoType "bool"}}false{{- else}}''{{- end}}
{{- end}}
{{- end}}
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await update{{.ClassName}}(form)
      ElMessage.success('更新成功')
    } else {
      await create{{.ClassName}}(form)
      ElMessage.success('创建成功')
    }
    visible.value = false
    emit('success')
  } finally {
    submitting.value = false
  }
}

defineExpose({ open })
</script>
