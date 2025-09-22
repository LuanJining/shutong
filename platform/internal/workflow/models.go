package workflow

import "time"

// FlowNodeType enumerates node behaviours.
type FlowNodeType string

const (
	// NodeTypeApproval indicates the node requires manual approval.
	NodeTypeApproval FlowNodeType = "approval"
	// NodeTypeAuto indicates the node auto-transitions.
	NodeTypeAuto FlowNodeType = "auto"
)

// FlowDefinition describes a workflow template for approvals.
type FlowDefinition struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Nodes       []FlowNode `json:"nodes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FlowNode represents a single node in a workflow.
type FlowNode struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         FlowNodeType `json:"type"`
	ApproverRole string       `json:"approver_role"`
	NextNodeID   string       `json:"next_node_id"`
}

// FlowInstance tracks a running workflow instance.
type FlowInstance struct {
	ID            string           `json:"id"`
	DefinitionID  string           `json:"definition_id"`
	BusinessID    string           `json:"business_id"`
	SpaceID       string           `json:"space_id"`
	Status        InstanceStatus   `json:"status"`
	CurrentNodeID string           `json:"current_node_id"`
	CreatedBy     string           `json:"created_by"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	CompletedAt   *time.Time       `json:"completed_at,omitempty"`
	History       []InstanceAction `json:"history"`
}

// InstanceAction captures a decision made on a task.
type InstanceAction struct {
	ID         string     `json:"id"`
	InstanceID string     `json:"instance_id"`
	NodeID     string     `json:"node_id"`
	ActorID    string     `json:"actor_id"`
	Action     TaskAction `json:"action"`
	Comment    string     `json:"comment"`
	CreatedAt  time.Time  `json:"created_at"`
}

// InstanceStatus represents workflow state.
type InstanceStatus string

const (
	// StatusPending indicates active workflow.
	StatusPending InstanceStatus = "pending"
	// StatusApproved indicates successful completion.
	StatusApproved InstanceStatus = "approved"
	// StatusRejected indicates the request was rejected.
	StatusRejected InstanceStatus = "rejected"
)

// TaskAction enumerates transition actions.
type TaskAction string

const (
	// ActionApprove approves a task.
	ActionApprove TaskAction = "approve"
	// ActionReject rejects a task.
	ActionReject TaskAction = "reject"
)
