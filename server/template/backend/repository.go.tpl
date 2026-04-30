package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/utils"
)

// {{.ClassName}}Repository {{.FunctionName}}仓储层
type {{.ClassName}}Repository struct {
	BaseRepository
}

// New{{.ClassName}}Repository 创建{{.FunctionName}}仓储实例
func New{{.ClassName}}Repository() *{{.ClassName}}Repository {
	return &{{.ClassName}}Repository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindByID 根据ID查询
func (r *{{.ClassName}}Repository) FindByID(id int64) (*entity.{{.ClassName}}, error) {
	var item entity.{{.ClassName}}
	err := r.DB.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindWithPage 分页查询
func (r *{{.ClassName}}Repository) FindWithPage(page, pageSize int, scopeInfo *utils.DataScopeInfo{{- range .Fields}}{{- if .IsQuery}}, {{.GoField}} {{.GoType}}{{- end}}{{- end}}) ([]entity.{{.ClassName}}, int64, error) {
	var items []entity.{{.ClassName}}
	var total int64

	query := r.DB.Model(&entity.{{.ClassName}}{})

	// 应用数据权限过滤
	if scopeInfo != nil {
		query = r.ApplyDataScope(query, scopeInfo, "", "dept_id")
	}

	// 条件过滤
{{- range .Fields}}
{{- if .IsQuery}}
{{- if eq .GoType "string"}}
	if {{.GoField}} != "" {
		{{- if eq .QueryType "LIKE"}}
		query = query.Where("{{.ColumnName}} LIKE ?", "%"+{{.GoField}}+"%")
		{{- else}}
		query = query.Where("{{.ColumnName}} {{.QueryType}} ?", {{.GoField}})
		{{- end}}
	}
{{- else if or (eq .GoType "int") (eq .GoType "int64")}}
	if {{.GoField}} > 0 {
		query = query.Where("{{.ColumnName}} {{.QueryType}} ?", {{.GoField}})
	}
{{- end}}
{{- end}}
{{- end}}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Create 创建
func (r *{{.ClassName}}Repository) Create(item *entity.{{.ClassName}}) error {
	return r.DB.Create(item).Error
}

// Update 更新
func (r *{{.ClassName}}Repository) Update(item *entity.{{.ClassName}}) error {
	return r.DB.Save(item).Error
}

// Delete 删除
func (r *{{.ClassName}}Repository) Delete(id int64) error {
	return r.DB.Delete(&entity.{{.ClassName}}{}, id).Error
}
