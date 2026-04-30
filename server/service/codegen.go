package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"

	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/repository"
	"go-admin/utils"
)

// TemplateData 模板渲染数据
type TemplateData struct {
	TableName    string                 // 数据库表名
	TableComment string                 // 表注释
	ClassName    string                 // 大驼峰类名
	BusinessName string                 // 业务名（小驼峰）
	FunctionName string                 // 功能名（中文）
	ModuleName   string                 // 模块名
	PackageName  string                 // 包名
	Author       string                 // 作者
	CreateTime   string                 // 创建时间
	Fields       []request.ColumnConfig // 字段配置列表
}

// CodegenService 代码生成服务
type CodegenService struct {
	repo *repository.CodegenRepository
}

// NewCodegenService 创建代码生成服务实例
func NewCodegenService() *CodegenService {
	return &CodegenService{
		repo: repository.NewCodegenRepository(),
	}
}

// GetAllTables 获取所有数据库表
func (s *CodegenService) GetAllTables() ([]response.TableInfoResponse, error) {
	return s.repo.GetAllTables()
}

// GetColumnsByTableName 获取表的列信息
func (s *CodegenService) GetColumnsByTableName(tableName string) ([]response.ColumnInfoResponse, error) {
	return s.repo.GetColumnsByTableName(tableName)
}

// PreviewCode 代码预览
func (s *CodegenService) PreviewCode(req request.CodegenPreviewRequest) (*response.CodegenPreviewResponse, error) {
	data := s.buildTemplateData(req)

	files, err := s.renderAllTemplates(data)
	if err != nil {
		return nil, fmt.Errorf("模板渲染失败：%w", err)
	}

	return &response.CodegenPreviewResponse{Files: files}, nil
}

// GenerateCode 生成代码（返回 ZIP 文件字节）
func (s *CodegenService) GenerateCode(req request.CodegenGenerateRequest) ([]byte, error) {
	data := s.buildTemplateData(req)

	files, err := s.renderAllTemplates(data)
	if err != nil {
		return nil, fmt.Errorf("模板渲染失败：%w", err)
	}

	// 打包为 ZIP
	zipData, err := s.createZip(files)
	if err != nil {
		return nil, fmt.Errorf("ZIP 打包失败：%w", err)
	}

	return zipData, nil
}

// SaveConfig 保存代码生成配置
func (s *CodegenService) SaveConfig(req request.CodegenSaveRequest) error {
	fieldsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		return fmt.Errorf("序列化字段配置失败：%w", err)
	}

	config := &entity.CodegenConfig{
		TblName:      req.TableName,
		TableComment: req.TableComment,
		ClassName:    req.ClassName,
		BusinessName: req.BusinessName,
		FunctionName: req.FunctionName,
		ModuleName:   req.ModuleName,
		PackageName:  req.PackageName,
		Author:       req.Author,
		Fields:       string(fieldsJSON),
	}

	return s.repo.SaveConfig(config)
}

// GetConfigByTableName 根据表名获取配置
func (s *CodegenService) GetConfigByTableName(tableName string) (*response.CodegenConfigResponse, error) {
	config, err := s.repo.GetConfigByTableName(tableName)
	if err != nil {
		return nil, err
	}

	var fields []request.ColumnConfig
	if config.Fields != "" {
		if err := json.Unmarshal([]byte(config.Fields), &fields); err != nil {
			return nil, fmt.Errorf("反序列化字段配置失败：%w", err)
		}
	}

	return &response.CodegenConfigResponse{
		ID:           config.ID,
		TableName:    config.TblName,
		TableComment: config.TableComment,
		ClassName:    config.ClassName,
		BusinessName: config.BusinessName,
		FunctionName: config.FunctionName,
		ModuleName:   config.ModuleName,
		PackageName:  config.PackageName,
		Author:       config.Author,
		Fields:       toColumnConfigResponse(fields),
		CreatedAt:    config.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetAllConfigs 获取所有配置
func (s *CodegenService) GetAllConfigs() ([]response.CodegenConfigResponse, error) {
	configs, err := s.repo.GetAllConfigs()
	if err != nil {
		return nil, err
	}

	var result []response.CodegenConfigResponse
	for _, config := range configs {
		var fields []request.ColumnConfig
		if config.Fields != "" {
			_ = json.Unmarshal([]byte(config.Fields), &fields)
		}
		result = append(result, response.CodegenConfigResponse{
			ID:           config.ID,
			TableName:    config.TblName,
			TableComment: config.TableComment,
			ClassName:    config.ClassName,
			BusinessName: config.BusinessName,
			FunctionName: config.FunctionName,
			ModuleName:   config.ModuleName,
			PackageName:  config.PackageName,
			Author:       config.Author,
			Fields:       toColumnConfigResponse(fields),
			CreatedAt:    config.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// DeleteConfig 删除配置
func (s *CodegenService) DeleteConfig(id int64) error {
	return s.repo.DeleteConfig(id)
}

// buildTemplateData 构建模板数据
func (s *CodegenService) buildTemplateData(req request.CodegenPreviewRequest) TemplateData {
	return TemplateData{
		TableName:    req.TableName,
		TableComment: req.TableComment,
		ClassName:    req.ClassName,
		BusinessName: req.BusinessName,
		FunctionName: req.FunctionName,
		ModuleName:   req.ModuleName,
		PackageName:  req.PackageName,
		Author:       req.Author,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Fields:       req.Fields,
	}
}

// getTemplateFuncMap 获取模板函数映射
func (s *CodegenService) getTemplateFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	// 添加自定义函数
	funcMap["firstLower"] = utils.FirstLower
	funcMap["firstUpper"] = utils.FirstUpper
	funcMap["snakeToCamel"] = utils.SnakeToCamel
	funcMap["camelToSnake"] = utils.CamelToSnake
	return funcMap
}

// renderAllTemplates 渲染所有模板
func (s *CodegenService) renderAllTemplates(data TemplateData) ([]response.CodegenFile, error) {
	funcMap := s.getTemplateFuncMap()

	// 后端模板列表
	backendTemplates := map[string]string{
		"model.go.tpl":      "server/model/entity/%s.go",
		"request.go.tpl":    "server/model/request/%s.go",
		"response.go.tpl":   "server/model/response/%s.go",
		"repository.go.tpl": "server/repository/%s.go",
		"service.go.tpl":    "server/service/%s.go",
		"controller.go.tpl": "server/controller/%s.go",
		"router.go.tpl":     "server/router/%s.go",
		"sql.go.tpl":        "server/sql/%s.sql",
	}

	// 前端模板列表
	frontendTemplates := map[string]string{
		"api.ts.tpl":     "web/src/api/%s/%s.ts",
		"index.vue.tpl":  "web/src/views/%s/%s/index.vue",
		"form.vue.tpl":   "web/src/views/%s/%s/components/%sForm.vue",
		"i18n-zh.ts.tpl": "web/src/i18n/zh-CN/modules/%s.ts",
		"i18n-en.ts.tpl": "web/src/i18n/en/modules/%s.ts",
	}

	var files []response.CodegenFile

	// 渲染后端模板
	for tplName, pathPattern := range backendTemplates {
		content, err := s.renderTemplate(tplName, funcMap, data)
		if err != nil {
			return nil, fmt.Errorf("渲染后端模板 %s 失败：%w", tplName, err)
		}

		var filePath string
		if tplName == "sql.go.tpl" {
			filePath = fmt.Sprintf(pathPattern, data.BusinessName)
		} else {
			filePath = fmt.Sprintf(pathPattern, data.BusinessName)
		}

		files = append(files, response.CodegenFile{
			FileName: filepath.Base(filePath),
			FilePath: filePath,
			Content:  content,
		})
	}

	// 渲染前端模板
	for tplName, pathPattern := range frontendTemplates {
		content, err := s.renderTemplate(tplName, funcMap, data)
		if err != nil {
			return nil, fmt.Errorf("渲染前端模板 %s 失败：%w", tplName, err)
		}

		var filePath string
		switch tplName {
		case "api.ts.tpl":
			filePath = fmt.Sprintf(pathPattern, data.ModuleName, data.BusinessName)
		case "index.vue.tpl", "form.vue.tpl":
			filePath = fmt.Sprintf(pathPattern, data.ModuleName, data.BusinessName, data.ClassName)
		case "i18n-zh.ts.tpl", "i18n-en.ts.tpl":
			filePath = fmt.Sprintf(pathPattern, data.BusinessName)
		default:
			filePath = fmt.Sprintf(pathPattern, data.BusinessName)
		}

		files = append(files, response.CodegenFile{
			FileName: filepath.Base(filePath),
			FilePath: filePath,
			Content:  content,
		})
	}

	return files, nil
}

// renderTemplate 渲染单个模板
func (s *CodegenService) renderTemplate(tplName string, funcMap template.FuncMap, data TemplateData) (string, error) {
	// 读取模板文件
	tplPath := filepath.Join("template", getTemplateSubDir(tplName), tplName)

	tmpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(tplPath)
	if err != nil {
		return "", fmt.Errorf("解析模板文件 %s 失败：%w", tplPath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tplName, data); err != nil {
		return "", fmt.Errorf("执行模板 %s 失败：%w", tplName, err)
	}

	return buf.String(), nil
}

// getTemplateSubDir 获取模板子目录
func getTemplateSubDir(tplName string) string {
	if strings.HasSuffix(tplName, ".go.tpl") || strings.HasSuffix(tplName, ".sql") {
		return "backend"
	}
	return "frontend"
}

// createZip 创建 ZIP 文件
func (s *CodegenService) createZip(files []response.CodegenFile) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	for _, file := range files {
		f, err := w.Create(file.FilePath)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write([]byte(file.Content)); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// toColumnConfigResponse 转换字段配置响应
func toColumnConfigResponse(fields []request.ColumnConfig) []response.ColumnConfig {
	var result []response.ColumnConfig
	for _, f := range fields {
		result = append(result, response.ColumnConfig{
			ColumnName: f.ColumnName,
			ColumnType: f.ColumnType,
			GoType:     f.GoType,
			GoField:    f.GoField,
			TsType:     f.TsType,
			Label:      f.Label,
			HtmlType:   f.HtmlType,
			QueryType:  f.QueryType,
			IsList:     f.IsList,
			IsQuery:    f.IsQuery,
			IsRequired: f.IsRequired,
			IsEdit:     f.IsEdit,
			DictType:   f.DictType,
			Comment:    f.Comment,
			IsPk:       f.IsPk,
			IsNull:     f.IsNull,
			Sort:       f.Sort,
		})
	}
	return result
}
