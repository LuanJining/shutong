package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gideonzy/knowledge-base/internal/common/middleware"
	"github.com/gideonzy/knowledge-base/internal/workflow"
	"github.com/gideonzy/knowledge-base/internal/workflow/service"
)

// Handler exposes workflow endpoints.
type Handler struct {
	service *service.Workflow
	logger  *slog.Logger
}

// New constructs a handler.
func New(svc *service.Workflow, logger *slog.Logger) *Handler {
	return &Handler{service: svc, logger: logger}
}

// Routes registers HTTP routes.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	mux.Handle("/api/flows", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleDefinitions))))
	mux.Handle("/api/flows/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleDefinitionInstances))))
	mux.Handle("/api/instances/", middleware.RequireAuth(h.logger)(middleware.RequestID(http.HandlerFunc(h.handleInstanceActions))))
	return middleware.Logging(h.logger)(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleDefinitions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		defs := h.service.ListDefinitions()
		writeJSON(w, http.StatusOK, defs)
	case http.MethodPost:
		var payload struct {
			Code        string              `json:"code"`
			Name        string              `json:"name"`
			Description string              `json:"description"`
			Nodes       []workflow.FlowNode `json:"nodes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		def, err := h.service.RegisterDefinition(payload.Code, payload.Name, payload.Description, payload.Nodes)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, def)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleDefinitionInstances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// Expect path /api/flows/{code}/instances
	suffix := strings.TrimPrefix(r.URL.Path, "/api/flows/")
	parts := strings.Split(suffix, "/")
	if len(parts) < 2 || parts[1] != "instances" {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}
	code := parts[0]
	var payload struct {
		BusinessID string `json:"business_id"`
		SpaceID    string `json:"space_id"`
		CreatedBy  string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	inst, err := h.service.StartInstance(code, payload.BusinessID, payload.SpaceID, payload.CreatedBy)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, inst)
}

func (h *Handler) handleInstanceActions(w http.ResponseWriter, r *http.Request) {
	suffix := strings.TrimPrefix(r.URL.Path, "/api/instances/")
	if suffix == "" {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}
	parts := strings.Split(suffix, "/")
	instanceID := parts[0]

	if len(parts) == 1 && r.Method == http.MethodGet {
		inst, err := h.service.GetInstance(instanceID)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, inst)
		return
	}

	if len(parts) == 2 && parts[1] == "actions" && r.Method == http.MethodPost {
		var payload struct {
			ActorID string              `json:"actor_id"`
			Comment string              `json:"comment"`
			Action  workflow.TaskAction `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		inst, err := h.service.ApplyAction(instanceID, payload.ActorID, payload.Comment, payload.Action)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, inst)
		return
	}

	writeError(w, http.StatusNotFound, "resource not found")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
