package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitJobRouter 初始化定时任务路由
func InitJobRouter(r *gin.RouterGroup) {
	jobCtrl := controller.NewJobController()

	jobGroup := r.Group("/tool/job")
	jobGroup.Use(middleware.Auth())
	{
		// 任务管理
		jobGroup.GET("/list", middleware.Permission("tool:job:list"), jobCtrl.GetList)
		jobGroup.GET("/:id", middleware.Permission("tool:job:list"), jobCtrl.GetByID)
		jobGroup.POST("", middleware.Permission("tool:job:add"), jobCtrl.Create)
		jobGroup.PUT("", middleware.Permission("tool:job:edit"), jobCtrl.Update)
		jobGroup.PUT("/change-status", middleware.Permission("tool:job:edit"), jobCtrl.ChangeStatus)
		jobGroup.DELETE("/:id", middleware.Permission("tool:job:delete"), jobCtrl.Delete)
		jobGroup.POST("/run-once/:id", middleware.Permission("tool:job:edit"), jobCtrl.RunOnce)

		// 任务日志
		jobGroup.GET("/log/list", middleware.Permission("tool:job:list"), jobCtrl.GetLogList)
		jobGroup.DELETE("/log/clean", middleware.Permission("tool:job:delete"), jobCtrl.CleanLogs)
	}
}
