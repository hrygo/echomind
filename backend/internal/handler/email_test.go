package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockEmailService implements handler.EmailServicer
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) ListEmails(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Email, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Email), args.Error(1)
}

func (m *MockEmailService) GetEmail(ctx context.Context, userID, emailID uuid.UUID) (*model.Email, error) {
	args := m.Called(ctx, userID, emailID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Email), args.Error(1)
}

func (m *MockEmailService) CreateEmail(ctx context.Context, email *model.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockEmailService) UpdateEmail(ctx context.Context, email *model.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockEmailService) DeleteAllUserEmails(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestEmailHandler_ListEmails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/emails?limit=10&offset=0", nil)

		expectedEmails := []model.Email{
			{ID: uuid.New(), Subject: "Email 1", Sender: "a@b.com", Date: time.Now()},
			{ID: uuid.New(), Subject: "Email 2", Sender: "c@d.com", Date: time.Now().Add(-time.Hour)},
		}
		mockService.On("ListEmails", mock.Anything, userID, 10, 0).Return(expectedEmails, nil)

		h.ListEmails(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualEmails []model.Email
		err := json.Unmarshal(w.Body.Bytes(), &actualEmails)
		assert.NoError(t, err)
		assert.Len(t, actualEmails, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/emails", nil)

		h.ListEmails(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "ListEmails")
	})
}

func TestEmailHandler_GetEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		emailID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Params = []gin.Param{{Key: "id", Value: emailID.String()}}
		c.Request = httptest.NewRequest("GET", "/api/v1/emails/"+emailID.String(), nil)

		expectedEmail := &model.Email{ID: emailID, Subject: "Single Email", Sender: "test@example.com"}
		mockService.On("GetEmail", mock.Anything, userID, emailID).Return(expectedEmail, nil)

		h.GetEmail(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var actualEmail model.Email
		err := json.Unmarshal(w.Body.Bytes(), &actualEmail)
		assert.NoError(t, err)
		assert.Equal(t, emailID, actualEmail.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		emailID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Params = []gin.Param{{Key: "id", Value: emailID.String()}}
		c.Request = httptest.NewRequest("GET", "/api/v1/emails/"+emailID.String(), nil)

		mockService.On("GetEmail", mock.Anything, userID, emailID).Return(nil, gorm.ErrRecordNotFound)

		h.GetEmail(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Email not found or not accessible", response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestEmailHandler_DeleteAllEmails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("DELETE", "/api/v1/emails/all", nil)

		mockService.On("DeleteAllUserEmails", mock.Anything, userID).Return(nil)

		h.DeleteAllEmails(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "All emails deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("DELETE", "/api/v1/emails/all", nil)

		mockService.On("DeleteAllUserEmails", mock.Anything, userID).Return(errors.New("db error"))

		h.DeleteAllEmails(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Failed to delete all emails")
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockService := new(MockEmailService)
		h := handler.NewEmailHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/v1/emails/all", nil)
		// No userID set

		h.DeleteAllEmails(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "DeleteAllUserEmails")
	})
}
