package handler

import (
	"net/http"
	"strconv"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/workflow/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 处理器结构体
type Handler struct {
	db              *gorm.DB
	workflowService *service.WorkflowService
}

// NewHandler 创建新的处理器
func NewHandler(db *gorm.DB, workflowService *service.WorkflowService) *Handler {
	return &Handler{
		db:              db,
		workflowService: workflowService,
	}
}

// CreateWorkflow 创建审批流程
// @Summary 创建审批流程
// @Description 创建新的审批流程定义
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param request body models.CreateWorkflowRequest true "创建流程请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows [post]
func (h *Handler) CreateWorkflow(c *gin.Context) {
	var req models.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID（需要中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to create workflow: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Workflow created successfully",
		Data:    workflow,
	})
}

// GetWorkflows 获取流程列表
// @Summary 获取流程列表
// @Description 获取审批流程列表，支持分页
// @Tags workflow
// @Produce json
// @Param space_id query int false "空间ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows [get]
func (h *Handler) GetWorkflows(c *gin.Context) {
	// 获取查询参数
	spaceIDStr := c.Query("space_id")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// 解析参数
	spaceID := uint(0)
	if spaceIDStr != "" {
		if id, err := strconv.ParseUint(spaceIDStr, 10, 32); err == nil {
			spaceID = uint(id)
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.workflowService.GetWorkflows(spaceID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to get workflows: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

// GetWorkflow 获取流程详情
// @Summary 获取流程详情
// @Description 根据ID获取审批流程详情
// @Tags workflow
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "流程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows/{id} [get]
func (h *Handler) GetWorkflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid workflow ID",
		})
		return
	}

	workflow, err := h.workflowService.GetWorkflow(uint(id))
	if err != nil {
		if err.Error() == "workflow not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Workflow not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to get workflow: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Success",
		Data:    workflow,
	})
}

// UpdateWorkflow 更新审批流程
// @Summary 更新审批流程
// @Description 更新审批流程定义
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "流程ID"
// @Param request body models.CreateWorkflowRequest true "更新流程请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows/{id} [put]
func (h *Handler) UpdateWorkflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid workflow ID",
		})
		return
	}

	var req models.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	workflow, err := h.workflowService.UpdateWorkflow(uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == "workflow not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Workflow not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to update workflow: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Workflow updated successfully",
		Data:    workflow,
	})
}

// DeleteWorkflow 删除审批流程
// @Summary 删除审批流程
// @Description 删除审批流程定义
// @Tags workflow
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "流程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows/{id} [delete]
func (h *Handler) DeleteWorkflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid workflow ID",
		})
		return
	}

	err = h.workflowService.DeleteWorkflow(uint(id))
	if err != nil {
		if err.Error() == "workflow not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Workflow not found",
			})
		} else if err.Error() == "cannot delete workflow with active instances" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "Cannot delete workflow with active instances",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to delete workflow: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Workflow deleted successfully",
	})
}

// StartWorkflow 启动流程实例
// @Summary 启动流程实例
// @Description 启动审批流程实例
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param request body models.StartWorkflowRequest true "启动流程请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/instances [post]
func (h *Handler) StartWorkflow(c *gin.Context) {
	var req models.StartWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	instance, err := h.workflowService.StartWorkflow(&req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to start workflow: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Workflow started successfully",
		Data:    instance,
	})
}

// GetInstances 获取流程实例列表
// @Summary 获取流程实例列表
// @Description 获取流程实例列表，支持分页和筛选
// @Tags workflow
// @Produce json
// @Param workflow_id query int false "流程ID"
// @Param space_id query int false "空间ID"
// @Param status query string false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/instances [get]
func (h *Handler) GetInstances(c *gin.Context) {
	// 获取查询参数
	workflowIDStr := c.Query("workflow_id")
	spaceIDStr := c.Query("space_id")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// 解析参数
	workflowID := uint(0)
	if workflowIDStr != "" {
		if id, err := strconv.ParseUint(workflowIDStr, 10, 32); err == nil {
			workflowID = uint(id)
		}
	}

	spaceID := uint(0)
	if spaceIDStr != "" {
		if id, err := strconv.ParseUint(spaceIDStr, 10, 32); err == nil {
			spaceID = uint(id)
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.workflowService.GetInstances(workflowID, spaceID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to get instances: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

// GetInstance 获取流程实例详情
// @Summary 获取流程实例详情
// @Description 根据ID获取流程实例详情
// @Tags workflow
// @Produce json
// @Param id path int true "实例ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/instances/{id} [get]
func (h *Handler) GetInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid instance ID",
		})
		return
	}

	instance, err := h.workflowService.GetInstance(uint(id))
	if err != nil {
		if err.Error() == "instance not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Instance not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to get instance: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Success",
		Data:    instance,
	})
}

// CancelInstance 取消流程实例
// @Summary 取消流程实例
// @Description 取消流程实例
// @Tags workflow
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "实例ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/instances/{id}/cancel [put]
func (h *Handler) CancelInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid instance ID",
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	err = h.workflowService.CancelInstance(uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "instance not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Instance not found",
			})
		} else if err.Error() == "unauthorized to cancel this instance" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Code:    403,
				Message: "Unauthorized to cancel this instance",
			})
		} else if err.Error() == "instance cannot be cancelled in current status" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "Instance cannot be cancelled in current status",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to cancel instance: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Instance cancelled successfully",
	})
}

// GetMyTasks 获取我的待办任务
// @Summary 获取我的待办任务
// @Description 获取当前用户的待办审批任务
// @Tags workflow
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/tasks [get]
func (h *Handler) GetMyTasks(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.workflowService.GetMyTasks(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to get tasks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	})
}

// ApproveTask 审批通过
// @Summary 审批通过
// @Description 审批通过任务
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "任务ID"
// @Param request body models.ApproveTaskRequest true "审批请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/tasks/{id}/approve [post]
func (h *Handler) ApproveTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid task ID",
		})
		return
	}

	var req models.ApproveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	err = h.workflowService.ApproveTask(uint(id), userID.(uint), req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to approve task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Task approved successfully",
	})
}

// RejectTask 审批拒绝
// @Summary 审批拒绝
// @Description 审批拒绝任务
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "任务ID"
// @Param request body models.RejectTaskRequest true "拒绝请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/tasks/{id}/reject [post]
func (h *Handler) RejectTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid task ID",
		})
		return
	}

	var req models.RejectTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	err = h.workflowService.RejectTask(uint(id), userID.(uint), req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to reject task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Task rejected successfully",
	})
}

// TransferTask 转交任务
// @Summary 转交任务
// @Description 转交审批任务给其他用户
// @Tags workflow
// @Accept json
// @Produce json
// @Param X-User-ID header string true "用户ID"
// @Param id path int true "任务ID"
// @Param request body models.TransferTaskRequest true "转交请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/tasks/{id}/transfer [post]
func (h *Handler) TransferTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid task ID",
		})
		return
	}

	var req models.TransferTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Code:    401,
			Message: "User not authenticated",
		})
		return
	}

	err = h.workflowService.TransferTask(uint(id), userID.(uint), req.ToUserID, req.Comment)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Task not found",
			})
		} else if err.Error() == "unauthorized to transfer this task" {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Code:    403,
				Message: "Unauthorized to transfer this task",
			})
		} else if err.Error() == "task is not pending" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "Task is not pending",
			})
		} else if err.Error() == "target user not found or inactive" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Code:    400,
				Message: "Target user not found or inactive",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    500,
				Message: "Failed to transfer task: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    200,
		Message: "Task transferred successfully",
	})
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
