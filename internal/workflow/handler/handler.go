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

// RegisterFlowRequest captures definition registration payload.
type RegisterFlowRequest struct {
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Nodes       []workflow.FlowNode `json:"nodes"`
}

// StartInstanceRequest captures instance start payload.
type StartInstanceRequest struct {
	BusinessID string `json:"business_id"`
	SpaceID    string `json:"space_id"`
	CreatedBy  string `json:"created_by"`
}

// InstanceActionRequest captures workflow action payload.
type InstanceActionRequest struct {
	ActorID string              `json:"actor_id"`
	Comment string              `json:"comment"`
	Action  workflow.TaskAction `json:"action"`
}

// WFErrorResponse models workflow error output.
type WFErrorResponse struct {
	Error string `json:"error"`
}

// Handler exposes workflow endpoints.
type Handler struct {
	service   *service.Workflow
	logger    *slog.Logger
	validator middleware.TokenValidator
}

// New constructs a handler.
func New(svc *service.Workflow, validator middleware.TokenValidator, logger *slog.Logger) *Handler {
	return &Handler{service: svc, validator: validator, logger: logger}
}

// Routes registers HTTP routes.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	authMiddleware := middleware.RequireAuth(h.logger, h.validator)
	mux.Handle("/api/flows", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleDefinitions))))
	mux.Handle("/api/flows/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleDefinitionInstances))))
	mux.Handle("/api/instances/", authMiddleware(middleware.RequestID(http.HandlerFunc(h.handleInstanceActions))))
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
		var payload RegisterFlowRequest
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
	var payload StartInstanceRequest
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
		var payload InstanceActionRequest
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

// docListDefinitions godoc
// @Summary List workflow definitions
// @Tags Definitions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} workflow.FlowDefinition
// @Router /api/flows [get]
func (h *Handler) docListDefinitions() {}

// docCreateDefinition godoc
// @Summary Register workflow definition
// @Tags Definitions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body RegisterFlowRequest true "Definition payload"
// @Success 201 {object} workflow.FlowDefinition
// @Failure 400 {object} WFErrorResponse
// @Router /api/flows [post]
func (h *Handler) docCreateDefinition() {}

// docStartInstance godoc
// @Summary Start workflow instance
// @Tags Instances
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "Definition code"
// @Param request body StartInstanceRequest true "Instance payload"
// @Success 201 {object} workflow.FlowInstance
// @Failure 400 {object} WFErrorResponse
// @Failure 404 {object} WFErrorResponse
// @Router /api/flows/{code}/instances [post]
func (h *Handler) docStartInstance() {}

// docGetInstance godoc
// @Summary Get workflow instance
// @Tags Instances
// @Security BearerAuth
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} workflow.FlowInstance
// @Failure 404 {object} WFErrorResponse
// @Router /api/instances/{id} [get]
func (h *Handler) docGetInstance() {}

// docActOnInstance godoc
// @Summary Submit action on workflow instance
// @Tags Instances
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Instance ID"
// @Param request body InstanceActionRequest true "Action payload"
// @Success 200 {object} workflow.FlowInstance
// @Failure 400 {object} WFErrorResponse
// @Failure 404 {object} WFErrorResponse
// @Router /api/instances/{id}/actions [post]
func (h *Handler) docActOnInstance() {}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, WFErrorResponse{Error: message})
}
