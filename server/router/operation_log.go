package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitOperationLogRouter 初始化操作日志路由
func InitOperationLogRouter(r *gin.RouterGroup) {
	logCtrl := controller.NewOperationLogController()

	logGroup := r.Group("/monitor/operation-log")
	logGroup.Use(middleware.Auth())
	{
		logGroup.GET("/list", middleware.Permission("monitor:operationLog:list"), logCtrl.GetList)
		logGroup.DELETE("/clear", middleware.Permission("monitor:operationLog:remove"), logCtrl.Clear)
	}
}
