package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gideonzy/knowledge-base/internal/common/auth"
	"github.com/gideonzy/knowledge-base/internal/common/middleware"
	"github.com/gideonzy/knowledge-base/internal/iam"
	"github.com/gideonzy/knowledge-base/internal/iam/service"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	mux.HandleFunc("/api/auth/login", h.handleLogin)
	authMiddleware := middleware.RequireAuth(h.logger, h.tokens)
	mux.Handle("/api/users", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleUsers))))
	mux.Handle("/api/users/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleUserByID))))
	mux.Handle("/api/roles", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleRoles))))
	mux.Handle("/api/roles/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleRoleByID))))
	mux.Handle("/api/spaces", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleSpaces))))
	mux.Handle("/api/spaces/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleSpaceByID))))
	mux.Handle("/api/policies", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handlePolicies))))
	mux.Handle("/api/policies/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handlePoliciesBySpace))))
	return middleware.Logging(h.logger)(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleLogin godoc
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
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if h.tokens == nil {
		writeError(w, http.StatusServiceUnavailable, "token manager not configured")
		return
	}
	var payload LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	user, err := h.service.Authenticate(payload.Phone, payload.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.tokens.Generate(user.ID, user.Phone)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	writeJSON(w, http.StatusOK, LoginResponse{Token: token, User: user})
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.service.ListUsers()
		writeJSON(w, http.StatusOK, users)
	case http.MethodPost:
		var payload CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		user, err := h.service.CreateUser(payload.Name, payload.Phone, payload.Password, payload.Roles, payload.Spaces)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, user)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "user id required")
		return
	}
	switch r.Method {
	case http.MethodPatch:
		var payload UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		user, err := h.service.UpdateUser(id, payload.Name, payload.Phone, payload.Password, payload.Roles, payload.Spaces)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)
	case http.MethodDelete:
		if err := h.service.DeleteUser(id); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		roles := h.service.ListRoles()
		writeJSON(w, http.StatusOK, roles)
	case http.MethodPost:
		var payload CreateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		role, err := h.service.CreateRole(payload.Name, payload.Description, payload.Permissions)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, role)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleRoleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/roles/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "role id required")
		return
	}
	if err := h.service.Roles.Delete(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleSpaces(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		spaces := h.service.ListSpaces()
		writeJSON(w, http.StatusOK, spaces)
	case http.MethodPost:
		var payload CreateSpaceRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		space, err := h.service.CreateSpace(payload.Name, payload.Description)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, space)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleSpaceByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/spaces/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "space id required")
		return
	}
	if err := h.service.Spaces.Delete(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handlePolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var payload CreatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	policy, err := h.service.AssignPolicy(payload.SpaceID, payload.RoleID, payload.Resource, payload.Action)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, policy)
}

// docListUsers godoc
// @Summary List users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.User
// @Router /api/users [get]
func (h *Handler) docListUsers() {}

// docCreateUser godoc
// @Summary Create user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user"
// @Success 201 {object} iam.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users [post]
func (h *Handler) docCreateUser() {}

// docUpdateUser godoc
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
func (h *Handler) docUpdateUser() {}

// docDeleteUser godoc
// @Summary Delete user
// @Tags Users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Router /api/users/{id} [delete]
func (h *Handler) docDeleteUser() {}

// docListRoles godoc
// @Summary List roles
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Role
// @Router /api/roles [get]
func (h *Handler) docListRoles() {}

// docCreateRole godoc
// @Summary Create role
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "Create role"
// @Success 201 {object} iam.Role
// @Failure 400 {object} ErrorResponse
// @Router /api/roles [post]
func (h *Handler) docCreateRole() {}

// docDeleteRole godoc
// @Summary Delete role
// @Tags Roles
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 204 {string} string "No Content"
// @Router /api/roles/{id} [delete]
func (h *Handler) docDeleteRole() {}

// docListSpaces godoc
// @Summary List spaces
// @Tags Spaces
// @Security BearerAuth
// @Produce json
// @Success 200 {array} iam.Space
// @Router /api/spaces [get]
func (h *Handler) docListSpaces() {}

// docCreateSpace godoc
// @Summary Create space
// @Tags Spaces
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSpaceRequest true "Create space"
// @Success 201 {object} iam.Space
// @Failure 400 {object} ErrorResponse
// @Router /api/spaces [post]
func (h *Handler) docCreateSpace() {}

// docDeleteSpace godoc
// @Summary Delete space
// @Tags Spaces
// @Security BearerAuth
// @Param id path string true "Space ID"
// @Success 204 {string} string "No Content"
// @Router /api/spaces/{id} [delete]
func (h *Handler) docDeleteSpace() {}

// docCreatePolicy godoc
// @Summary Create policy
// @Tags Policies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreatePolicyRequest true "Create policy"
// @Success 201 {object} iam.Policy
// @Failure 400 {object} ErrorResponse
// @Router /api/policies [post]
func (h *Handler) docCreatePolicy() {}

// docListPolicies godoc
// @Summary List policies for a space
// @Tags Policies
// @Security BearerAuth
// @Param spaceId path string true "Space ID"
// @Success 200 {array} iam.Policy
// @Router /api/policies/{spaceId} [get]
func (h *Handler) docListPolicies() {}

func (h *Handler) handlePoliciesBySpace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	spaceID := strings.TrimPrefix(r.URL.Path, "/api/policies/")
	if spaceID == "" {
		writeError(w, http.StatusBadRequest, "space id required")
		return
	}
	policies := h.service.ListPoliciesBySpace(spaceID)
	writeJSON(w, http.StatusOK, policies)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
