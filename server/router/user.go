package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitUserRouter 初始化用户管理路由
func InitUserRouter(r *gin.RouterGroup) {
	userCtrl := controller.NewUserController()

	userGroup := r.Group("/system/user")
	userGroup.Use(middleware.Auth(), middleware.DataScope(), middleware.OperationLog())
	{
		userGroup.GET("/list", middleware.Permission("system:user:list"), userCtrl.GetList)
		userGroup.GET("/:id", middleware.Permission("system:user:query"), userCtrl.GetByID)
		userGroup.POST("", middleware.Permission("system:user:add"), userCtrl.Create)
		userGroup.PUT("", middleware.Permission("system:user:edit"), userCtrl.Update)
		userGroup.DELETE("/:id", middleware.Permission("system:user:remove"), userCtrl.Delete)
		userGroup.PUT("/reset-password", middleware.Permission("system:user:resetPwd"), userCtrl.ResetPassword)
		userGroup.PUT("/change-status", middleware.Permission("system:user:edit"), userCtrl.ChangeStatus)
	}
}
