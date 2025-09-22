package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/gideonzy/knowledge-base/internal/workflow"
	"github.com/gideonzy/knowledge-base/internal/workflow/repository"
)

// Workflow encapsulates workflow orchestration logic.
type Workflow struct {
	Definitions repository.DefinitionRepository
	Instances   repository.InstanceRepository
}

// New creates a new workflow service.
func New(defRepo repository.DefinitionRepository, instRepo repository.InstanceRepository) *Workflow {
	return &Workflow{Definitions: defRepo, Instances: instRepo}
}

// RegisterDefinition registers or updates a flow definition.
func (s *Workflow) RegisterDefinition(code, name, description string, nodes []workflow.FlowNode) (workflow.FlowDefinition, error) {
	if strings.TrimSpace(code) == "" {
		return workflow.FlowDefinition{}, errors.New("code is required")
	}
	if len(nodes) == 0 {
		return workflow.FlowDefinition{}, errors.New("at least one node is required")
	}
	preparedNodes := make([]workflow.FlowNode, len(nodes))
	for i, node := range nodes {
		if node.ID == "" {
			node.ID = generateID()
		}
		if node.Type == "" {
			node.Type = workflow.NodeTypeApproval
		}
		preparedNodes[i] = node
	}

	now := time.Now().UTC()
	existing, found := s.Definitions.GetByCode(code)
	def := workflow.FlowDefinition{
		ID:          generateID(),
		Code:        code,
		Name:        name,
		Description: description,
		Nodes:       preparedNodes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if found {
		def.ID = existing.ID
		def.CreatedAt = existing.CreatedAt
	}
	if err := s.Definitions.Save(def); err != nil {
		return workflow.FlowDefinition{}, err
	}
	return def, nil
}

// ListDefinitions returns all defined workflows.
func (s *Workflow) ListDefinitions() []workflow.FlowDefinition {
	return s.Definitions.List()
}

// StartInstance launches a workflow instance.
func (s *Workflow) StartInstance(definitionCode, businessID, spaceID, createdBy string) (workflow.FlowInstance, error) {
	def, ok := s.Definitions.GetByCode(definitionCode)
	if !ok {
		return workflow.FlowInstance{}, errors.New("definition not found")
	}
	if businessID == "" {
		return workflow.FlowInstance{}, errors.New("business id required")
	}
	firstNodeID := ""
	if len(def.Nodes) > 0 {
		firstNodeID = def.Nodes[0].ID
	}
	now := time.Now().UTC()
	inst := workflow.FlowInstance{
		ID:            generateID(),
		DefinitionID:  def.ID,
		BusinessID:    businessID,
		SpaceID:       spaceID,
		Status:        workflow.StatusPending,
		CurrentNodeID: firstNodeID,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.Instances.Save(inst); err != nil {
		return workflow.FlowInstance{}, err
	}
	return inst, nil
}

// ApplyAction updates workflow status for the current node.
func (s *Workflow) ApplyAction(instanceID, actorID, comment string, action workflow.TaskAction) (workflow.FlowInstance, error) {
	inst, ok := s.Instances.Get(instanceID)
	if !ok {
		return workflow.FlowInstance{}, errors.New("instance not found")
	}
	if inst.Status != workflow.StatusPending {
		return workflow.FlowInstance{}, errors.New("instance already closed")
	}

	definition, ok := s.Definitions.GetByID(inst.DefinitionID)
	if !ok {
		return workflow.FlowInstance{}, errors.New("definition not found")
	}

	now := time.Now().UTC()
	historyEntry := workflow.InstanceAction{
		ID:         generateID(),
		InstanceID: inst.ID,
		NodeID:     inst.CurrentNodeID,
		ActorID:    actorID,
		Action:     action,
		Comment:    comment,
		CreatedAt:  now,
	}

	inst.History = append(inst.History, historyEntry)

	switch action {
	case workflow.ActionApprove:
		nextNodeID := findNextNode(definition, inst.CurrentNodeID)
		if nextNodeID == "" {
			inst.Status = workflow.StatusApproved
			inst.CurrentNodeID = ""
			inst.CompletedAt = &now
		} else {
			inst.CurrentNodeID = nextNodeID
		}
	case workflow.ActionReject:
		inst.Status = workflow.StatusRejected
		inst.CurrentNodeID = ""
		inst.CompletedAt = &now
	default:
		return workflow.FlowInstance{}, errors.New("unsupported action")
	}

	inst.UpdatedAt = now
	if err := s.Instances.Save(inst); err != nil {
		return workflow.FlowInstance{}, err
	}
	return inst, nil
}

// GetInstance returns a workflow instance.
func (s *Workflow) GetInstance(id string) (workflow.FlowInstance, error) {
	inst, ok := s.Instances.Get(id)
	if !ok {
		return workflow.FlowInstance{}, errors.New("instance not found")
	}
	return inst, nil
}

func findNextNode(def workflow.FlowDefinition, currentID string) string {
	for _, node := range def.Nodes {
		if node.ID == currentID {
			return node.NextNodeID
		}
	}
	return ""
}

func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
