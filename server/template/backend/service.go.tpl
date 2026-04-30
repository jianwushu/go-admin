package service

import (
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
	"go-admin/utils"
)

// {{.ClassName}}Service {{.FunctionName}}服务
type {{.ClassName}}Service struct {
	repo *repository.{{.ClassName}}Repository
}

// New{{.ClassName}}Service 创建{{.FunctionName}}服务实例
func New{{.ClassName}}Service() *{{.ClassName}}Service {
	return &{{.ClassName}}Service{
		repo: repository.New{{.ClassName}}Repository(),
	}
}

// GetList 获取{{.FunctionName}}列表（分页）
func (s *{{.ClassName}}Service) GetList(req request.{{.ClassName}}ListRequest, scopeInfo *utils.DataScopeInfo) ([]response.{{.ClassName}}Response, int64, error) {
	items, total, err := s.repo.FindWithPage(
		req.GetPage(),
		req.GetPageSize(),
		scopeInfo,
{{- range .Fields}}
{{- if .IsQuery}}
		req.{{.GoField}},
{{- end}}
{{- end}}
	)
	if err != nil {
		return nil, 0, err
	}

	var result []response.{{.ClassName}}Response
	for _, item := range items {
		result = append(result, s.toResponse(&item))
	}

	return result, total, nil
}

// GetByID 根据ID获取详情
func (s *{{.ClassName}}Service) GetByID(id int64) (*response.{{.ClassName}}Response, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := s.toResponse(item)
	return &resp, nil
}

// Create 创建{{.FunctionName}}
func (s *{{.ClassName}}Service) Create(req request.{{.ClassName}}CreateRequest) error {
	item := entity.{{.ClassName}}{
{{- range .Fields}}
{{- if not .IsPk}}
		{{.GoField}}: req.{{.GoField}},
{{- end}}
{{- end}}
	}
	return s.repo.Create(&item)
}

// Update 更新{{.FunctionName}}
func (s *{{.ClassName}}Service) Update(req request.{{.ClassName}}UpdateRequest) error {
	item, err := s.repo.FindByID(req.ID)
	if err != nil {
		return err
	}

{{- range .Fields}}
{{- if and (not .IsPk) .IsEdit}}
	item.{{.GoField}} = req.{{.GoField}}
{{- end}}
{{- end}}

	return s.repo.Update(item)
}

// Delete 删除{{.FunctionName}}
func (s *{{.ClassName}}Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

// toResponse 实体转响应
func (s *{{.ClassName}}Service) toResponse(item *entity.{{.ClassName}}) response.{{.ClassName}}Response {
	return response.{{.ClassName}}Response{
		ID: item.ID,
{{- range .Fields}}
{{- if .IsList}}
		{{.GoField}}: item.{{.GoField}},
{{- end}}
{{- end}}
		CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
