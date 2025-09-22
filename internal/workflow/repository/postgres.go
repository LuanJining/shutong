package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/gideonzy/knowledge-base/internal/workflow"
)

// PostgresDefinitionRepo manages flow definitions in Postgres.
type PostgresDefinitionRepo struct{ db *sql.DB }

// NewPostgresDefinitionRepo creates a repo.
func NewPostgresDefinitionRepo(db *sql.DB) *PostgresDefinitionRepo {
	return &PostgresDefinitionRepo{db: db}
}

func (r *PostgresDefinitionRepo) GetByCode(code string) (workflow.FlowDefinition, bool) {
	const query = `SELECT id, code, name, description, nodes, created_at, updated_at FROM wf_definitions WHERE code=$1`
	var (
		def       workflow.FlowDefinition
		nodesData []byte
	)
	if err := r.db.QueryRow(query, code).Scan(&def.ID, &def.Code, &def.Name, &def.Description, &nodesData, &def.CreatedAt, &def.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return workflow.FlowDefinition{}, false
		}
		return workflow.FlowDefinition{}, false
	}
	_ = json.Unmarshal(nodesData, &def.Nodes)
	return def, true
}

func (r *PostgresDefinitionRepo) GetByID(id string) (workflow.FlowDefinition, bool) {
	const query = `SELECT id, code, name, description, nodes, created_at, updated_at FROM wf_definitions WHERE id=$1`
	var (
		def       workflow.FlowDefinition
		nodesData []byte
	)
	if err := r.db.QueryRow(query, id).Scan(&def.ID, &def.Code, &def.Name, &def.Description, &nodesData, &def.CreatedAt, &def.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return workflow.FlowDefinition{}, false
		}
		return workflow.FlowDefinition{}, false
	}
	_ = json.Unmarshal(nodesData, &def.Nodes)
	return def, true
}

func (r *PostgresDefinitionRepo) Save(def workflow.FlowDefinition) error {
	nodesData, err := json.Marshal(def.Nodes)
	if err != nil {
		return err
	}
	const query = `INSERT INTO wf_definitions (id, code, name, description, nodes, created_at, updated_at)
                    VALUES ($1,$2,$3,$4,$5,$6,$7)
                    ON CONFLICT (code) DO UPDATE SET name=EXCLUDED.name, description=EXCLUDED.description, nodes=EXCLUDED.nodes, updated_at=EXCLUDED.updated_at`
	_, err = r.db.Exec(query, def.ID, def.Code, def.Name, def.Description, nodesData, def.CreatedAt, def.UpdatedAt)
	return err
}

func (r *PostgresDefinitionRepo) List() []workflow.FlowDefinition {
	const query = `SELECT id, code, name, description, nodes, created_at, updated_at FROM wf_definitions ORDER BY created_at`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	defs := make([]workflow.FlowDefinition, 0)
	for rows.Next() {
		var def workflow.FlowDefinition
		var nodesData []byte
		if err := rows.Scan(&def.ID, &def.Code, &def.Name, &def.Description, &nodesData, &def.CreatedAt, &def.UpdatedAt); err != nil {
			continue
		}
		_ = json.Unmarshal(nodesData, &def.Nodes)
		defs = append(defs, def)
	}
	return defs
}

// PostgresInstanceRepo stores workflow instances.
type PostgresInstanceRepo struct{ db *sql.DB }

// NewPostgresInstanceRepo creates repo.
func NewPostgresInstanceRepo(db *sql.DB) *PostgresInstanceRepo {
	return &PostgresInstanceRepo{db: db}
}

func (r *PostgresInstanceRepo) Get(id string) (workflow.FlowInstance, bool) {
	const query = `SELECT id, definition_id, business_id, space_id, status, current_node_id, created_by, created_at, updated_at, completed_at
                    FROM wf_instances WHERE id=$1`
	var (
		inst      workflow.FlowInstance
		completed sql.NullTime
	)
	if err := r.db.QueryRow(query, id).Scan(&inst.ID, &inst.DefinitionID, &inst.BusinessID, &inst.SpaceID, &inst.Status, &inst.CurrentNodeID, &inst.CreatedBy, &inst.CreatedAt, &inst.UpdatedAt, &completed); err != nil {
		if err == sql.ErrNoRows {
			return workflow.FlowInstance{}, false
		}
		return workflow.FlowInstance{}, false
	}
	if completed.Valid {
		inst.CompletedAt = &completed.Time
	}
	inst.History = r.loadHistory(inst.ID)
	return inst, true
}

func (r *PostgresInstanceRepo) Save(inst workflow.FlowInstance) error {
	const query = `INSERT INTO wf_instances (id, definition_id, business_id, space_id, status, current_node_id, created_by, created_at, updated_at, completed_at)
                    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
                    ON CONFLICT (id) DO UPDATE SET status=EXCLUDED.status, current_node_id=EXCLUDED.current_node_id, updated_at=EXCLUDED.updated_at, completed_at=EXCLUDED.completed_at`
	var completed interface{}
	if inst.CompletedAt != nil {
		completed = *inst.CompletedAt
	} else {
		completed = nil
	}
	if _, err := r.db.Exec(query, inst.ID, inst.DefinitionID, inst.BusinessID, inst.SpaceID, inst.Status, inst.CurrentNodeID, inst.CreatedBy, inst.CreatedAt, inst.UpdatedAt, completed); err != nil {
		return err
	}
	if len(inst.History) > 0 {
		if err := r.saveHistory(inst.History[len(inst.History)-1]); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresInstanceRepo) ListByDefinition(definitionID string) []workflow.FlowInstance {
	const query = `SELECT id, definition_id, business_id, space_id, status, current_node_id, created_by, created_at, updated_at, completed_at
                    FROM wf_instances WHERE definition_id=$1`
	rows, err := r.db.Query(query, definitionID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	instances := make([]workflow.FlowInstance, 0)
	for rows.Next() {
		var (
			inst      workflow.FlowInstance
			completed sql.NullTime
		)
		if err := rows.Scan(&inst.ID, &inst.DefinitionID, &inst.BusinessID, &inst.SpaceID, &inst.Status, &inst.CurrentNodeID, &inst.CreatedBy, &inst.CreatedAt, &inst.UpdatedAt, &completed); err != nil {
			continue
		}
		if completed.Valid {
			inst.CompletedAt = &completed.Time
		}
		inst.History = r.loadHistory(inst.ID)
		instances = append(instances, inst)
	}
	return instances
}

func (r *PostgresInstanceRepo) saveHistory(entry workflow.InstanceAction) error {
	const query = `INSERT INTO wf_history (id, instance_id, node_id, actor_id, action, comment, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(query, entry.ID, entry.InstanceID, entry.NodeID, entry.ActorID, entry.Action, entry.Comment, entry.CreatedAt)
	return err
}

func (r *PostgresInstanceRepo) loadHistory(instanceID string) []workflow.InstanceAction {
	const query = `SELECT id, instance_id, node_id, actor_id, action, comment, created_at FROM wf_history WHERE instance_id=$1 ORDER BY created_at`
	rows, err := r.db.Query(query, instanceID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	history := make([]workflow.InstanceAction, 0)
	for rows.Next() {
		var entry workflow.InstanceAction
		if err := rows.Scan(&entry.ID, &entry.InstanceID, &entry.NodeID, &entry.ActorID, &entry.Action, &entry.Comment, &entry.CreatedAt); err != nil {
			continue
		}
		history = append(history, entry)
	}
	return history
}
