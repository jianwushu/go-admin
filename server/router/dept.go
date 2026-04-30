package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitDeptRouter 初始化部门管理路由
func InitDeptRouter(r *gin.RouterGroup) {
	deptCtrl := controller.NewDeptController()

	deptGroup := r.Group("/system/dept")
	deptGroup.Use(middleware.Auth(), middleware.OperationLog())
	{
		deptGroup.GET("/tree", middleware.Permission("system:dept:list"), deptCtrl.GetTree)
		deptGroup.GET("/:id", middleware.Permission("system:dept:query"), deptCtrl.GetByID)
		deptGroup.POST("", middleware.Permission("system:dept:add"), deptCtrl.Create)
		deptGroup.PUT("", middleware.Permission("system:dept:edit"), deptCtrl.Update)
		deptGroup.DELETE("/:id", middleware.Permission("system:dept:remove"), deptCtrl.Delete)
	}
}
