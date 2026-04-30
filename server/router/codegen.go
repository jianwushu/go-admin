package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitCodegenRouter 初始化代码生成路由
func InitCodegenRouter(r *gin.RouterGroup) {
	codegenCtrl := controller.NewCodegenController()

	codegenGroup := r.Group("/codegen")
	codegenGroup.Use(middleware.Auth(), middleware.OperationLog())
	{
		// 表信息
		codegenGroup.GET("/tables", codegenCtrl.GetTables)
		codegenGroup.GET("/columns/:tableName", codegenCtrl.GetColumns)

		// 代码预览与生成
		codegenGroup.POST("/preview", codegenCtrl.Preview)
		codegenGroup.POST("/generate", codegenCtrl.Generate)

		// 配置管理
		codegenGroup.GET("/configs", codegenCtrl.GetAllConfigs)
		codegenGroup.GET("/config/:tableName", codegenCtrl.GetConfig)
		codegenGroup.POST("/config", codegenCtrl.SaveConfig)
		codegenGroup.DELETE("/config/:id", codegenCtrl.DeleteConfig)
	}
}
