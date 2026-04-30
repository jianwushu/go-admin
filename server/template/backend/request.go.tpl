package request

// {{.ClassName}}ListRequest {{.FunctionName}}列表查询请求
type {{.ClassName}}ListRequest struct {
	PageRequest
{{- range .Fields}}
{{- if .IsQuery}}
	{{.GoField}} {{.GoType}} `json:"{{firstLower .GoField}}" form:"{{firstLower .GoField}}"` // {{.Label}}
{{- end}}
{{- end}}
}

// {{.ClassName}}CreateRequest 创建{{.FunctionName}}请求
type {{.ClassName}}CreateRequest struct {
{{- range .Fields}}
{{- if not .IsPk}}
	{{.GoField}} {{.GoType}} `json:"{{firstLower .GoField}}"{{- if .IsRequired}} binding:"required"{{- end}}` // {{.Label}}
{{- end}}
{{- end}}
}

// {{.ClassName}}UpdateRequest 更新{{.FunctionName}}请求
type {{.ClassName}}UpdateRequest struct {
	ID int64 `json:"id" binding:"required"` // ID
{{- range .Fields}}
{{- if and (not .IsPk) .IsEdit}}
	{{.GoField}} {{.GoType}} `json:"{{firstLower .GoField}}"{{- if .IsRequired}} binding:"required"{{- end}}` // {{.Label}}
{{- end}}
{{- end}}
}
