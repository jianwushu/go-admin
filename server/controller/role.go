package controller

import (
	"strconv"

	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *service.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController() *RoleController {
	return &RoleController{
		roleService: service.NewRoleService(),
	}
}

// GetList 获取角色列表
// @Summary 获取角色列表
// @Description 分页查询角色列表，支持按名称、标识、状态筛选
// @Tags 角色管理
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param name query string false "角色名称"
// @Param code query string false "角色标识"
// @Param status query int false "状态"
// @Success 200 {object} response.PageResponse{data=[]response.RoleResponse}
// @Router /api/v1/system/role/list [get]
func (ctrl *RoleController) GetList(c *gin.Context) {
	var req request.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	roles, total, err := ctrl.roleService.GetList(req)
	if err != nil {
		utils.Fail(c, 500, "查询角色列表失败："+err.Error())
		return
	}

	utils.PageSuccess(c, roles, total, req.GetPage(), req.GetPageSize())
}

// GetAll 获取所有角色（不分页，用于下拉选择）
// @Summary 获取所有角色
// @Description 获取所有角色列表（不分页）
// @Tags 角色管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.RoleResponse}
// @Router /api/v1/system/role/all [get]
func (ctrl *RoleController) GetAll(c *gin.Context) {
	roles, err := ctrl.roleService.GetAll()
	if err != nil {
		utils.Fail(c, 500, "查询角色列表失败："+err.Error())
		return
	}

	utils.Success(c, roles)
}

// GetByID 根据ID获取角色详情
// @Summary 获取角色详情
// @Description 根据角色ID获取角色详细信息
// @Tags 角色管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=response.RoleResponse}
// @Router /api/v1/system/role/{id} [get]
func (ctrl *RoleController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的角色ID")
		return
	}

	role, err := ctrl.roleService.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.Success(c, role)
}

// Create 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.RoleCreateRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/role [post]
func (ctrl *RoleController) Create(c *gin.Context) {
	var req request.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.roleService.Create(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.RoleUpdateRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/role [put]
func (ctrl *RoleController) Update(c *gin.Context) {
	var req request.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.roleService.Update(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除角色
// @Summary 删除角色
// @Description 根据角色ID删除角色
// @Tags 角色管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/v1/system/role/{id} [delete]
func (ctrl *RoleController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的角色ID")
		return
	}

	if err := ctrl.roleService.Delete(id); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ChangeStatus 修改角色状态
// @Summary 修改角色状态
// @Description 启用或禁用角色
// @Tags 角色管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.RoleChangeStatusRequest true "状态参数"
// @Success 200 {object} response.Response
// @Router /api/v1/system/role/change-status [put]
func (ctrl *RoleController) ChangeStatus(c *gin.Context) {
	var req request.RoleChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.roleService.ChangeStatus(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "状态修改成功", nil)
}
