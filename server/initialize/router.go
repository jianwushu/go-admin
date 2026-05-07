package initialize

import (
	"go-admin/global"
	"go-admin/middleware"
	"go-admin/router"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	if global.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 注册全局中间件
	r.Use(middleware.Recovery()) // 自定义异常恢复中间件（替代 gin.Recovery）
	r.Use(middleware.Cors())     // CORS 跨域中间件

	// 注册 API 路由
	apiGroup := r.Group("/api/v1")
	{
		router.InitAuthRouter(apiGroup)
		router.InitUserRouter(apiGroup)
		router.InitRoleRouter(apiGroup)
		router.InitMenuRouter(apiGroup)
		router.InitDeptRouter(apiGroup)
		router.InitDashboardRouter(apiGroup)
		router.InitMonitorRouter(apiGroup)
		router.InitOperationLogRouter(apiGroup)
		router.InitLoginLogRouter(apiGroup)
		router.InitCodegenRouter(apiGroup)
		router.InitSystemConfigRouter(apiGroup)
		router.InitJobRouter(apiGroup)
	}

	return r
}
