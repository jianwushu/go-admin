package controller

import (
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// OperationLogController 操作日志控制器
type OperationLogController struct {
	service *service.OperationLogService
}

// NewOperationLogController 创建操作日志控制器实例
func NewOperationLogController() *OperationLogController {
	return &OperationLogController{
		service: service.NewOperationLogService(),
	}
}

// GetList 获取操作日志列表
// @Summary 获取操作日志列表
// @Description 分页查询操作日志，支持按模块、操作人、状态、方法筛选
// @Tags 操作日志
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param module query string false "模块名"
// @Param operator query string false "操作人"
// @Param status query int false "状态"
// @Param method query string false "请求方法"
// @Success 200 {object} response.PageResponse{data=[]response.OperationLogResponse}
// @Router /api/v1/monitor/operation-log/list [get]
func (ctrl *OperationLogController) GetList(c *gin.Context) {
	var req request.OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	logs, total, err := ctrl.service.GetList(req)
	if err != nil {
		utils.Fail(c, 500, "查询操作日志失败："+err.Error())
		return
	}

	utils.PageSuccess(c, logs, total, req.GetPage(), req.GetPageSize())
}

// Clear 清空操作日志
// @Summary 清空操作日志
// @Description 清空所有操作日志记录
// @Tags 操作日志
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/monitor/operation-log/clear [delete]
func (ctrl *OperationLogController) Clear(c *gin.Context) {
	if err := ctrl.service.Clear(); err != nil {
		utils.Fail(c, 500, "清空操作日志失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "清空成功", nil)
}
