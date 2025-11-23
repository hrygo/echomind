package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Automatically log in the user after successful registration to get a token and account status
	token, _, hasAccount, err := h.userService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Log this error, but registration was successful
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registered, but failed to generate login token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "User registered successfully",
		"token":       token,
		"user": gin.H{
			"id": user.ID,
			"email": user.Email,
			"name": user.Name,
			"role": user.Role,
			"has_account": hasAccount,
		},
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login and returns a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, hasAccount, err := h.userService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user": gin.H{
			"id": user.ID,
			"email": user.Email,
			"name": user.Name,
			"role": user.Role,
			"has_account": hasAccount,
		},
	})
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=executive manager dealmaker"`
}

// UpdateUserRole handles updating the authenticated user's role.
func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUserRole(c.Request.Context(), userID, req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}
