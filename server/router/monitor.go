package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitMonitorRouter 初始化服务器监控相关路由
func InitMonitorRouter(r *gin.RouterGroup) {
	monitorCtrl := controller.NewMonitorController()

	monitorGroup := r.Group("/monitor")
	monitorGroup.Use(middleware.Auth())
	{
		monitorGroup.GET("/server", monitorCtrl.GetServerInfo)
	}
}
