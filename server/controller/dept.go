package controller

import (
	"strconv"

	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// DeptController 部门控制器
type DeptController struct {
	deptService *service.DeptService
}

// NewDeptController 创建部门控制器实例
func NewDeptController() *DeptController {
	return &DeptController{
		deptService: service.NewDeptService(),
	}
}

// GetTree 获取部门树形列表
// @Summary 获取部门树形列表
// @Description 获取所有部门的树形结构
// @Tags 部门管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.DeptResponse}
// @Router /api/v1/system/dept/tree [get]
func (ctrl *DeptController) GetTree(c *gin.Context) {
	tree, err := ctrl.deptService.GetTree()
	if err != nil {
		utils.Fail(c, 500, "查询部门树失败："+err.Error())
		return
	}

	utils.Success(c, tree)
}

// GetByID 根据ID获取部门详情
// @Summary 获取部门详情
// @Description 根据部门ID获取部门详细信息
// @Tags 部门管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response{data=response.DeptResponse}
// @Router /api/v1/system/dept/{id} [get]
func (ctrl *DeptController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的部门ID")
		return
	}

	dept, err := ctrl.deptService.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.Success(c, dept)
}

// Create 创建部门
// @Summary 创建部门
// @Description 创建新部门
// @Tags 部门管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.DeptCreateRequest true "部门信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/dept [post]
func (ctrl *DeptController) Create(c *gin.Context) {
	var req request.DeptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.deptService.Create(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 部门管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.DeptUpdateRequest true "部门信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/dept [put]
func (ctrl *DeptController) Update(c *gin.Context) {
	var req request.DeptUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.deptService.Update(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除部门
// @Summary 删除部门
// @Description 根据部门ID删除部门
// @Tags 部门管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response
// @Router /api/v1/system/dept/{id} [delete]
func (ctrl *DeptController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的部门ID")
		return
	}

	if err := ctrl.deptService.Delete(id); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
