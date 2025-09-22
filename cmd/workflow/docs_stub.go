package main

import (
	"github.com/gideonzy/knowledge-base/internal/workflow"
	workflowhandler "github.com/gideonzy/knowledge-base/internal/workflow/handler"
	"github.com/gin-gonic/gin"
)

// ListDefinitionsDoc godoc
// @Summary List workflow definitions
// @Tags Definitions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} workflow.FlowDefinition
// @Router /api/flows [get]
func ListDefinitionsDoc(_ *gin.Context) {}

// CreateDefinitionDoc godoc
// @Summary Register workflow definition
// @Tags Definitions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body workflowhandler.RegisterFlowRequest true "Definition payload"
// @Success 201 {object} workflow.FlowDefinition
// @Failure 400 {object} workflowhandler.WFErrorResponse
// @Router /api/flows [post]
func CreateDefinitionDoc(_ *gin.Context) {}

// StartInstanceDoc godoc
// @Summary Start workflow instance
// @Tags Instances
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "Definition code"
// @Param request body workflowhandler.StartInstanceRequest true "Instance payload"
// @Success 201 {object} workflow.FlowInstance
// @Failure 400 {object} workflowhandler.WFErrorResponse
// @Failure 404 {object} workflowhandler.WFErrorResponse
// @Router /api/flows/{code}/instances [post]
func StartInstanceDoc(_ *gin.Context) {}

// GetInstanceDoc godoc
// @Summary Get workflow instance
// @Tags Instances
// @Security BearerAuth
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} workflow.FlowInstance
// @Failure 404 {object} workflowhandler.WFErrorResponse
// @Router /api/instances/{id} [get]
func GetInstanceDoc(_ *gin.Context) {}

// ActOnInstanceDoc godoc
// @Summary Submit action on workflow instance
// @Tags Instances
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Instance ID"
// @Param request body workflowhandler.InstanceActionRequest true "Action payload"
// @Success 200 {object} workflow.FlowInstance
// @Failure 400 {object} workflowhandler.WFErrorResponse
// @Failure 404 {object} workflowhandler.WFErrorResponse
// @Router /api/instances/{id}/actions [post]
func ActOnInstanceDoc(_ *gin.Context) {}
