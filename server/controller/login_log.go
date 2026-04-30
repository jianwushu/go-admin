package controller

import (
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// LoginLogController 登录日志控制器
type LoginLogController struct {
	service *service.LoginLogService
}

// NewLoginLogController 创建登录日志控制器实例
func NewLoginLogController() *LoginLogController {
	return &LoginLogController{
		service: service.NewLoginLogService(),
	}
}

// GetList 获取登录日志列表
// @Summary 获取登录日志列表
// @Description 分页查询登录日志，支持按用户名、IP、状态筛选
// @Tags 登录日志
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param username query string false "用户名"
// @Param ip query string false "IP地址"
// @Param status query int false "状态"
// @Success 200 {object} response.PageResponse{data=[]response.LoginLogResponse}
// @Router /api/v1/monitor/login-log/list [get]
func (ctrl *LoginLogController) GetList(c *gin.Context) {
	var req request.LoginLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	logs, total, err := ctrl.service.GetList(req)
	if err != nil {
		utils.Fail(c, 500, "查询登录日志失败："+err.Error())
		return
	}

	utils.PageSuccess(c, logs, total, req.GetPage(), req.GetPageSize())
}

// Clear 清空登录日志
// @Summary 清空登录日志
// @Description 清空所有登录日志记录
// @Tags 登录日志
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/monitor/login-log/clear [delete]
func (ctrl *LoginLogController) Clear(c *gin.Context) {
	if err := ctrl.service.Clear(); err != nil {
		utils.Fail(c, 500, "清空登录日志失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "清空成功", nil)
}
