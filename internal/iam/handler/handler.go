package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gideonzy/knowledge-base/internal/common/auth"
	"github.com/gideonzy/knowledge-base/internal/common/middleware"
	"github.com/gideonzy/knowledge-base/internal/iam/service"
)

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
	mux.Handle("/api/users", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleUsers))))
	mux.Handle("/api/users/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleUserByID))))
	mux.Handle("/api/roles", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleRoles))))
	mux.Handle("/api/roles/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleRoleByID))))
	mux.Handle("/api/spaces", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleSpaces))))
	mux.Handle("/api/spaces/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleSpaceByID))))
	mux.Handle("/api/policies", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handlePolicies))))
	mux.Handle("/api/policies/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handlePoliciesBySpace))))
	return middleware.Logging(h.logger)(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if h.tokens == nil {
		writeError(w, http.StatusServiceUnavailable, "token manager not configured")
		return
	}
	var payload struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
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
	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id":     user.ID,
			"name":   user.Name,
			"phone":  user.Phone,
			"roles":  user.Roles,
			"spaces": user.Spaces,
		},
	})
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.service.ListUsers()
		writeJSON(w, http.StatusOK, users)
	case http.MethodPost:
		var payload struct {
			Name     string   `json:"name"`
			Phone    string   `json:"phone"`
			Password string   `json:"password"`
			Roles    []string `json:"roles"`
			Spaces   []string `json:"spaces"`
		}
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
		var payload struct {
			Name     string   `json:"name"`
			Phone    string   `json:"phone"`
			Password string   `json:"password"`
			Roles    []string `json:"roles"`
			Spaces   []string `json:"spaces"`
		}
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
		var payload struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Permissions []string `json:"permissions"`
		}
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
		var payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
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
	var payload struct {
		SpaceID  string `json:"space_id"`
		RoleID   string `json:"role_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
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
	writeJSON(w, status, map[string]string{"error": message})
}
