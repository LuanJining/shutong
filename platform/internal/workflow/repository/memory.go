package repository

import (
	"errors"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/common/storage"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/workflow"
)

// DefinitionRepo is an in-memory FlowDefinition repo.
type DefinitionRepo struct {
	store *storage.InMemory[workflow.FlowDefinition]
}

// NewDefinitionRepo creates a new repo instance.
func NewDefinitionRepo(store *storage.InMemory[workflow.FlowDefinition]) *DefinitionRepo {
	return &DefinitionRepo{store: store}
}

// GetByCode retrieves definition by code.
func (r *DefinitionRepo) GetByCode(code string) (workflow.FlowDefinition, bool) {
	defs := r.store.List()
	for _, def := range defs {
		if def.Code == code {
			return def, true
		}
	}
	return workflow.FlowDefinition{}, false
}

// GetByID retrieves a definition by id.
func (r *DefinitionRepo) GetByID(id string) (workflow.FlowDefinition, bool) {
	return r.store.Get(id)
}

// Save stores definition.
func (r *DefinitionRepo) Save(def workflow.FlowDefinition) error {
	if def.ID == "" {
		return errors.New("definition id required")
	}
	r.store.Set(def.ID, def)
	return nil
}

// List returns definitions.
func (r *DefinitionRepo) List() []workflow.FlowDefinition {
	return r.store.List()
}

// InstanceRepo is an in-memory FlowInstance repo.
type InstanceRepo struct {
	store *storage.InMemory[workflow.FlowInstance]
}

// NewInstanceRepo creates repo.
func NewInstanceRepo(store *storage.InMemory[workflow.FlowInstance]) *InstanceRepo {
	return &InstanceRepo{store: store}
}

// Get retrieves instance by id.
func (r *InstanceRepo) Get(id string) (workflow.FlowInstance, bool) {
	return r.store.Get(id)
}

// Save stores instance.
func (r *InstanceRepo) Save(instance workflow.FlowInstance) error {
	if instance.ID == "" {
		return errors.New("instance id required")
	}
	r.store.Set(instance.ID, instance)
	return nil
}

// ListByDefinition returns instances filtered by definition.
func (r *InstanceRepo) ListByDefinition(definitionID string) []workflow.FlowInstance {
	list := r.store.List()
	filtered := make([]workflow.FlowInstance, 0)
	for _, inst := range list {
		if inst.DefinitionID == definitionID {
			filtered = append(filtered, inst)
		}
	}
	return filtered
}
