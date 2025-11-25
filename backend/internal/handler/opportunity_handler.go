package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
)

type OpportunityHandler struct {
	opportunityService *service.OpportunityService
}

func NewOpportunityHandler(opportunityService *service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{
		opportunityService: opportunityService,
	}
}

// CreateOpportunity godoc
// @Summary Create a new opportunity
// @Description Creates a new business opportunity for the current user
// @Tags Opportunities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOpportunityRequest true "Opportunity creation request"
// @Success 201 {object} model.Opportunity
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opportunities [post]
func (h *OpportunityHandler) CreateOpportunity(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// For now, we'll use a default team and org ID
	// In a real implementation, these should come from the JWT or user context
	teamID := uuid.New() // TODO: Get from actual user context
	orgID := uuid.New()  // TODO: Get from actual user context

	var req CreateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oppType := model.OpportunityType(req.Type)
	if !isValidOpportunityType(oppType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid opportunity type"})
		return
	}

	opportunity, err := h.opportunityService.CreateOpportunity(
		c.Request.Context(),
		userID.String(),
		teamID.String(),
		orgID.String(),
		req.Title,
		req.Description,
		req.Company,
		req.Value,
		oppType,
		req.Confidence,
		req.SourceEmailID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create opportunity"})
		return
	}

	c.JSON(http.StatusCreated, opportunity)
}

// ListOpportunities godoc
// @Summary List opportunities
// @Description Retrieves a list of opportunities for the current user with optional filtering
// @Tags Opportunities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(new, active, won, lost, on_hold)
// @Param type query string false "Filter by type" Enums(buying, partnership, renewal, strategic)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {array} model.Opportunity
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opportunities [get]
func (h *OpportunityHandler) ListOpportunities(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var status *model.OpportunityStatus
	var oppType *model.OpportunityType

	if statusStr := c.Query("status"); statusStr != "" {
		s := model.OpportunityStatus(statusStr)
		if isValidOpportunityStatus(s) {
			status = &s
		}
	}

	if typeStr := c.Query("type"); typeStr != "" {
		t := model.OpportunityType(typeStr)
		if isValidOpportunityType(t) {
			oppType = &t
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	opportunities, err := h.opportunityService.ListOpportunities(
		c.Request.Context(),
		userID.String(),
		status,
		oppType,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list opportunities"})
		return
	}

	c.JSON(http.StatusOK, opportunities)
}

// GetOpportunity godoc
// @Summary Get an opportunity
// @Description Retrieves a specific opportunity by ID
// @Tags Opportunities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Opportunity ID"
// @Success 200 {object} model.Opportunity
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opportunities/{id} [get]
func (h *OpportunityHandler) GetOpportunity(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Opportunity ID is required"})
		return
	}

	opportunity, err := h.opportunityService.GetOpportunityByID(c.Request.Context(), userID.String(), id)
	if err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Opportunity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get opportunity"})
		return
	}

	c.JSON(http.StatusOK, opportunity)
}

// UpdateOpportunity godoc
// @Summary Update an opportunity
// @Description Updates an existing opportunity
// @Tags Opportunities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Opportunity ID"
// @Param request body UpdateOpportunityRequest true "Opportunity update request"
// @Success 200 {object} model.Opportunity
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opportunities/{id} [patch]
func (h *OpportunityHandler) UpdateOpportunity(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Opportunity ID is required"})
		return
	}

	var req UpdateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Value != nil {
		updates["value"] = *req.Value
	}
	if req.Status != nil {
		status := model.OpportunityStatus(*req.Status)
		if !isValidOpportunityStatus(status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid opportunity status"})
			return
		}
		updates["status"] = status
	}
	if req.Confidence != nil {
		updates["confidence"] = *req.Confidence
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	opportunity, err := h.opportunityService.UpdateOpportunity(c.Request.Context(), userID.String(), id, updates)
	if err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Opportunity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update opportunity"})
		return
	}

	c.JSON(http.StatusOK, opportunity)
}

// DeleteOpportunity godoc
// @Summary Delete an opportunity
// @Description Deletes an opportunity
// @Tags Opportunities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Opportunity ID"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opportunities/{id} [delete]
func (h *OpportunityHandler) DeleteOpportunity(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Opportunity ID is required"})
		return
	}

	err := h.opportunityService.DeleteOpportunity(c.Request.Context(), userID.String(), id)
	if err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Opportunity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete opportunity"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Helper functions

func isValidOpportunityType(t model.OpportunityType) bool {
	switch t {
	case model.OpportunityTypeBuying, model.OpportunityTypePartnership, model.OpportunityTypeRenewal, model.OpportunityTypeStrategic:
		return true
	default:
		return false
	}
}

func isValidOpportunityStatus(s model.OpportunityStatus) bool {
	switch s {
	case model.OpportunityStatusNew, model.OpportunityStatusActive, model.OpportunityStatusWon, model.OpportunityStatusLost, model.OpportunityStatusOnHold:
		return true
	default:
		return false
	}
}

// Request types

type CreateOpportunityRequest struct {
	Title         string  `json:"title" binding:"required"`
	Description   string  `json:"description"`
	Company       string  `json:"company" binding:"required"`
	Value         string  `json:"value"`
	Type          string  `json:"type"`
	Confidence    int     `json:"confidence"`
	SourceEmailID *string `json:"source_email_id"`
}

type UpdateOpportunityRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Value       *string `json:"value"`
	Status      *string `json:"status"`
	Confidence  *int    `json:"confidence"`
}
