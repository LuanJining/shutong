package repository

import (
	"errors"
	"fmt"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/storage"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam"
)

// InMemoryRepositories aggregates memory-backed repos for IAM service.
type InMemoryRepositories struct {
	Users    *storage.InMemory[iam.User]
	Roles    *storage.InMemory[iam.Role]
	Spaces   *storage.InMemory[iam.Space]
	Policies *storage.InMemory[iam.Policy]
}

// NewInMemoryRepositories creates repository instances.
func NewInMemoryRepositories() *InMemoryRepositories {
	return &InMemoryRepositories{
		Users:    storage.NewInMemory[iam.User](),
		Roles:    storage.NewInMemory[iam.Role](),
		Spaces:   storage.NewInMemory[iam.Space](),
		Policies: storage.NewInMemory[iam.Policy](),
	}
}

// UserRepo implements UserRepository.
type UserRepo struct{ store *storage.InMemory[iam.User] }

// NewUserRepo returns a new UserRepo.
func NewUserRepo(store *storage.InMemory[iam.User]) *UserRepo {
	return &UserRepo{store: store}
}

// Get retrieves a user by id.
func (r *UserRepo) Get(id string) (iam.User, bool) {
	return r.store.Get(id)
}

// GetByPhone retrieves a user by phone number.
func (r *UserRepo) GetByPhone(phone string) (iam.User, bool) {
	users := r.store.List()
	for _, user := range users {
		if user.Phone == phone {
			return user, true
		}
	}
	return iam.User{}, false
}

// List returns all users.
func (r *UserRepo) List() []iam.User {
	return r.store.List()
}

// Save stores the user.
func (r *UserRepo) Save(user iam.User) error {
	if user.ID == "" {
		return errors.New("user id required")
	}
	r.store.Set(user.ID, user)
	return nil
}

// Delete removes a user.
func (r *UserRepo) Delete(id string) error {
	r.store.Delete(id)
	return nil
}

// RoleRepo implements RoleRepository.
type RoleRepo struct{ store *storage.InMemory[iam.Role] }

// NewRoleRepo returns a new RoleRepo.
func NewRoleRepo(store *storage.InMemory[iam.Role]) *RoleRepo {
	return &RoleRepo{store: store}
}

func (r *RoleRepo) Get(id string) (iam.Role, bool) { return r.store.Get(id) }
func (r *RoleRepo) List() []iam.Role               { return r.store.List() }
func (r *RoleRepo) Save(role iam.Role) error {
	if role.ID == "" {
		return errors.New("role id required")
	}
	r.store.Set(role.ID, role)
	return nil
}
func (r *RoleRepo) Delete(id string) error { r.store.Delete(id); return nil }

// SpaceRepo implements SpaceRepository.
type SpaceRepo struct{ store *storage.InMemory[iam.Space] }

// NewSpaceRepo returns a new SpaceRepo.
func NewSpaceRepo(store *storage.InMemory[iam.Space]) *SpaceRepo {
	return &SpaceRepo{store: store}
}

func (r *SpaceRepo) Get(id string) (iam.Space, bool) { return r.store.Get(id) }
func (r *SpaceRepo) List() []iam.Space               { return r.store.List() }
func (r *SpaceRepo) Save(space iam.Space) error {
	if space.ID == "" {
		return errors.New("space id required")
	}
	r.store.Set(space.ID, space)
	return nil
}
func (r *SpaceRepo) Delete(id string) error { r.store.Delete(id); return nil }

// PolicyRepo implements PolicyRepository.
type PolicyRepo struct{ store *storage.InMemory[iam.Policy] }

// NewPolicyRepo returns a new PolicyRepo.
func NewPolicyRepo(store *storage.InMemory[iam.Policy]) *PolicyRepo {
	return &PolicyRepo{store: store}
}

func (r *PolicyRepo) ListBySpace(spaceID string) []iam.Policy {
	policies := r.store.List()
	filtered := make([]iam.Policy, 0)
	for _, p := range policies {
		if p.SpaceID == spaceID {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (r *PolicyRepo) Save(policy iam.Policy) error {
	if policy.ID == "" {
		return errors.New("policy id required")
	}
	if policy.SpaceID == "" {
		return fmt.Errorf("policy %s missing space id", policy.ID)
	}
	r.store.Set(policy.ID, policy)
	return nil
}

func (r *PolicyRepo) Delete(id string) error {
	r.store.Delete(id)
	return nil
}
