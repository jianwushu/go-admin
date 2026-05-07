package controller

import (
	"strconv"
	"strings"

	"go-admin/middleware"
	"go-admin/model/request"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// JobController 定时任务控制器
type JobController struct {
	service *service.JobService
}

// NewJobController 创建定时任务控制器实例
func NewJobController() *JobController {
	return &JobController{
		service: service.NewJobService(),
	}
}

// GetList 获取任务列表（分页）
// @Summary 获取定时任务列表
// @Description 分页查询定时任务
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param name query string false "任务名称"
// @Param jobType query int false "任务类型"
// @Param status query int false "状态"
// @Success 200 {object} response.PageResponse{data=[]response.JobResponse}
// @Router /api/v1/tool/job/list [get]
func (ctrl *JobController) GetList(c *gin.Context) {
	var req request.JobListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	jobs, total, err := ctrl.service.GetList(req)
	if err != nil {
		utils.Fail(c, 500, "查询任务失败："+err.Error())
		return
	}

	utils.PageSuccess(c, jobs, total, req.GetPage(), req.GetPageSize())
}

// GetByID 根据ID获取任务详情
// @Summary 获取任务详情
// @Description 根据ID获取定时任务详情
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response{data=response.JobResponse}
// @Router /api/v1/tool/job/{id} [get]
func (ctrl *JobController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的任务ID")
		return
	}

	job, err := ctrl.service.GetByID(id)
	if err != nil {
		utils.Fail(c, 500, "查询任务失败："+err.Error())
		return
	}

	utils.Success(c, job)
}

// Create 创建任务
// @Summary 创建定时任务
// @Description 创建新的定时任务
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param data body request.JobCreateRequest true "任务创建请求"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job [post]
func (ctrl *JobController) Create(c *gin.Context) {
	var req request.JobCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.service.Create(req, userID); err != nil {
		utils.Fail(c, 500, "创建任务失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新任务
// @Summary 更新定时任务
// @Description 更新定时任务信息
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param data body request.JobUpdateRequest true "任务更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job [put]
func (ctrl *JobController) Update(c *gin.Context) {
	var req request.JobUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.service.Update(req, userID); err != nil {
		utils.Fail(c, 500, "更新任务失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// ChangeStatus 修改任务状态
// @Summary 修改任务状态
// @Description 启用/禁用定时任务
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param data body request.JobChangeStatusRequest true "状态修改请求"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job/change-status [put]
func (ctrl *JobController) ChangeStatus(c *gin.Context) {
	var req request.JobChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	userID := middleware.GetCurrentUserID(c)

	if err := ctrl.service.ChangeStatus(req, userID); err != nil {
		utils.Fail(c, 500, "修改状态失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "操作成功", nil)
}

// Delete 删除任务
// @Summary 删除定时任务
// @Description 删除定时任务（支持批量）
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param id path string true "任务ID，多个用逗号分隔"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job/{id} [delete]
func (ctrl *JobController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.Fail(c, 400, "任务ID不能为空")
		return
	}

	// 支持批量删除
	idStrs := strings.Split(idStr, ",")
	var ids []int64
	for _, s := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			utils.Fail(c, 400, "参数错误：无效的任务ID")
			return
		}
		ids = append(ids, id)
	}

	if len(ids) == 1 {
		if err := ctrl.service.Delete(ids[0]); err != nil {
			utils.Fail(c, 500, "删除任务失败："+err.Error())
			return
		}
	} else {
		if err := ctrl.service.DeleteBatch(ids); err != nil {
			utils.Fail(c, 500, "批量删除任务失败："+err.Error())
			return
		}
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// RunOnce 手动执行一次任务
// @Summary 手动执行任务
// @Description 手动触发执行一次定时任务
// @Tags 定时任务
// @Security BearerAuth
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job/run-once/{id} [post]
func (ctrl *JobController) RunOnce(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的任务ID")
		return
	}

	if err := ctrl.service.RunOnce(id); err != nil {
		utils.Fail(c, 500, "执行任务失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "任务已触发执行", nil)
}

// ==================== 任务日志 ====================

// GetLogList 获取任务日志列表（分页）
// @Summary 获取任务日志列表
// @Description 分页查询任务执行日志
// @Tags 定时任务日志
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param jobId query int false "任务ID"
// @Param status query int false "状态"
// @Success 200 {object} response.PageResponse{data=[]response.JobLogResponse}
// @Router /api/v1/tool/job/log/list [get]
func (ctrl *JobController) GetLogList(c *gin.Context) {
	var req request.JobLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误："+err.Error())
		return
	}

	logs, total, err := ctrl.service.GetLogList(req)
	if err != nil {
		utils.Fail(c, 500, "查询日志失败："+err.Error())
		return
	}

	utils.PageSuccess(c, logs, total, req.GetPage(), req.GetPageSize())
}

// CleanLogs 清理指定任务的日志
// @Summary 清理任务日志
// @Description 清理指定任务的所有执行日志
// @Tags 定时任务日志
// @Security BearerAuth
// @Produce json
// @Param jobId query int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/tool/job/log/clean [delete]
func (ctrl *JobController) CleanLogs(c *gin.Context) {
	jobID, err := strconv.ParseInt(c.Query("jobId"), 10, 64)
	if err != nil {
		utils.Fail(c, 400, "参数错误：无效的任务ID")
		return
	}

	if err := ctrl.service.CleanLogs(jobID); err != nil {
		utils.Fail(c, 500, "清理日志失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "清理成功", nil)
}
