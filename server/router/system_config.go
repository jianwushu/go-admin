package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitSystemConfigRouter 初始化系统配置路由
func InitSystemConfigRouter(r *gin.RouterGroup) {
	configCtrl := controller.NewSystemConfigController()

	configGroup := r.Group("/system/config")
	configGroup.Use(middleware.Auth())
	{
		configGroup.GET("/all", middleware.Permission("system:config:list"), configCtrl.GetAll)
		configGroup.GET("/list", middleware.Permission("system:config:list"), configCtrl.GetList)
		configGroup.GET("/key/:key", middleware.Permission("system:config:list"), configCtrl.GetByKey)
		configGroup.GET("/keys", middleware.Permission("system:config:list"), configCtrl.GetByKeys)
		configGroup.PUT("", middleware.Permission("system:config:edit"), configCtrl.Update)
		configGroup.PUT("/batch", middleware.Permission("system:config:edit"), configCtrl.BatchUpdate)
		configGroup.POST("/upload", middleware.Permission("system:config:edit"), configCtrl.UploadLogo)
	}
}
