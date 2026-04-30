package middleware

import (
	"context"
	"fmt"
	"strings"

	"go-admin/global"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// Auth JWT 认证中间件（有状态设计）
// 解析 Authorization Header 中的 Bearer Token
// 先验证 JWT 签名，再检查 Redis 中是否存在对应 Token
// 校验通过后将用户信息写入 gin.Context
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Fail(c, 401, "未提供认证信息")
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Fail(c, 401, "认证格式错误，应为 Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 1. 解析并验证 JWT 签名
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Fail(c, 401, "Token 已失效或无效："+err.Error())
			c.Abort()
			return
		}

		// 2. 检查 Redis 中是否存在该 Token（有状态验证）
		ctx := context.Background()
		userTokenKey := fmt.Sprintf("jwt:user:%d", claims.UserID)
		uuid, err := global.Redis.Get(ctx, userTokenKey).Result()
		if err != nil {
			utils.Fail(c, 401, "Token 已失效，请重新登录")
			c.Abort()
			return
		}

		tokenKey := fmt.Sprintf("jwt:token:%d:%s", claims.UserID, uuid)
		storedToken, err := global.Redis.Get(ctx, tokenKey).Result()
		if err != nil || storedToken != tokenString {
			utils.Fail(c, 401, "Token 已失效，请重新登录")
			c.Abort()
			return
		}

		// 将用户信息写入 Context，供后续 Handler 使用
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetCurrentUserID 从 gin.Context 中获取当前登录用户ID
func GetCurrentUserID(c *gin.Context) int64 {
	if userID, exists := c.Get("userId"); exists {
		if id, ok := userID.(int64); ok {
			return id
		}
	}
	return 0
}

// GetCurrentUsername 从 gin.Context 中获取当前登录用户名
func GetCurrentUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}
