package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"go-admin/global"
	"go-admin/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 全局异常恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是连接断开错误，这类错误不需要记录堆栈
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 检查 URL 解码错误
				if ue, ok := err.(*url.Error); ok {
					if strings.Contains(ue.Error(), "invalid URL escape") {
						brokenPipe = true
					}
				}

				// 获取请求信息
				httpRequest, _ := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path), c.Request.URL.RawQuery

				if brokenPipe {
					global.Log.Warn("连接断开",
						zap.String("request", httpRequest),
						zap.Any("error", err),
					)
					c.Abort()
					return
				}

				// 记录完整的 panic 堆栈信息
				global.Log.Error("请求处理异常",
					zap.String("request", httpRequest),
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code: 5000,
					Msg:  "系统内部错误",
					Data: nil,
				})
			}
		}()

		c.Next()
	}
}
