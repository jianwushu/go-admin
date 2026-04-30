package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitDashboardRouter 初始化仪表盘相关路由
func InitDashboardRouter(r *gin.RouterGroup) {
	dashboardCtrl := controller.NewDashboardController()

	dashboardGroup := r.Group("/dashboard")
	dashboardGroup.Use(middleware.Auth())
	{
		dashboardGroup.GET("/stats", dashboardCtrl.GetStats)
	}
}
