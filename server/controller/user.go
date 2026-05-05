package controller

import (
	"strconv"

	"go-admin/middleware"
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// GetList 获取用户列表
// @Summary 获取用户列表
// @Description 分页查询用户列表，支持按用户名、状态、部门筛选
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param username query string false "用户名"
// @Param status query int false "状态"
// @Param deptId query int false "部门ID"
// @Success 200 {object} response.PageResponse{data=[]response.UserResponse}
// @Router /api/v1/system/user/list [get]
func (ctrl *UserController) GetList(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	// 获取数据权限信息
	scopeInfo := middleware.GetDataScopeFromContext(c)

	users, total, err := ctrl.userService.GetList(req, scopeInfo)
	if err != nil {
		utils.Fail(c, 500, "查询用户列表失败："+err.Error())
		return
	}

	utils.PageSuccess(c, users, total, req.GetPage(), req.GetPageSize())
}

// GetByID 根据ID获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=response.UserResponse}
// @Router /api/v1/system/user/{id} [get]
func (ctrl *UserController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的用户ID")
		return
	}

	user, err := ctrl.userService.GetByID(id)
	if err != nil {
		utils.Fail(c, 404, err.Error())
		return
	}

	utils.Success(c, user)
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.UserCreateRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/user [post]
func (ctrl *UserController) Create(c *gin.Context) {
	var req request.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.Create(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.UserUpdateRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/system/user [put]
func (ctrl *UserController) Update(c *gin.Context) {
	var req request.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.Update(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 根据用户ID删除用户
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/system/user/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的用户ID")
		return
	}

	if err := ctrl.userService.Delete(id); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ResetPassword 重置用户密码
// @Summary 重置密码
// @Description 重置指定用户的密码
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.UserResetPasswordRequest true "重置密码参数"
// @Success 200 {object} response.Response
// @Router /api/v1/system/user/reset-password [put]
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	var req request.UserResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.ResetPassword(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码重置成功", nil)
}

// ChangeStatus 修改用户状态
// @Summary 修改用户状态
// @Description 启用或禁用用户
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.UserChangeStatusRequest true "状态参数"
// @Success 200 {object} response.Response
// @Router /api/v1/system/user/change-status [put]
func (ctrl *UserController) ChangeStatus(c *gin.Context) {
	var req request.UserChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.ChangeStatus(req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "状态修改成功", nil)
}

// GetProfile 获取当前登录用户的个人资料
// @Summary 获取个人资料
// @Description 获取当前登录用户的个人资料信息
// @Tags 个人中心
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=response.UserProfileResponse}
// @Router /api/v1/user/profile [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.Fail(c, 401, "未登录")
		return
	}

	profile, err := ctrl.userService.GetProfile(userID)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.Success(c, profile)
}

// UpdateProfile 更新当前登录用户的个人资料
// @Summary 更新个人资料
// @Description 更新当前登录用户的昵称、邮箱、手机号、头像等信息
// @Tags 个人中心
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.UserProfileUpdateRequest true "个人资料"
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [put]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.Fail(c, 401, "未登录")
		return
	}

	var req request.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.UpdateProfile(userID, req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改当前登录用户的密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码，需要验证旧密码
// @Tags 个人中心
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body request.ChangePasswordRequest true "密码参数"
// @Success 200 {object} response.Response
// @Router /api/v1/user/change-password [put]
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		utils.Fail(c, 401, "未登录")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID, req); err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}
