package middleware

import (
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// DataScopeKey 数据权限信息在 Context 中的 Key
const DataScopeKey = "dataScope"

// DataScope 数据权限中间件
// 在 Auth 中间件之后使用，将数据权限信息写入 Context
// 后续 Repository 层可通过 GetDataScopeFromContext 获取权限信息进行数据过滤
func DataScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID（由 Auth 中间件设置）
		userID := GetCurrentUserID(c)
		if userID == 0 {
			utils.Fail(c, 401, "用户未登录")
			c.Abort()
			return
		}

		// 获取用户数据权限信息
		scopeInfo := utils.GetDataScopeInfo(userID)

		// 将数据权限信息写入 Context
		c.Set(DataScopeKey, scopeInfo)

		c.Next()
	}
}

// GetDataScopeFromContext 从 gin.Context 中获取数据权限信息
func GetDataScopeFromContext(c *gin.Context) *utils.DataScopeInfo {
	if scope, exists := c.Get(DataScopeKey); exists {
		if scopeInfo, ok := scope.(*utils.DataScopeInfo); ok {
			return scopeInfo
		}
	}
	return nil
}
