package handler

import (
	"net/http"
	"strconv"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
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
// @Param request body model.CreateWorkflowRequest true "创建流程请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/workflow/workflows [post]
func (h *Handler) CreateWorkflow(c *gin.Context) {
	var req model.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户信息（需要中间件设置）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    401,
			Message: "用户未认证",
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "用户信息格式错误",
		})
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(&req, userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "创建工作流失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "工作流创建成功",
		Data:    workflow,
	})
}

func (h *Handler) StartWorkflow(c *gin.Context) {
	var req model.StartWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户信息（需要中间件设置）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    401,
			Message: "用户未认证",
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "用户信息格式错误",
		})
		return
	}

	workflow, err := h.workflowService.StartWorkflow(&req, userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "启动工作流失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "工作流启动成功",
		Data:    workflow,
	})
}

func (h *Handler) GetTasks(c *gin.Context) {

	// 从上下文获取用户信息（需要中间件设置）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    401,
			Message: "用户未认证",
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "用户信息格式错误",
		})
		return
	}

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

	tasks, err := h.workflowService.GetTasks(userModel, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "获取任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "获取任务成功",
		Data:    tasks,
	})
}

func (h *Handler) ApproveTask(c *gin.Context) {
	var req model.ApproveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 从上下文获取用户信息（需要中间件设置）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Code:    401,
			Message: "用户未认证",
		})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "用户信息格式错误",
		})
		return
	}

	task, err := h.workflowService.ApproveTask(&req, userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    500,
			Message: "审批任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "审批任务成功",
		Data:    task,
	})
}
