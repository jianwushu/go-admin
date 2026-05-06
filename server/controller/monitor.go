package controller

import (
	"go-admin/service"
	"go-admin/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MonitorController 服务器监控控制器
type MonitorController struct {
	monitorService *service.MonitorService
	onlineService  *service.OnlineService
}

// NewMonitorController 创建监控控制器实例
func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: service.NewMonitorService(),
		onlineService:  service.NewOnlineService(),
	}
}

// GetServerInfo 获取服务器监控信息
// @Summary 获取服务器监控信息
// @Description 获取 CPU、内存、磁盘、Go运行时、数据库、Redis 状态
// @Tags 服务器监控
// @Produce json
// @Success 200 {object} response.Response{data=response.ServerMonitorResponse}
// @Router /api/v1/monitor/server [get]
func (ctrl *MonitorController) GetServerInfo(c *gin.Context) {
	result, err := ctrl.monitorService.GetServerInfo()
	if err != nil {
		utils.Fail(c, 500, "获取服务器信息失败："+err.Error())
		return
	}

	utils.Success(c, result)
}

// GetOnlineUsers 获取在线用户列表
// @Summary 获取在线用户列表
// @Description 查询当前在线的用户列表，支持按用户名搜索
// @Tags 在线用户
// @Produce json
// @Param username query string false "用户名（模糊搜索）"
// @Success 200 {object} response.Response{data=[]response.OnlineUser}
// @Router /api/v1/monitor/online [get]
func (ctrl *MonitorController) GetOnlineUsers(c *gin.Context) {
	username := c.Query("username")

	result, err := ctrl.onlineService.GetOnlineUsers(username)
	if err != nil {
		utils.Fail(c, 500, "获取在线用户列表失败："+err.Error())
		return
	}

	utils.Success(c, result)
}

// ForceLogoutOnlineUser 强制用户下线
// @Summary 强制用户下线
// @Description 强制指定用户下线，删除其 Token
// @Tags 在线用户
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/monitor/online/{userId} [delete]
func (ctrl *MonitorController) ForceLogoutOnlineUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.Fail(c, 400, "无效的用户ID")
		return
	}

	if err := ctrl.onlineService.ForceLogout(userID); err != nil {
		utils.Fail(c, 500, "强制下线失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "强制下线成功", nil)
}
