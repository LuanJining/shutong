package models

import (
	"time"
)

// WorkflowDefinition 审批流程定义
type WorkflowDefinition struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"` // 流程名称
	Description string         `json:"description" gorm:"size:500"`   // 流程描述
	SpaceID     uint           `json:"space_id" gorm:"not null"`      // 所属空间ID
	IsActive    bool           `json:"is_active" gorm:"default:true"` // 是否启用
	Priority    int            `json:"priority" gorm:"default:0"`     // 优先级
	CreatedBy   uint           `json:"created_by" gorm:"not null"`    // 创建人
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Steps       []WorkflowStep `json:"steps" gorm:"foreignKey:WorkflowID"` // 流程步骤
}

// WorkflowStep 审批流程步骤
type WorkflowStep struct {
	ID           uint               `json:"id" gorm:"primaryKey"`
	WorkflowID   uint               `json:"workflow_id" gorm:"not null"`           // 流程ID
	Workflow     WorkflowDefinition `json:"workflow" gorm:"foreignKey:WorkflowID"` // 流程关联
	StepName     string             `json:"step_name" gorm:"size:100;not null"`    // 步骤名称
	StepOrder    int                `json:"step_order" gorm:"not null"`            // 步骤顺序
	ApproverType string             `json:"approver_type" gorm:"size:50;not null"` // 审批人类型: user, role, space_admin
	ApproverID   uint               `json:"approver_id"`                           // 审批人ID（用户ID或角色ID）
	IsRequired   bool               `json:"is_required" gorm:"default:true"`       // 是否必须
	TimeoutHours int                `json:"timeout_hours" gorm:"default:72"`       // 超时时间（小时）
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// WorkflowInstance 审批流程实例
type WorkflowInstance struct {
	ID           uint               `json:"id" gorm:"primaryKey"`
	WorkflowID   uint               `json:"workflow_id" gorm:"not null"`              // 流程定义ID
	Workflow     WorkflowDefinition `json:"workflow" gorm:"foreignKey:WorkflowID"`    // 流程定义关联
	Title        string             `json:"title" gorm:"size:200;not null"`           // 审批标题
	Description  string             `json:"description" gorm:"size:1000"`             // 审批描述
	ResourceType string             `json:"resource_type" gorm:"size:50;not null"`    // 资源类型: document, space, user
	ResourceID   uint               `json:"resource_id" gorm:"not null"`              // 资源ID
	SpaceID      uint               `json:"space_id" gorm:"not null"`                 // 所属空间ID
	Status       string             `json:"status" gorm:"size:20;default:'pending'"`  // 状态: pending, approved, rejected, cancelled, timeout
	Priority     string             `json:"priority" gorm:"size:20;default:'normal'"` // 优先级: normal, urgent
	CreatedBy    uint               `json:"created_by" gorm:"not null"`               // 申请人
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	CompletedAt  *time.Time         `json:"completed_at"`                       // 完成时间
	Tasks        []WorkflowTask     `json:"tasks" gorm:"foreignKey:InstanceID"` // 审批任务
}

// WorkflowTask 审批任务
type WorkflowTask struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	InstanceID  uint             `json:"instance_id" gorm:"not null"`             // 流程实例ID
	Instance    WorkflowInstance `json:"instance" gorm:"foreignKey:InstanceID"`   // 流程实例关联
	StepID      uint             `json:"step_id" gorm:"not null"`                 // 流程步骤ID
	Step        WorkflowStep     `json:"step" gorm:"foreignKey:StepID"`           // 流程步骤关联
	AssigneeID  uint             `json:"assignee_id" gorm:"not null"`             // 审批人ID
	Status      string           `json:"status" gorm:"size:20;default:'pending'"` // 状态: pending, approved, rejected, transferred
	Comment     string           `json:"comment" gorm:"size:1000"`                // 审批意见
	AssignedAt  time.Time        `json:"assigned_at"`                             // 分配时间
	CompletedAt *time.Time       `json:"completed_at"`                            // 完成时间
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// WorkflowNotification 审批通知
type WorkflowNotification struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	TaskID    uint         `json:"task_id" gorm:"not null"`        // 任务ID
	Task      WorkflowTask `json:"task" gorm:"foreignKey:TaskID"`  // 任务关联
	UserID    uint         `json:"user_id" gorm:"not null"`        // 接收人ID
	Type      string       `json:"type" gorm:"size:50;not null"`   // 通知类型: task_assigned, task_approved, task_rejected, task_timeout
	Title     string       `json:"title" gorm:"size:200;not null"` // 通知标题
	Content   string       `json:"content" gorm:"size:1000"`       // 通知内容
	IsRead    bool         `json:"is_read" gorm:"default:false"`   // 是否已读
	CreatedAt time.Time    `json:"created_at"`
	ReadAt    *time.Time   `json:"read_at"` // 阅读时间
}

// 审批状态常量
const (
	StatusPending   = "pending"   // 待审批
	StatusApproved  = "approved"  // 已通过
	StatusRejected  = "rejected"  // 已拒绝
	StatusCancelled = "cancelled" // 已取消
	StatusTimeout   = "timeout"   // 已超时
)

// 优先级常量
const (
	PriorityNormal = "normal" // 一般
	PriorityUrgent = "urgent" // 紧急
)

// 审批人类型常量
const (
	ApproverTypeUser       = "user"        // 指定用户
	ApproverTypeRole       = "role"        // 指定角色
	ApproverTypeSpaceAdmin = "space_admin" // 空间管理员
)

// 资源类型常量
const (
	ResourceTypeDocument = "document" // 文档
	ResourceTypeSpace    = "space"    // 空间
	ResourceTypeUser     = "user"     // 用户
)

// 通知类型常量
const (
	NotificationTypeTaskAssigned = "task_assigned" // 任务分配
	NotificationTypeTaskApproved = "task_approved" // 任务通过
	NotificationTypeTaskRejected = "task_rejected" // 任务拒绝
	NotificationTypeTaskTimeout  = "task_timeout"  // 任务超时
)

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
	StepName     string `json:"step_name" binding:"required"`
	StepOrder    int    `json:"step_order" binding:"required"`
	ApproverType string `json:"approver_type" binding:"required"`
	ApproverID   uint   `json:"approver_id"`
	IsRequired   bool   `json:"is_required"`
	TimeoutHours int    `json:"timeout_hours"`
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
	NewAssigneeID uint   `json:"new_assignee_id" binding:"required"`
	Comment       string `json:"comment"`
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
