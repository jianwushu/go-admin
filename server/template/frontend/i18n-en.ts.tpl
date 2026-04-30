export default {
  {{.BusinessName}}: {
    // List page
    list: '{{.FunctionName}} List',
    add: 'Add {{.FunctionName}}',
    edit: 'Edit {{.FunctionName}}',
    delete: 'Delete {{.FunctionName}}',
    deleteConfirm: 'Are you sure to delete this {{.FunctionName}} record?',
    // Field labels
{{- range .Fields}}
{{- if .IsList}}
    {{firstLower .GoField}}: '{{.GoField}}',
{{- end}}
{{- end}}
    createdAt: 'Created At',
    operation: 'Operation',
  },
}
