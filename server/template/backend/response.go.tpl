package response

// {{.ClassName}}Response {{.FunctionName}}响应
type {{.ClassName}}Response struct {
	ID        int64  `json:"id"` // ID
{{- range .Fields}}
{{- if .IsList}}
	{{.GoField}} {{.GoType}} `json:"{{firstLower .GoField}}"` // {{.Label}}
{{- end}}
{{- end}}
	CreatedAt string `json:"createdAt"` // 创建时间
}
