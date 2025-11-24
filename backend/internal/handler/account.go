package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"gorm.io/gorm"
)

// AccountHandler handles requests related to user email accounts.
type AccountHandler struct {
	accountService *service.AccountService
}

// NewAccountHandler creates a new AccountHandler.
func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// ConnectAndSaveAccount handles the POST request to connect or update an email account.
func (h *AccountHandler) ConnectAndSaveAccount(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var input model.EmailAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.ConnectAndSaveAccount(c.Request.Context(), userID, &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account connected successfully", "account_id": account.ID})
}

// GetAccountStatus handles the GET request to retrieve the status of a user's email account.
func (h *AccountHandler) GetAccountStatus(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	account, err := h.accountService.GetAccountByUserID(c.Request.Context(), userID)
	if err != nil {
		// If account not found, return has_account: false
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"has_account": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a simplified view of the account status (NO PASSWORD)
	response := gin.H{
		"has_account":    true,
		"email":          account.Email,
		"server_address": account.ServerAddress,
		"server_port":    account.ServerPort,
		"username":       account.Username,
		"is_connected":   account.IsConnected,
		"last_sync_at":   account.LastSyncAt,
		"error_message":  account.ErrorMessage,
	}
	c.JSON(http.StatusOK, response)
}

// DisconnectAccount handles the DELETE request to remove a user's email account.
func (h *AccountHandler) DisconnectAccount(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := h.accountService.DisconnectAccount(c.Request.Context(), userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No email account found to disconnect"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account disconnected successfully"})
}
