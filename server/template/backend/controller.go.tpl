package controller

import (
	"strconv"

	"go-admin/middleware"
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// {{.ClassName}}Controller {{.FunctionName}}控制器
type {{.ClassName}}Controller struct {
	service *service.{{.ClassName}}Service
}

// New{{.ClassName}}Controller 创建{{.FunctionName}}控制器实例
func New{{.ClassName}}Controller() *{{.ClassName}}Controller {
	return &{{.ClassName}}Controller{
		service: service.New{{.ClassName}}Service(),
	}
}

// GetList 获取{{.FunctionName}}列表
func (ctrl *{{.ClassName}}Controller) GetList(c *gin.Context) {
	var req request.{{.ClassName}}ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	scopeInfo := middleware.GetDataScopeFromContext(c)

	items, total, err := ctrl.service.GetList(req, scopeInfo)
	if err != nil {
		utils.Fail(c, 500, "查询列表失败："+err.Error())
		return
	}

	utils.PageSuccess(c, items, total, req.GetPage(), req.GetPageSize())
}

// GetByID 根据ID获取详情
func (ctrl *{{.ClassName}}Controller) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的ID")
		return
	}

	item, err := ctrl.service.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.Success(c, item)
}

// Create 创建{{.FunctionName}}
func (ctrl *{{.ClassName}}Controller) Create(c *gin.Context) {
	var req request.{{.ClassName}}CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.service.Create(req); err != nil {
		utils.Fail(c, 500, "创建失败："+err.Error())
		return
	}

	utils.Success(c, nil)
}

// Update 更新{{.FunctionName}}
func (ctrl *{{.ClassName}}Controller) Update(c *gin.Context) {
	var req request.{{.ClassName}}UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.service.Update(req); err != nil {
		utils.Fail(c, 500, "更新失败："+err.Error())
		return
	}

	utils.Success(c, nil)
}

// Delete 删除{{.FunctionName}}
func (ctrl *{{.ClassName}}Controller) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的ID")
		return
	}

	if err := ctrl.service.Delete(id); err != nil {
		utils.Fail(c, 500, "删除失败："+err.Error())
		return
	}

	utils.Success(c, nil)
}
