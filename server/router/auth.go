package router

import (
	"go-admin/controller"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

// InitAuthRouter 初始化认证相关路由
func InitAuthRouter(r *gin.RouterGroup) {
	authCtrl := controller.NewAuthController()

	authGroup := r.Group("/auth")
	{
		// 公开接口（无需认证）
		authGroup.POST("/login", authCtrl.Login)
		authGroup.POST("/refresh", authCtrl.RefreshToken)

		// 需要认证的接口
		authGroup.Use(middleware.Auth())
		{
			authGroup.POST("/logout", authCtrl.Logout)
		}
	}

	// 用户相关路由（需要认证）
	userCtrl := controller.NewUserController()
	userGroup := r.Group("/user")
	userGroup.Use(middleware.Auth())
	{
		userGroup.GET("/info", authCtrl.GetUserInfo)
		userGroup.GET("/menus", authCtrl.GetUserMenus)
		userGroup.GET("/profile", userCtrl.GetProfile)
		userGroup.PUT("/profile", userCtrl.UpdateProfile)
		userGroup.PUT("/change-password", userCtrl.ChangePassword)
	}
}
