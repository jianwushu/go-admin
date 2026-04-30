package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// Init{{.ClassName}}Router 初始化{{.FunctionName}}路由
func Init{{.ClassName}}Router(r *gin.RouterGroup) {
	ctrl := controller.New{{.ClassName}}Controller()

	group := r.Group("/{{.ModuleName}}/{{.BusinessName}}")
	group.Use(middleware.Auth(), middleware.DataScope(), middleware.OperationLog())
	{
		group.GET("/list", middleware.Permission("{{.ModuleName}}:{{.BusinessName}}:list"), ctrl.GetList)
		group.GET("/:id", middleware.Permission("{{.ModuleName}}:{{.BusinessName}}:query"), ctrl.GetByID)
		group.POST("", middleware.Permission("{{.ModuleName}}:{{.BusinessName}}:add"), ctrl.Create)
		group.PUT("", middleware.Permission("{{.ModuleName}}:{{.BusinessName}}:edit"), ctrl.Update)
		group.DELETE("/:id", middleware.Permission("{{.ModuleName}}:{{.BusinessName}}:remove"), ctrl.Delete)
	}
}
