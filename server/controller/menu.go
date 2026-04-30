package controller

import (
	"strconv"

	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *service.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController() *MenuController {
	return &MenuController{
		menuService: service.NewMenuService(),
	}
}

// GetTree 获取菜单树形列表
// @Summary 获取菜单树形列表
// @Description 获取所有菜单的树形结构
// @Tags 菜单管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.MenuResponse}
// @Router /api/v1/system/menu/tree [get]
func (ctrl *MenuController) GetTree(c *gin.Context) {
	tree, err := ctrl.menuService.GetTree()
	if err != nil {
		utils.Fail(c, 500, "查询菜单树失败："+err.Error())
		return
	}

	utils.Success(c, tree)
}

// GetByID 根据ID获取菜单详情
// @Summary 获取菜单详情
// @Description 根据菜单ID获取菜单详细信息
// @Tags 菜单管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response{data=response.MenuResponse}
// @Router /api/v1/system/menu/{id} [get]
func (ctrl *MenuController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的菜单ID")
		return
	}

	menu, err := ctrl.menuService.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.Success(c, menu)
}

// Create 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.MenuCreateRequest true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/menu [post]
func (ctrl *MenuController) Create(c *gin.Context) {
	var req request.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.menuService.Create(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.MenuUpdateRequest true "菜单信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/menu [put]
func (ctrl *MenuController) Update(c *gin.Context) {
	var req request.MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.menuService.Update(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除菜单
// @Summary 删除菜单
// @Description 根据菜单ID删除菜单
// @Tags 菜单管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/system/menu/{id} [delete]
func (ctrl *MenuController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的菜单ID")
		return
	}

	if err := ctrl.menuService.Delete(id); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
