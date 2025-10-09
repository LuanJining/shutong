package model

type Workflow struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" binding:"required"`
	Description   string         `json:"description"`
	SpaceID       uint           `json:"space_id" binding:"required"`
	Status        WorkflowStatus `json:"status" binding:"required"`
	CurrentStepID uint           `json:"current_step_id"` // 当前步骤ID

	// 关联字段
	CreatedBy       uint   `json:"created_by"`
	CreatorNickName string `json:"creator_nick_name"`
	Steps           []Step `json:"steps" gorm:"foreignKey:WorkflowID"` // 一对多关联
}

type Step struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	StepName     string     `json:"step_name" binding:"required"`
	StepOrder    int        `json:"step_order" binding:"required"`
	StepRole     string     `json:"step_role" binding:"required"`
	IsRequired   bool       `json:"is_required" binding:"required"`
	TimeoutHours int        `json:"timeout_hours" binding:"required"`
	Tasks        []Task     `json:"tasks" gorm:"foreignKey:StepID"` // 一对多关联
	Status       StepStatus `json:"status" binding:"required"`

	// 关联字段
	WorkflowID uint `json:"workflow_id"`
}

type Task struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	TaskName         string     `json:"task_name" binding:"required"`
	IsRequired       bool       `json:"is_required" binding:"required"`
	TimeoutHours     int        `json:"timeout_hours" binding:"required"`
	Status           TaskStatus `json:"status" binding:"required"`
	ApproverID       uint       `json:"approver_id"`
	ApproverNickName string     `json:"approver_nick_name"`

	// 关联字段
	WorkflowID uint `json:"workflow_id"`
	StepID     uint `json:"step_id"`
}

type WorkflowStatus string

const (
	WorkflowStatusProcessing WorkflowStatus = "processing"
	WorkflowStatusCompleted  WorkflowStatus = "completed"
	WorkflowStatusCancelled  WorkflowStatus = "cancelled"
)

type StepStatus string

const (
	StepStatusProcessing StepStatus = "processing"
	StepStatusApproved   StepStatus = "approved"
	StepStatusRejected   StepStatus = "rejected"
)

type TaskStatus string

const (
	TaskStatusProcessing      TaskStatus = "processing"
	TaskStatusApproved        TaskStatus = "approved"
	TaskStatusApprovedByOther TaskStatus = "approved_by_others"
	TaskStatusRejected        TaskStatus = "rejected"
	TaskStatusRejectedByOther TaskStatus = "rejected_by_others"
)

type CreateWorkflowRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SpaceID     uint   `json:"space_id" binding:"required"`
	Priority    int    `json:"priority"`
	Steps       []Step `json:"steps" binding:"required"`
}

type StartWorkflowRequest struct {
	WorkflowID uint `json:"workflow_id" binding:"required"`
}
