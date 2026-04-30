package entity

// {{.ClassName}} {{.TableComment}}
type {{.ClassName}} struct {
	BaseModel
{{- range .Fields}}
	{{.GoField}} {{.GoType}} `json:"{{firstLower .GoField}}" gorm:"{{- if .IsPk}}primaryKey;{{- end}}{{- if eq .GoType "string"}}size:256;{{- end}}{{- if not .IsNull}}not null;{{- end}}column:{{.ColumnName}}"`
{{- end}}
}

// TableName 返回带前缀的表名
func ({{.ClassName}}) TableName() string {
	return TableName("{{.BusinessName}}")
}
