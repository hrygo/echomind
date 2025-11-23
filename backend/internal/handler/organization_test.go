package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.User{}, &model.Organization{}, &model.OrganizationMember{})
	return db
}

func TestCreateOrganization(t *testing.T) {
	db := setupTestDB()
	svc := service.NewOrganizationService(db)
	handler := NewOrganizationHandler(svc)

	// Create a user
	userID := uuid.New()
	user := &model.User{ID: userID, Email: "test@example.com", Name: "Test User"}
	db.Create(user)

	// Setup Router
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, userID)
		c.Next()
	})
	r.POST("/api/v1/orgs", handler.CreateOrganization)

	// Make Request
	payload := map[string]string{"name": "New Corp"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/orgs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Organization
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Corp", response.Name)
	assert.Equal(t, userID, response.OwnerID)
}

func TestListOrganizations(t *testing.T) {
	db := setupTestDB()
	svc := service.NewOrganizationService(db)
	handler := NewOrganizationHandler(svc)

	userID := uuid.New()
	user := &model.User{ID: userID, Email: "test@example.com"}
	db.Create(user)

	// Create existing org
	_, _ = svc.CreateOrganization(toContext(), "Org A", userID)
	_, _ = svc.CreateOrganization(toContext(), "Org B", userID)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, userID)
		c.Next()
	})
	r.GET("/api/v1/orgs", handler.ListOrganizations)

	req, _ := http.NewRequest("GET", "/api/v1/orgs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var orgs []model.Organization
	_ = json.Unmarshal(w.Body.Bytes(), &orgs)
	assert.Len(t, orgs, 2)
}

// Helper to get context
func toContext() context.Context {
	return context.Background()
}
