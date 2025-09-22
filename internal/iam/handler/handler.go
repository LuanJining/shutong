package handler

import (
	"net/http"

	"log/slog"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/auth"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/middleware"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// LoginRequest represents the payload for login.
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginResponse wraps login success output.
type LoginResponse struct {
	Token string   `json:"token"`
	User  iam.User `json:"user"`
}

// CreateUserRequest captures create user payload.
type CreateUserRequest struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Spaces   []string `json:"spaces"`
}

// UpdateUserRequest captures update user payload.
type UpdateUserRequest struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Spaces   []string `json:"spaces"`
}

// CreateRoleRequest defines role creation payload.
type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// CreateSpaceRequest defines space creation payload.
type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreatePolicyRequest defines policy creation payload.
type CreatePolicyRequest struct {
	SpaceID  string `json:"space_id"`
	RoleID   string `json:"role_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// ErrorResponse represents error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Handler exposes IAM HTTP endpoints.
type Handler struct {
	service *service.IAM
	tokens  *auth.Manager
	logger  *slog.Logger
}

// New creates a new handler instance.
func New(svc *service.IAM, tokens *auth.Manager, logger *slog.Logger) *Handler {
	return &Handler{service: svc, tokens: tokens, logger: logger}
}

// Routes returns a http.Handler with IAM routes mounted.
func (h *Handler) Routes() http.Handler {
	router := gin.New()
	router.Use(middleware.GinRequestID(), middleware.GinLogging(h.logger))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))


	router.GET("/healthz", h.HandleHealth)
	router.POST("/api/auth/login", h.Login)

	authGroup := router.Group("/api")
	authGroup.Use(middleware.GinRequireAuth(h.logger, h.tokens))

	authGroup.GET("/users", h.ListUsers)
	authGroup.POST("/users", h.CreateUser)
	authGroup.PATCH("/users/:id", h.UpdateUser)
	authGroup.DELETE("/users/:id", h.DeleteUser)

	authGroup.GET("/roles", h.ListRoles)
	authGroup.POST("/roles", h.CreateRole)
	authGroup.DELETE("/roles/:id", h.DeleteRole)

	authGroup.GET("/spaces", h.ListSpaces)
	authGroup.POST("/spaces", h.CreateSpace)
	authGroup.DELETE("/spaces/:id", h.DeleteSpace)

	authGroup.POST("/policies", h.CreatePolicy)
	authGroup.GET("/policies/:spaceId", h.ListPoliciesBySpace)

	return router
}

func (h *Handler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Login godoc
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
func (h *Handler) Login(c *gin.Context) {
	var payload LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	if h.tokens == nil {
		respondError(c, http.StatusServiceUnavailable, "token manager not configured")
		return
	}
	user, err := h.service.Authenticate(payload.Phone, payload.Password)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.tokens.Generate(user.ID, user.Phone)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token, User: user})
}

// ListUsers godoc
// @Summary List users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.User
// @Router /api/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users := h.service.ListUsers()
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user"
// @Success 201 {object} iam.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var payload CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	user, err := h.service.CreateUser(payload.Name, payload.Phone, payload.Password, payload.Roles, payload.Spaces)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
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
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var payload UpdateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	user, err := h.service.UpdateUser(id, payload.Name, payload.Phone, payload.Password, payload.Roles, payload.Spaces)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Router /api/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// ListRoles godoc
// @Summary List roles
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Role
// @Router /api/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	roles := h.service.ListRoles()
	c.JSON(http.StatusOK, roles)
}

// CreateRole godoc
// @Summary Create role
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "Create role"
// @Success 201 {object} iam.Role
// @Failure 400 {object} ErrorResponse
// @Router /api/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var payload CreateRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	role, err := h.service.CreateRole(payload.Name, payload.Description, payload.Permissions)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, role)
}

// DeleteRole godoc
// @Summary Delete role
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 204 {string} string "No Content"
// @Router /api/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Roles.Delete(id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// ListSpaces godoc
// @Summary List spaces
// @Tags Spaces
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Space
// @Router /api/spaces [get]
func (h *Handler) ListSpaces(c *gin.Context) {
	spaces := h.service.ListSpaces()
	c.JSON(http.StatusOK, spaces)
}

// CreateSpace godoc
// @Summary Create space
// @Tags Spaces
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSpaceRequest true "Create space"
// @Success 201 {object} iam.Space
// @Failure 400 {object} ErrorResponse
// @Router /api/spaces [post]
func (h *Handler) CreateSpace(c *gin.Context) {
	var payload CreateSpaceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	space, err := h.service.CreateSpace(payload.Name, payload.Description)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, space)
}

// DeleteSpace godoc
// @Summary Delete space
// @Tags Spaces
// @Security BearerAuth
// @Param id path string true "Space ID"
// @Success 204 {string} string "No Content"
// @Router /api/spaces/{id} [delete]
func (h *Handler) DeleteSpace(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Spaces.Delete(id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// CreatePolicy godoc
// @Summary Create policy
// @Tags Policies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreatePolicyRequest true "Create policy"
// @Success 201 {object} iam.Policy
// @Failure 400 {object} ErrorResponse
// @Router /api/policies [post]
func (h *Handler) CreatePolicy(c *gin.Context) {
	var payload CreatePolicyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid json")
		return
	}
	policy, err := h.service.AssignPolicy(payload.SpaceID, payload.RoleID, payload.Resource, payload.Action)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, policy)
}

// ListPoliciesBySpace godoc
// @Summary List policies for a space
// @Tags Policies
// @Security BearerAuth
// @Param spaceId path string true "Space ID"
// @Success 200 {array} iam.Policy
// @Router /api/policies/{spaceId} [get]
func (h *Handler) ListPoliciesBySpace(c *gin.Context) {
	spaceID := c.Param("spaceId")
	policies := h.service.ListPoliciesBySpace(spaceID)
	c.JSON(http.StatusOK, policies)
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Error: message})
}
