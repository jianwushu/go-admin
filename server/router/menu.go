package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitMenuRouter 初始化菜单管理路由
func InitMenuRouter(r *gin.RouterGroup) {
	menuCtrl := controller.NewMenuController()

	menuGroup := r.Group("/system/menu")
	menuGroup.Use(middleware.Auth(), middleware.OperationLog())
	{
		menuGroup.GET("/tree", middleware.Permission("system:menu:list"), menuCtrl.GetTree)
		menuGroup.GET("/:id", middleware.Permission("system:menu:query"), menuCtrl.GetByID)
		menuGroup.POST("", middleware.Permission("system:menu:add"), menuCtrl.Create)
		menuGroup.PUT("", middleware.Permission("system:menu:edit"), menuCtrl.Update)
		menuGroup.DELETE("/:id", middleware.Permission("system:menu:remove"), menuCtrl.Delete)
	}
}
