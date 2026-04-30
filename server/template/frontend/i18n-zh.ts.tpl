export default {
  {{.BusinessName}}: {
    // 列表页
    list: '{{.FunctionName}}列表',
    add: '新增{{.FunctionName}}',
    edit: '编辑{{.FunctionName}}',
    delete: '删除{{.FunctionName}}',
    deleteConfirm: '确认删除该{{.FunctionName}}记录？',
    // 字段标签
{{- range .Fields}}
{{- if .IsList}}
    {{firstLower .GoField}}: '{{.Label}}',
{{- end}}
{{- end}}
    createdAt: '创建时间',
    operation: '操作',
  },
}
