package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitLoginLogRouter 初始化登录日志路由
func InitLoginLogRouter(r *gin.RouterGroup) {
	logCtrl := controller.NewLoginLogController()

	logGroup := r.Group("/monitor/login-log")
	logGroup.Use(middleware.Auth())
	{
		logGroup.GET("/list", middleware.Permission("monitor:loginLog:list"), logCtrl.GetList)
		logGroup.DELETE("/clear", middleware.Permission("monitor:loginLog:remove"), logCtrl.Clear)
	}
}
