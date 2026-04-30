package controller

import (
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// MonitorController 服务器监控控制器
type MonitorController struct {
	monitorService *service.MonitorService
}

// NewMonitorController 创建监控控制器实例
func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: service.NewMonitorService(),
	}
}

// GetServerInfo 获取服务器监控信息
// @Summary 获取服务器监控信息
// @Description 获取 CPU、内存、磁盘、Go运行时、数据库、Redis 状态
// @Tags 服务器监控
// @Produce json
// @Success 200 {object} response.Response{data=response.ServerMonitorResponse}
// @Router /api/v1/monitor/server [get]
func (ctrl *MonitorController) GetServerInfo(c *gin.Context) {
	result, err := ctrl.monitorService.GetServerInfo()
	if err != nil {
		utils.Fail(c, 500, "获取服务器信息失败："+err.Error())
		return
	}

	utils.Success(c, result)
}
