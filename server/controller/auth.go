package controller

import (
	"strings"

	"go-admin/middleware"
	"go-admin/model/entity"
	"go-admin/model/request"
	"go-admin/model/response"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

// AuthController 认证控制器
type AuthController struct {
	authService     *service.AuthService
	loginLogService *service.LoginLogService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		authService:     service.NewAuthService(),
		loginLogService: service.NewLoginLogService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户名密码登录，返回 accessToken 和 refreshToken
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param data body request.LoginRequest true "登录参数"
// @Success 200 {object} response.Response{data=response.LoginResponse}
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	browserName, osName := parseUserAgent(ua)

	result, err := ctrl.authService.Login(req)
	if err != nil {
		// 记录登录失败日志
		ctrl.loginLogService.Create(&entity.LoginLog{
			Username: req.Username,
			IP:       ip,
			Browser:  browserName,
			OS:       osName,
			Status:   0,
			Msg:      err.Error(),
		})
		utils.Fail(c, 401, err.Error())
		return
	}

	// 记录登录成功日志
	ctrl.loginLogService.Create(&entity.LoginLog{
		Username: req.Username,
		IP:       ip,
		Browser:  browserName,
		OS:       osName,
		Status:   1,
		Msg:      "登录成功",
	})

	utils.SuccessWithMessage(c, "登录成功", result)
}

// parseUserAgent 解析 User-Agent 获取浏览器和操作系统信息
func parseUserAgent(ua string) (browser, os string) {
	agent := useragent.New(ua)
	browserName, _ := agent.Browser()
	if browserName == "" {
		browserName = "Unknown"
	}
	osName := agent.OS()
	if osName == "" {
		osName = "Unknown"
	}
	// 清理操作系统名称中的版本号
	if idx := strings.Index(osName, " "); idx > 0 {
		osName = osName[:idx]
	}
	return browserName, osName
}

// Logout 用户登出
// @Summary 用户登出
// @Description 将当前 Token 加入黑名单
// @Tags 认证管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 从 Header 获取当前 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Fail(c, 400, "未提供Token")
		return
	}

	token := authHeader[7:] // 去掉 "Bearer " 前缀

	if err := ctrl.authService.Logout(token); err != nil {
		utils.Fail(c, 500, "登出失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 使用 refreshToken 获取新的 accessToken
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param data body request.RefreshTokenRequest true "刷新Token参数"
// @Success 200 {object} response.Response{data=response.LoginResponse}
// @Router /api/v1/auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	result, err := ctrl.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		utils.Fail(c, 401, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "刷新成功", result)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息（含角色和权限）
// @Tags 认证管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=response.UserInfoResponse}
// @Router /api/v1/auth/userinfo [get]
func (ctrl *AuthController) GetUserInfo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.Fail(c, 401, "未登录")
		return
	}

	result, err := ctrl.authService.GetUserInfo(userID)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.Success(c, result)
}

// GetUserMenus 获取当前用户的菜单列表
// @Summary 获取用户菜单
// @Description 获取当前登录用户可访问的菜单列表（树形结构）
// @Tags 认证管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]response.MenuResponse}
// @Router /api/v1/user/menus [get]
func (ctrl *AuthController) GetUserMenus(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.Fail(c, 401, "未登录")
		return
	}

	menus, err := ctrl.authService.GetUserMenus(userID)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	// 将扁平菜单转换为树形结构
	tree := buildMenuTree(menus, 0)
	utils.Success(c, tree)
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []entity.Menu, parentID int64) []response.MenuResponse {
	var tree []response.MenuResponse
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := response.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Type:      menu.Type,
				Sort:      menu.Sort,
				Visible:   menu.Visible,
				Perms:     menu.Perms,
			}
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}
