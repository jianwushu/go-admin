package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"go-admin/global"
	"go-admin/model/entity"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLog 操作日志中间件
// 自动记录请求的模块、操作、方法、URL、IP、请求参数、响应数据、状态、耗时等信息
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的请求（如 GET 查询类接口可选择性跳过）
		method := c.Request.Method
		// 只记录写操作（POST、PUT、DELETE）
		if method == "GET" {
			c.Next()
			return
		}

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 恢复请求体，供后续 Handler 使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 截断过长的请求参数
		if len(requestBody) > 2000 {
			requestBody = requestBody[:2000] + "...(truncated)"
		}

		// 记录开始时间
		startTime := time.Now()

		// 使用自定义 ResponseWriter 捕获响应数据
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 执行后续 Handler
		c.Next()

		// 计算耗时（毫秒）
		duration := time.Since(startTime).Milliseconds()

		// 获取响应数据
		responseBody := blw.body.String()
		if len(responseBody) > 2000 {
			responseBody = responseBody[:2000] + "...(truncated)"
		}

		// 解析模块和操作（从 URL 路径推断）
		module, action := parseModuleAndAction(c.FullPath(), method)

		// 获取操作人
		operator := GetCurrentUsername(c)
		if operator == "" {
			operator = "unknown"
		}

		// 判断请求状态
		status := 1 // 成功
		var errorMsg string
		if len(c.Errors) > 0 {
			status = 0
			errorMsg = c.Errors.String()
		}

		// 异步写入日志（不阻塞响应）
		log := entity.OperationLog{
			Module:       module,
			Action:       action,
			Method:       method,
			URL:          c.Request.URL.Path,
			IP:           c.ClientIP(),
			Operator:     operator,
			RequestParam: requestBody,
			ResponseData: responseBody,
			Status:       status,
			ErrorMsg:     errorMsg,
			Duration:     duration,
		}

		go func() {
			if err := global.DB.Create(&log).Error; err != nil {
				global.Log.Error("写入操作日志失败", zap.Error(err))
			}
		}()
	}
}

// bodyLogWriter 自定义 ResponseWriter，用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// parseModuleAndAction 根据 URL 路径和请求方法解析模块和操作名称
func parseModuleAndAction(path string, method string) (string, string) {
	// 默认值
	module := "系统"
	action := "操作"

	// 根据路径解析模块
	switch {
	case strings.Contains(path, "/system/user"):
		module = "用户管理"
	case strings.Contains(path, "/system/role"):
		module = "角色管理"
	case strings.Contains(path, "/system/menu"):
		module = "菜单管理"
	case strings.Contains(path, "/system/dept"):
		module = "部门管理"
	case strings.Contains(path, "/auth"):
		module = "认证"
	case strings.Contains(path, "/monitor"):
		module = "系统监控"
	case strings.Contains(path, "/codegen"):
		module = "代码生成"
	}

	// 根据路径和方法解析操作
	switch {
	case strings.HasSuffix(path, "/list") || strings.HasSuffix(path, "/tree") || strings.HasSuffix(path, "/all"):
		action = "查询"
	case strings.HasSuffix(path, "/change-status"):
		action = "修改状态"
	case strings.HasSuffix(path, "/reset-password"):
		action = "重置密码"
	default:
		switch method {
		case "POST":
			action = "新增"
		case "PUT":
			action = "修改"
		case "DELETE":
			action = "删除"
		case "GET":
			action = "查询"
		}
	}

	return module, action
}
