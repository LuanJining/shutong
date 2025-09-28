package model

// 请求/响应结构体

// CreateWorkflowRequest 创建流程请求
type CreateWorkflowRequest struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	SpaceID     uint                `json:"space_id" binding:"required"`
	Priority    int                 `json:"priority"`
	Steps       []CreateStepRequest `json:"steps" binding:"required"`
}

// CreateStepRequest 创建步骤请求
type CreateStepRequest struct {
	StepName         string `json:"step_name" binding:"required"`
	StepOrder        int    `json:"step_order" binding:"required"`
	ApproverType     string `json:"approver_type" binding:"required"`
	ApproverID       uint   `json:"approver_id"`
	IsRequired       bool   `json:"is_required"`
	TimeoutHours     int    `json:"timeout_hours"`
	ApprovalStrategy string `json:"approval_strategy"` // 审批策略: any, all, majority
}

// StartWorkflowRequest 启动流程请求
type StartWorkflowRequest struct {
	WorkflowID   uint   `json:"workflow_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ResourceType string `json:"resource_type" binding:"required"`
	ResourceID   uint   `json:"resource_id" binding:"required"`
	SpaceID      uint   `json:"space_id" binding:"required"`
	Priority     string `json:"priority"`
}

// ApproveTaskRequest 审批任务请求
type ApproveTaskRequest struct {
	Comment string `json:"comment"`
}

// RejectTaskRequest 拒绝任务请求
type RejectTaskRequest struct {
	Comment string `json:"comment" binding:"required"`
}

// TransferTaskRequest 转交任务请求
type TransferTaskRequest struct {
	ToUserID uint   `json:"to_user_id" binding:"required"`
	Comment  string `json:"comment"`
}

// API响应结构体
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PaginationResponse struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}
