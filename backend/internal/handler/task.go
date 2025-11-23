package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Request structs for API input
type CreateTaskRequest struct {
	Title         string     `json:"title" binding:"required"`
	Description   string     `json:"description"`
	SourceEmailID *uuid.UUID `json:"source_email_id"`
	DueDate       *time.Time `json:"due_date" time_format:"2006-01-02T15:04:05Z07:00"` // RFC3339
}

type UpdateTaskStatusRequest struct {
	Status model.TaskStatus `json:"status" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	Priority    *model.TaskPriority `json:"priority"`
	DueDate     *time.Time          `json:"due_date" time_format:"2006-01-02T15:04:05Z07:00"`
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task data"
// @Success 201 {object} model.Task
// @Failure 400 {object} gin.H{"error":string}
// @Failure 401 {object} gin.H{"error":string}
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), userID, req.Title, req.Description, req.SourceEmailID, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ListTasks godoc
// @Summary List tasks
// @Description Get a list of tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Filter by task status (todo, in_progress, done)"
// @Param priority query string false "Filter by task priority (high, medium, low)"
// @Param limit query int false "Limit the number of results (default 20)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} model.Task
// @Failure 401 {object} gin.H{"error":string}
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	status := c.Query("status")
	priority := c.Query("priority")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	parsedLimit := 20
	if _, err := fmt.Sscanf(limit, "%d", &parsedLimit); err != nil {
		parsedLimit = 20
	}
	parsedOffset := 0
	if _, err := fmt.Sscanf(offset, "%d", &parsedOffset); err != nil {
		parsedOffset = 0
	}

	tasks, err := h.taskService.ListTasks(c.Request.Context(), userID, status, priority, parsedLimit, parsedOffset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update the status of a specific task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param status body UpdateTaskStatusRequest true "New status"
// @Success 200 {object} model.Task
// @Failure 400 {object} gin.H{"error":string}
// @Failure 401 {object} gin.H{"error":string}
// @Failure 404 {object} gin.H{"error":string}
// @Router /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.UpdateTaskStatus(c.Request.Context(), userID, taskID, req.Status); err != nil {
		if err.Error() == "task not found or unauthorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update fields of a specific task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body UpdateTaskRequest true "Task fields to update"
// @Success 200 {object} model.Task
// @Failure 400 {object} gin.H{"error":string}
// @Failure 401 {object} gin.H{"error":string}
// @Failure 404 {object} gin.H{"error":string}
// @Router /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := make(map[string]interface{})
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.Priority != nil {
		updateData["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		updateData["due_date"] = *req.DueDate
	}

	if err := h.taskService.UpdateTask(c.Request.Context(), userID, taskID, updateData); err != nil {
		if err.Error() == "task not found or unauthorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a specific task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H{"error":string}
// @Failure 401 {object} gin.H{"error":string}
// @Failure 404 {object} gin.H{"error":string}
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), userID, taskID); err != nil {
		if err.Error() == "task not found or unauthorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
