package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-admin/middleware"
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// SystemConfigController 系统配置控制器
type SystemConfigController struct {
	service *service.SystemConfigService
}

// NewSystemConfigController 创建系统配置控制器实例
func NewSystemConfigController() *SystemConfigController {
	return &SystemConfigController{
		service: service.NewSystemConfigService(),
	}
}

// GetAll 获取所有配置
// @Summary 获取所有系统配置
// @Description 获取所有系统配置项
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.SystemConfigResponse}
// @Router /api/v1/system/config/all [get]
func (ctrl *SystemConfigController) GetAll(c *gin.Context) {
	configs, err := ctrl.service.GetAll()
	if err != nil {
		utils.Fail(c, 500, "查询配置失败："+err.Error())
		return
	}
	utils.Success(c, configs)
}

// GetList 获取配置列表（分页）
// @Summary 获取系统配置列表
// @Description 分页查询系统配置
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param configKey query string false "配置键"
// @Param configType query string false "配置类型"
// @Success 200 {object} response.PageResponse{data=[]response.SystemConfigResponse}
// @Router /api/v1/system/config/list [get]
func (ctrl *SystemConfigController) GetList(c *gin.Context) {
	var req request.SystemConfigListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	configs, total, err := ctrl.service.GetList(req)
	if err != nil {
		utils.Fail(c, 500, "查询配置失败："+err.Error())
		return
	}

	utils.PageSuccess(c, configs, total, req.GetPage(), req.GetPageSize())
}

// GetByKey 根据配置键获取配置
// @Summary 根据键获取配置
// @Description 根据配置键获取单个配置值
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Param key path string true "配置键"
// @Success 200 {object} response.Response{data=response.SystemConfigResponse}
// @Router /api/v1/system/config/key/{key} [get]
func (ctrl *SystemConfigController) GetByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.Fail(c, 400, "配置键不能为空")
		return
	}

	config, err := ctrl.service.GetByKey(key)
	if err != nil {
		utils.Fail(c, 500, "查询配置失败："+err.Error())
		return
	}

	utils.Success(c, config)
}

// GetByKeys 批量获取配置
// @Summary 批量获取配置
// @Description 根据多个配置键批量获取配置值
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Param keys query string true "配置键列表，逗号分隔"
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/system/config/keys [get]
func (ctrl *SystemConfigController) GetByKeys(c *gin.Context) {
	keysStr := c.Query("keys")
	if keysStr == "" {
		utils.Fail(c, 400, "配置键不能为空")
		return
	}

	keys := strings.Split(keysStr, ",")
	result, err := ctrl.service.GetByKeys(keys)
	if err != nil {
		utils.Fail(c, 500, "查询配置失败："+err.Error())
		return
	}

	utils.Success(c, result)
}

// Update 更新单个配置
// @Summary 更新系统配置
// @Description 更新单个系统配置项
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Param data body request.SystemConfigUpdateRequest true "配置更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/system/config [put]
func (ctrl *SystemConfigController) Update(c *gin.Context) {
	var req request.SystemConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.service.Update(req, userID); err != nil {
		utils.Fail(c, 500, "更新配置失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// BatchUpdate 批量更新配置
// @Summary 批量更新系统配置
// @Description 批量更新系统配置项
// @Tags 系统配置
// @Security BearerAuth
// @Produce json
// @Param data body request.SystemConfigBatchUpdateRequest true "批量更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/system/config/batch [put]
func (ctrl *SystemConfigController) BatchUpdate(c *gin.Context) {
	var req request.SystemConfigBatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.service.BatchUpdate(req, userID); err != nil {
		utils.Fail(c, 500, "批量更新配置失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "批量更新成功", nil)
}

// UploadLogo 上传Logo图片
// @Summary 上传Logo
// @Description 上传系统Logo或Favicon图片
// @Tags 系统配置
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Param type query string true "配置类型 logo/favicon"
// @Success 200 {object} response.Response{data=string}
// @Router /api/v1/system/config/upload [post]
func (ctrl *SystemConfigController) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, 400, "请选择要上传的文件")
		return
	}

	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".ico": true}
	if !allowedExts[ext] {
		utils.Fail(c, 400, "不支持的文件类型，仅支持 png/jpg/jpeg/gif/svg/ico")
		return
	}

	// 校验文件大小（最大2MB）
	if file.Size > 2*1024*1024 {
		utils.Fail(c, 400, "文件大小不能超过2MB")
		return
	}

	// 创建上传目录
	uploadDir := "uploads/logo"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Fail(c, 500, "创建上传目录失败："+err.Error())
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.Fail(c, 500, "保存文件失败："+err.Error())
		return
	}

	// 返回可访问的URL路径
	urlPath := "/" + strings.ReplaceAll(savePath, "\\", "/")
	utils.Success(c, urlPath)
}
