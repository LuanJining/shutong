package repository

import "github.com/gideonzy/knowledge-base/internal/iam"

// UserRepository manages users.
type UserRepository interface {
	Get(id string) (iam.User, bool)
	List() []iam.User
	Save(user iam.User) error
	Delete(id string) error
}

// RoleRepository manages roles.
type RoleRepository interface {
	Get(id string) (iam.Role, bool)
	List() []iam.Role
	Save(role iam.Role) error
	Delete(id string) error
}

// SpaceRepository manages spaces.
type SpaceRepository interface {
	Get(id string) (iam.Space, bool)
	List() []iam.Space
	Save(space iam.Space) error
	Delete(id string) error
}

// PolicyRepository manages resource policies.
type PolicyRepository interface {
	ListBySpace(spaceID string) []iam.Policy
	Save(policy iam.Policy) error
	Delete(id string) error
}
