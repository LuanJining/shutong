package handler

import (
	"net/http"

	"github.com/gideonzy/knowledge-base/internal/common/middleware"
	"github.com/gideonzy/knowledge-base/internal/workflow"
	"github.com/gideonzy/knowledge-base/internal/workflow/service"
	"github.com/gin-gonic/gin"
	"log/slog"
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
	router := gin.New()
	router.Use(middleware.GinRequestID(), middleware.GinLogging(h.logger))

	router.GET("/healthz", h.HandleHealth)

	protected := router.Group("/api")
	protected.Use(middleware.GinRequireAuth(h.logger, h.validator))

	protected.GET("/flows", h.ListDefinitions)
	protected.POST("/flows", h.CreateDefinition)
	protected.POST("/flows/:code/instances", h.StartInstance)
	protected.GET("/instances/:id", h.GetInstance)
	protected.POST("/instances/:id/actions", h.ActOnInstance)

	return router
}

func (h *Handler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListDefinitions godoc
// @Summary List workflow definitions
// @Tags Definitions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} workflow.FlowDefinition
// @Router /api/flows [get]
func (h *Handler) ListDefinitions(c *gin.Context) {
	defs := h.service.ListDefinitions()
	c.JSON(http.StatusOK, defs)
}

// CreateDefinition godoc
// @Summary Register workflow definition
// @Tags Definitions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body RegisterFlowRequest true "Definition payload"
// @Success 201 {object} workflow.FlowDefinition
// @Failure 400 {object} WFErrorResponse
// @Router /api/flows [post]
func (h *Handler) CreateDefinition(c *gin.Context) {
	var payload RegisterFlowRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWorkflowError(c, http.StatusBadRequest, "invalid json")
		return
	}
	def, err := h.service.RegisterDefinition(payload.Code, payload.Name, payload.Description, payload.Nodes)
	if err != nil {
		respondWorkflowError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, def)
}

// StartInstance godoc
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
func (h *Handler) StartInstance(c *gin.Context) {
	var payload StartInstanceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWorkflowError(c, http.StatusBadRequest, "invalid json")
		return
	}
	code := c.Param("code")
	inst, err := h.service.StartInstance(code, payload.BusinessID, payload.SpaceID, payload.CreatedBy)
	if err != nil {
		respondWorkflowError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, inst)
}

// GetInstance godoc
// @Summary Get workflow instance
// @Tags Instances
// @Security BearerAuth
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} workflow.FlowInstance
// @Failure 404 {object} WFErrorResponse
// @Router /api/instances/{id} [get]
func (h *Handler) GetInstance(c *gin.Context) {
	inst, err := h.service.GetInstance(c.Param("id"))
	if err != nil {
		respondWorkflowError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, inst)
}

// ActOnInstance godoc
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
func (h *Handler) ActOnInstance(c *gin.Context) {
	var payload InstanceActionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWorkflowError(c, http.StatusBadRequest, "invalid json")
		return
	}
	inst, err := h.service.ApplyAction(c.Param("id"), payload.ActorID, payload.Comment, payload.Action)
	if err != nil {
		respondWorkflowError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, inst)
}

func respondWorkflowError(c *gin.Context, status int, message string) {
	c.JSON(status, WFErrorResponse{Error: message})
}
