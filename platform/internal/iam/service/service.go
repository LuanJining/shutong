package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/iam"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/platform/internal/iam/repository"
)

// IAM provides core account management operations.
type IAM struct {
	Users    repository.UserRepository
	Roles    repository.RoleRepository
	Spaces   repository.SpaceRepository
	Policies repository.PolicyRepository
}

// New creates a new IAM service.
func New(users repository.UserRepository, roles repository.RoleRepository, spaces repository.SpaceRepository, policies repository.PolicyRepository) *IAM {
	return &IAM{Users: users, Roles: roles, Spaces: spaces, Policies: policies}
}

// CreateUser creates a new user.
func (s *IAM) CreateUser(name, phone, password string, roles, spaces []string) (iam.User, error) {
	if strings.TrimSpace(name) == "" {
		return iam.User{}, errors.New("name is required")
	}
	if strings.TrimSpace(phone) == "" {
		return iam.User{}, errors.New("phone is required")
	}
	if len(password) < 6 {
		return iam.User{}, errors.New("password must be at least 6 characters")
	}
	phone = normalizePhone(phone)
	if _, exists := s.Users.GetByPhone(phone); exists {
		return iam.User{}, errors.New("phone already registered")
	}
	hash, err := hashPassword(password)
	if err != nil {
		return iam.User{}, err
	}
	now := time.Now().UTC()
	user := iam.User{
		ID:           generateID(),
		Name:         name,
		Phone:        phone,
		PasswordHash: hash,
		Roles:        uniqueStrings(roles),
		Spaces:       uniqueStrings(spaces),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.Users.Save(user); err != nil {
		return iam.User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user.
func (s *IAM) UpdateUser(id string, name, phone, password string, roles, spaces []string) (iam.User, error) {
	existing, ok := s.Users.Get(id)
	if !ok {
		return iam.User{}, errors.New("user not found")
	}
	if name != "" {
		existing.Name = name
	}
	if phone != "" {
		existing.Phone = normalizePhone(phone)
	}
	if password != "" {
		if len(password) < 6 {
			return iam.User{}, errors.New("password must be at least 6 characters")
		}
		hash, err := hashPassword(password)
		if err != nil {
			return iam.User{}, err
		}
		existing.PasswordHash = hash
	}
	if roles != nil {
		existing.Roles = uniqueStrings(roles)
	}
	if spaces != nil {
		existing.Spaces = uniqueStrings(spaces)
	}
	existing.UpdatedAt = time.Now().UTC()
	if err := s.Users.Save(existing); err != nil {
		return iam.User{}, err
	}
	return existing, nil
}

// DeleteUser removes a user.
func (s *IAM) DeleteUser(id string) error {
	return s.Users.Delete(id)
}

// ListUsers returns all users.
func (s *IAM) ListUsers() []iam.User {
	return s.Users.List()
}

// CreateRole registers a role definition.
func (s *IAM) CreateRole(name, description string, permissions []string) (iam.Role, error) {
	if strings.TrimSpace(name) == "" {
		return iam.Role{}, errors.New("name is required")
	}
	now := time.Now().UTC()
	role := iam.Role{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Permissions: uniqueStrings(permissions),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.Roles.Save(role); err != nil {
		return iam.Role{}, err
	}
	return role, nil
}

// ListRoles returns all roles.
func (s *IAM) ListRoles() []iam.Role {
	return s.Roles.List()
}

// CreateSpace registers a new space.
func (s *IAM) CreateSpace(name, description string) (iam.Space, error) {
	if strings.TrimSpace(name) == "" {
		return iam.Space{}, errors.New("name is required")
	}
	now := time.Now().UTC()
	space := iam.Space{
		ID:          generateID(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.Spaces.Save(space); err != nil {
		return iam.Space{}, err
	}
	return space, nil
}

// ListSpaces returns all spaces.
func (s *IAM) ListSpaces() []iam.Space {
	return s.Spaces.List()
}

// AssignPolicy stores a policy binding.
func (s *IAM) AssignPolicy(spaceID, roleID, resource, action string) (iam.Policy, error) {
	if strings.TrimSpace(spaceID) == "" {
		return iam.Policy{}, errors.New("space id required")
	}
	if strings.TrimSpace(roleID) == "" {
		return iam.Policy{}, errors.New("role id required")
	}
	now := time.Now().UTC()
	policy := iam.Policy{
		ID:        generateID(),
		SpaceID:   spaceID,
		RoleID:    roleID,
		Resource:  resource,
		Action:    action,
		CreatedAt: now,
	}
	if err := s.Policies.Save(policy); err != nil {
		return iam.Policy{}, err
	}
	return policy, nil
}

// ListPoliciesBySpace returns policies for a space.
func (s *IAM) ListPoliciesBySpace(spaceID string) []iam.Policy {
	return s.Policies.ListBySpace(spaceID)
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	if strings.HasPrefix(phone, "+") {
		return phone
	}
	if strings.HasPrefix(phone, "00") {
		return "+" + phone[2:]
	}
	return phone
}

// Authenticate verifies phone/password combination.
func (s *IAM) Authenticate(phone, password string) (iam.User, error) {
	if strings.TrimSpace(phone) == "" || strings.TrimSpace(password) == "" {
		return iam.User{}, errors.New("phone and password required")
	}
	phone = normalizePhone(phone)
	user, ok := s.Users.GetByPhone(phone)
	if !ok {
		return iam.User{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return iam.User{}, errors.New("invalid credentials")
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
