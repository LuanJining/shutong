package repository

import (
	"database/sql"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/iam"
	"github.com/lib/pq"
)

// PostgresRepositories aggregates SQL-backed repositories.
type PostgresRepositories struct {
	Users    *PostgresUserRepo
	Roles    *PostgresRoleRepo
	Spaces   *PostgresSpaceRepo
	Policies *PostgresPolicyRepo
}

// NewPostgresRepositories wires Postgres repositories.
func NewPostgresRepositories(db *sql.DB) *PostgresRepositories {
	return &PostgresRepositories{
		Users:    &PostgresUserRepo{db: db},
		Roles:    &PostgresRoleRepo{db: db},
		Spaces:   &PostgresSpaceRepo{db: db},
		Policies: &PostgresPolicyRepo{db: db},
	}
}

// PostgresUserRepo implements UserRepository using PostgreSQL.
type PostgresUserRepo struct{ db *sql.DB }

func (r *PostgresUserRepo) Get(id string) (iam.User, bool) {
	const query = `SELECT id, name, phone, password_hash, roles, spaces, created_at, updated_at FROM iam_users WHERE id=$1`
	var user iam.User
	var roles pq.StringArray
	var spaces pq.StringArray
	if err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Phone, &user.PasswordHash, &roles, &spaces, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return iam.User{}, false
		}
		return iam.User{}, false
	}
	user.Roles = []string(roles)
	user.Spaces = []string(spaces)
	return user, true
}

// GetByPhone fetches a user by phone number.
func (r *PostgresUserRepo) GetByPhone(phone string) (iam.User, bool) {
	const query = `SELECT id, name, phone, password_hash, roles, spaces, created_at, updated_at FROM iam_users WHERE phone=$1`
	var user iam.User
	var roles pq.StringArray
	var spaces pq.StringArray
	if err := r.db.QueryRow(query, phone).Scan(&user.ID, &user.Name, &user.Phone, &user.PasswordHash, &roles, &spaces, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return iam.User{}, false
		}
		return iam.User{}, false
	}
	user.Roles = []string(roles)
	user.Spaces = []string(spaces)
	return user, true
}

func (r *PostgresUserRepo) List() []iam.User {
	const query = `SELECT id, name, phone, password_hash, roles, spaces, created_at, updated_at FROM iam_users ORDER BY created_at`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	users := make([]iam.User, 0)
	for rows.Next() {
		var user iam.User
		var roles pq.StringArray
		var spaces pq.StringArray
		if err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.PasswordHash, &roles, &spaces, &user.CreatedAt, &user.UpdatedAt); err != nil {
			continue
		}
		user.Roles = []string(roles)
		user.Spaces = []string(spaces)
		users = append(users, user)
	}
	return users
}

func (r *PostgresUserRepo) Save(user iam.User) error {
	const query = `INSERT INTO iam_users (id, name, phone, password_hash, roles, spaces, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, phone=EXCLUDED.phone, password_hash=EXCLUDED.password_hash, roles=EXCLUDED.roles, spaces=EXCLUDED.spaces, updated_at=EXCLUDED.updated_at`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Phone, user.PasswordHash, pq.StringArray(user.Roles), pq.StringArray(user.Spaces), user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM iam_users WHERE id=$1`, id)
	return err
}

// PostgresRoleRepo implements RoleRepository.
type PostgresRoleRepo struct{ db *sql.DB }

func (r *PostgresRoleRepo) Get(id string) (iam.Role, bool) {
	const query = `SELECT id, name, description, permissions, created_at, updated_at FROM iam_roles WHERE id=$1`
	var role iam.Role
	var permissions pq.StringArray
	if err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &permissions, &role.CreatedAt, &role.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return iam.Role{}, false
		}
		return iam.Role{}, false
	}
	role.Permissions = []string(permissions)
	return role, true
}

func (r *PostgresRoleRepo) List() []iam.Role {
	const query = `SELECT id, name, description, permissions, created_at, updated_at FROM iam_roles ORDER BY created_at`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	roles := make([]iam.Role, 0)
	for rows.Next() {
		var role iam.Role
		var permissions pq.StringArray
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &permissions, &role.CreatedAt, &role.UpdatedAt); err != nil {
			continue
		}
		role.Permissions = []string(permissions)
		roles = append(roles, role)
	}
	return roles
}

func (r *PostgresRoleRepo) Save(role iam.Role) error {
	const query = `INSERT INTO iam_roles (id, name, description, permissions, created_at, updated_at)
                    VALUES ($1,$2,$3,$4,$5,$6)
                    ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, description=EXCLUDED.description, permissions=EXCLUDED.permissions, updated_at=EXCLUDED.updated_at`
	_, err := r.db.Exec(query, role.ID, role.Name, role.Description, pq.StringArray(role.Permissions), role.CreatedAt, role.UpdatedAt)
	return err
}

func (r *PostgresRoleRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM iam_roles WHERE id=$1`, id)
	return err
}

// PostgresSpaceRepo implements SpaceRepository.
type PostgresSpaceRepo struct{ db *sql.DB }

func (r *PostgresSpaceRepo) Get(id string) (iam.Space, bool) {
	const query = `SELECT id, name, description, created_at, updated_at FROM iam_spaces WHERE id=$1`
	var space iam.Space
	if err := r.db.QueryRow(query, id).Scan(&space.ID, &space.Name, &space.Description, &space.CreatedAt, &space.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return iam.Space{}, false
		}
		return iam.Space{}, false
	}
	return space, true
}

func (r *PostgresSpaceRepo) List() []iam.Space {
	const query = `SELECT id, name, description, created_at, updated_at FROM iam_spaces ORDER BY created_at`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	spaces := make([]iam.Space, 0)
	for rows.Next() {
		var space iam.Space
		if err := rows.Scan(&space.ID, &space.Name, &space.Description, &space.CreatedAt, &space.UpdatedAt); err != nil {
			continue
		}
		spaces = append(spaces, space)
	}
	return spaces
}

func (r *PostgresSpaceRepo) Save(space iam.Space) error {
	const query = `INSERT INTO iam_spaces (id, name, description, created_at, updated_at)
                    VALUES ($1,$2,$3,$4,$5)
                    ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, description=EXCLUDED.description, updated_at=EXCLUDED.updated_at`
	_, err := r.db.Exec(query, space.ID, space.Name, space.Description, space.CreatedAt, space.UpdatedAt)
	return err
}

func (r *PostgresSpaceRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM iam_spaces WHERE id=$1`, id)
	return err
}

// PostgresPolicyRepo implements PolicyRepository.
type PostgresPolicyRepo struct{ db *sql.DB }

func (r *PostgresPolicyRepo) ListBySpace(spaceID string) []iam.Policy {
	const query = `SELECT id, space_id, role_id, resource, action, created_at FROM iam_policies WHERE space_id=$1`
	rows, err := r.db.Query(query, spaceID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	policies := make([]iam.Policy, 0)
	for rows.Next() {
		var policy iam.Policy
		if err := rows.Scan(&policy.ID, &policy.SpaceID, &policy.RoleID, &policy.Resource, &policy.Action, &policy.CreatedAt); err != nil {
			continue
		}
		policies = append(policies, policy)
	}
	return policies
}

func (r *PostgresPolicyRepo) Save(policy iam.Policy) error {
	const query = `INSERT INTO iam_policies (id, space_id, role_id, resource, action, created_at)
                    VALUES ($1,$2,$3,$4,$5,$6)
                    ON CONFLICT (id) DO UPDATE SET space_id=EXCLUDED.space_id, role_id=EXCLUDED.role_id, resource=EXCLUDED.resource, action=EXCLUDED.action`
	_, err := r.db.Exec(query, policy.ID, policy.SpaceID, policy.RoleID, policy.Resource, policy.Action, policy.CreatedAt)
	return err
}

func (r *PostgresPolicyRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM iam_policies WHERE id=$1`, id)
	return err
}
