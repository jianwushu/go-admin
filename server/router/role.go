package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoleRouter 初始化角色管理路由
func InitRoleRouter(r *gin.RouterGroup) {
	roleCtrl := controller.NewRoleController()

	roleGroup := r.Group("/system/role")
	roleGroup.Use(middleware.Auth(), middleware.OperationLog())
	{
		roleGroup.GET("/list", middleware.Permission("system:role:list"), roleCtrl.GetList)
		roleGroup.GET("/all", middleware.Permission("system:role:list"), roleCtrl.GetAll)
		roleGroup.GET("/:id", middleware.Permission("system:role:query"), roleCtrl.GetByID)
		roleGroup.POST("", middleware.Permission("system:role:add"), roleCtrl.Create)
		roleGroup.PUT("", middleware.Permission("system:role:edit"), roleCtrl.Update)
		roleGroup.DELETE("/:id", middleware.Permission("system:role:remove"), roleCtrl.Delete)
		roleGroup.PUT("/change-status", middleware.Permission("system:role:edit"), roleCtrl.ChangeStatus)
	}
}
