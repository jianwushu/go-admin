package controller

import (
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	dashboardService *service.DashboardService
}

// NewDashboardController 创建仪表盘控制器实例
func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: service.NewDashboardService(),
	}
}

// GetStats 获取仪表盘统计数据
// @Summary 获取仪表盘统计数据
// @Description 获取用户数、角色数、菜单数、部门数、今日登录次数
// @Tags 仪表盘
// @Produce json
// @Success 200 {object} response.Response{data=response.DashboardResponse}
// @Router /api/v1/dashboard/stats [get]
func (ctrl *DashboardController) GetStats(c *gin.Context) {
	result, err := ctrl.dashboardService.GetDashboardStats()
	if err != nil {
		utils.Fail(c, 500, "获取统计数据失败："+err.Error())
		return
	}

	utils.Success(c, result)
}
