package utils

import (
	"net/http"

	"go-admin/model/response"

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response.Response{
		Code: 0,
		Msg:  "操作成功",
		Data: data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, response.Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, response.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, response.PageResponse{
		Code:  0,
		Msg:   "查询成功",
		Data:  data,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, response.Response{
		Code: 5003,
		Msg:  "系统内部错误",
		Data: nil,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Code: 1001,
		Msg:  "Token已过期或无效",
		Data: nil,
	})
}

// Forbidden 无权限
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.Response{
		Code: 1003,
		Msg:  "无操作权限",
		Data: nil,
	})
}
