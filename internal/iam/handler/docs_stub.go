package handler

import (
	"github.com/gideonzy/knowledge-base/internal/iam"
	"github.com/gin-gonic/gin"
)

// LoginDoc godoc
// @Summary Obtain JWT token
// @Description Authenticate with phone and password to receive JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func LoginDoc(_ *gin.Context) {}

// ListUsersDoc godoc
// @Summary List users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.User
// @Router /api/users [get]
func ListUsersDoc(_ *gin.Context) {}

// CreateUserDoc godoc
// @Summary Create user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user"
// @Success 201 {object} iam.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users [post]
func CreateUserDoc(_ *gin.Context) {}

// UpdateUserDoc godoc
// @Summary Update user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update payload"
// @Success 200 {object} iam.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [patch]
func UpdateUserDoc(_ *gin.Context) {}

// DeleteUserDoc godoc
// @Summary Delete user
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Router /api/users/{id} [delete]
func DeleteUserDoc(_ *gin.Context) {}

// ListRolesDoc godoc
// @Summary List roles
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Role
// @Router /api/roles [get]
func ListRolesDoc(_ *gin.Context) {}

// CreateRoleDoc godoc
// @Summary Create role
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "Create role"
// @Success 201 {object} iam.Role
// @Failure 400 {object} ErrorResponse
// @Router /api/roles [post]
func CreateRoleDoc(_ *gin.Context) {}

// DeleteRoleDoc godoc
// @Summary Delete role
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 204 {string} string "No Content"
// @Router /api/roles/{id} [delete]
func DeleteRoleDoc(_ *gin.Context) {}

// ListSpacesDoc godoc
// @Summary List spaces
// @Tags Spaces
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Space
// @Router /api/spaces [get]
func ListSpacesDoc(_ *gin.Context) {}

// CreateSpaceDoc godoc
// @Summary Create space
// @Tags Spaces
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSpaceRequest true "Create space"
// @Success 201 {object} iam.Space
// @Failure 400 {object} ErrorResponse
// @Router /api/spaces [post]
func CreateSpaceDoc(_ *gin.Context) {}

// DeleteSpaceDoc godoc
// @Summary Delete space
// @Tags Spaces
// @Security BearerAuth
// @Param id path string true "Space ID"
// @Success 204 {string} string "No Content"
// @Router /api/spaces/{id} [delete]
func DeleteSpaceDoc(_ *gin.Context) {}

// CreatePolicyDoc godoc
// @Summary Create policy
// @Tags Policies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreatePolicyRequest true "Create policy"
// @Success 201 {object} iam.Policy
// @Failure 400 {object} ErrorResponse
// @Router /api/policies [post]
func CreatePolicyDoc(_ *gin.Context) {}

// ListPoliciesBySpaceDoc godoc
// @Summary List policies for a space
// @Tags Policies
// @Security BearerAuth
// @Param spaceId path string true "Space ID"
// @Success 200 {array} iam.Policy
// @Router /api/policies/{spaceId} [get]
func ListPoliciesBySpaceDoc(_ *gin.Context) {}
