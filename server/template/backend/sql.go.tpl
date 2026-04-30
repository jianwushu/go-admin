-- {{.TableComment}} ({{.TableName}})
CREATE TABLE IF NOT EXISTS `{{.TableName}}` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
{{- range .Fields}}
{{- if not .IsPk}}
  `{{.ColumnName}}` {{- if eq .GoType "string"}} VARCHAR({{- if gt .MaxLength 0}}{{.MaxLength}}{{- else}}256{{- end}}){{- else if eq .GoType "int64"}} BIGINT{{- else if eq .GoType "int"}} INTEGER{{- else if eq .GoType "float64"}} DOUBLE{{- else if eq .GoType "bool"}} TINYINT(1){{- else if eq .GoType "time.Time"}} DATETIME{{- else}} VARCHAR(256){{- end}}{{- if not .IsNull}} NOT NULL{{- end}}{{- if .Comment}} -- {{.Comment}}{{- end}},
{{- end}}
{{- end}}
  `created_by` BIGINT DEFAULT 0,
  `updated_by` BIGINT DEFAULT 0,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME
);

-- 索引
CREATE INDEX IF NOT EXISTS `idx_{{.TableName}}_deleted_at` ON `{{.TableName}}` (`deleted_at`);
