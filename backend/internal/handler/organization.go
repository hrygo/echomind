package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

// CreateOrganization handles POST /api/v1/orgs
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet(middleware.ContextUserIDKey).(uuid.UUID)

	org, err := h.orgService.CreateOrganization(c.Request.Context(), req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// ListOrganizations handles GET /api/v1/orgs
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(uuid.UUID)

	orgs, err := h.orgService.GetUserOrganizations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch organizations"})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// GetOrganization handles GET /api/v1/orgs/:id
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	userID := c.MustGet(middleware.ContextUserIDKey).(uuid.UUID)

	org, err := h.orgService.GetOrganizationByID(c.Request.Context(), orgID, userID)
	if err != nil {
		// Distinguish between Not Found and Unauthorized/Error?
		// Service returns error if not found or not member
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// GetMembers handles GET /api/v1/orgs/:id/members
func (h *OrganizationHandler) GetMembers(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	// TODO: Verify requester has permission to view members (usually all members can view)
	// For now, we rely on the fact that if they can access the org, they can see members

	members, err := h.orgService.GetOrganizationMembers(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// InviteMember handles POST /api/v1/orgs/:id/invites
func (h *OrganizationHandler) InviteMember(c *gin.Context) {
	// TODO: Implement actual email invitation.
	// For now, we'll just add the user if they exist by email (Mock flow for testing)
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logic: Find user by email, if exists, add to org.
	// This belongs in Service layer properly, but putting here for rapid prototype
	// Need to expose UserService or move logic to OrganizationService

	c.JSON(http.StatusNotImplemented, gin.H{"error": "invitation not implemented yet"})
}
