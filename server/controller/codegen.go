package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// CodegenController 代码生成控制器
type CodegenController struct {
	codegenService *service.CodegenService
}

// NewCodegenController 创建代码生成控制器实例
func NewCodegenController() *CodegenController {
	return &CodegenController{
		codegenService: service.NewCodegenService(),
	}
}

// GetTables 获取数据库表列表
// @Summary 获取数据库表列表
// @Description 获取所有数据库表信息
// @Tags 代码生成
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.TableInfoResponse}
// @Router /api/v1/codegen/tables [get]
func (ctrl *CodegenController) GetTables(c *gin.Context) {
	tables, err := ctrl.codegenService.GetAllTables()
	if err != nil {
		utils.Fail(c, 500, "获取表列表失败："+err.Error())
		return
	}

	utils.Success(c, tables)
}

// GetColumns 获取表的列信息
// @Summary 获取表的列信息
// @Description 根据表名获取列信息
// @Tags 代码生成
// @Security BearerAuth
// @Produce json
// @Param tableName path string true "表名"
// @Success 200 {object} response.Response{data=[]response.ColumnInfoResponse}
// @Router /api/v1/codegen/columns/{tableName} [get]
func (ctrl *CodegenController) GetColumns(c *gin.Context) {
	tableName := c.Param("tableName")
	if tableName == "" {
		utils.Fail(c, 400, "参数错误：表名不能为空")
		return
	}

	columns, err := ctrl.codegenService.GetColumnsByTableName(tableName)
	if err != nil {
		utils.Fail(c, 500, "获取列信息失败："+err.Error())
		return
	}

	utils.Success(c, columns)
}

// Preview 代码预览
// @Summary 代码预览
// @Description 预览生成的代码
// @Tags 代码生成
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.CodegenPreviewRequest true "代码生成配置"
// @Success 200 {object} response.Response{data=response.CodegenPreviewResponse}
// @Router /api/v1/codegen/preview [post]
func (ctrl *CodegenController) Preview(c *gin.Context) {
	var req request.CodegenPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	result, err := ctrl.codegenService.PreviewCode(req)
	if err != nil {
		utils.Fail(c, 500, "代码预览失败："+err.Error())
		return
	}

	utils.Success(c, result)
}

// Generate 生成代码并下载
// @Summary 生成代码并下载
// @Description 生成代码并返回 ZIP 文件
// @Tags 代码生成
// @Security BearerAuth
// @Accept json
// @Param data body request.CodegenGenerateRequest true "代码生成配置"
// @Router /api/v1/codegen/generate [post]
func (ctrl *CodegenController) Generate(c *gin.Context) {
	var req request.CodegenGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	zipData, err := ctrl.codegenService.GenerateCode(req)
	if err != nil {
		utils.Fail(c, 500, "代码生成失败："+err.Error())
		return
	}

	// 设置响应头，返回 ZIP 文件
	fileName := fmt.Sprintf("%s_%s.zip", req.BusinessName, time.Now().Format("20060102150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Length", fmt.Sprintf("%d", len(zipData)))
	c.Data(http.StatusOK, "application/zip", zipData)
}

// SaveConfig 保存代码生成配置
// @Summary 保存代码生成配置
// @Description 保存代码生成配置
// @Tags 代码生成
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.CodegenSaveRequest true "配置信息"
// @Success 200 {object} response.Response
// @Router /api/v1/codegen/config [post]
func (ctrl *CodegenController) SaveConfig(c *gin.Context) {
	var req request.CodegenSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.codegenService.SaveConfig(req); err != nil {
		utils.Fail(c, 500, "保存配置失败："+err.Error())
		return
	}

	utils.Success(c, nil)
}

// GetConfig 根据表名获取配置
// @Summary 获取代码生成配置
// @Description 根据表名获取代码生成配置
// @Tags 代码生成
// @Security BearerAuth
// @Produce json
// @Param tableName path string true "表名"
// @Success 200 {object} response.Response{data=response.CodegenConfigResponse}
// @Router /api/v1/codegen/config/{tableName} [get]
func (ctrl *CodegenController) GetConfig(c *gin.Context) {
	tableName := c.Param("tableName")
	if tableName == "" {
		utils.Fail(c, 400, "参数错误：表名不能为空")
		return
	}

	config, err := ctrl.codegenService.GetConfigByTableName(tableName)
	if err != nil {
		utils.Fail(c, 404, "配置不存在："+err.Error())
		return
	}

	utils.Success(c, config)
}

// GetAllConfigs 获取所有配置
// @Summary 获取所有代码生成配置
// @Description 获取所有代码生成配置列表
// @Tags 代码生成
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.CodegenConfigResponse}
// @Router /api/v1/codegen/configs [get]
func (ctrl *CodegenController) GetAllConfigs(c *gin.Context) {
	configs, err := ctrl.codegenService.GetAllConfigs()
	if err != nil {
		utils.Fail(c, 500, "获取配置列表失败："+err.Error())
		return
	}

	utils.Success(c, configs)
}

// DeleteConfig 删除配置
// @Summary 删除代码生成配置
// @Description 根据ID删除代码生成配置
// @Tags 代码生成
// @Security BearerAuth
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response
// @Router /api/v1/codegen/config/{id} [delete]
func (ctrl *CodegenController) DeleteConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的ID")
		return
	}

	if err := ctrl.codegenService.DeleteConfig(id); err != nil {
		utils.Fail(c, 500, "删除配置失败："+err.Error())
		return
	}

	utils.Success(c, nil)
}
