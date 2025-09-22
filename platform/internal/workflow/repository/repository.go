package repository

import "gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow"

// DefinitionRepository handles workflow definitions.
type DefinitionRepository interface {
	GetByCode(code string) (workflow.FlowDefinition, bool)
	GetByID(id string) (workflow.FlowDefinition, bool)
	Save(def workflow.FlowDefinition) error
	List() []workflow.FlowDefinition
}

// InstanceRepository handles workflow instances.
type InstanceRepository interface {
	Get(id string) (workflow.FlowInstance, bool)
	Save(instance workflow.FlowInstance) error
	ListByDefinition(definitionID string) []workflow.FlowInstance
}
