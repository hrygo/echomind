package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask creates a new task for a user.
func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, title, description string, sourceEmailID *uuid.UUID, dueDate *time.Time) (*model.Task, error) {
	task := &model.Task{
		ID:            uuid.New(),
		UserID:        userID,
		SourceEmailID: sourceEmailID,
		Title:         title,
		Description:   description,
		Status:        model.TaskStatusTodo,
		Priority:      model.TaskPriorityMedium, // Default priority
		DueDate:       dueDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// ListTasks lists tasks for a user, with optional status and priority filters.
func (s *TaskService) ListTasks(ctx context.Context, userID uuid.UUID, status, priority string, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task

	query := s.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if limit <= 0 {
		limit = 20
	} // Default limit
	if offset < 0 {
		offset = 0
	} // Default offset

	// Order by DueDate (ascending) and CreatedAt (descending) for consistency
	if err := query.Limit(limit).Offset(offset).Order("due_date ASC, created_at DESC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTaskStatus updates the status of a specific task.
func (s *TaskService) UpdateTaskStatus(ctx context.Context, userID, taskID uuid.UUID, newStatus model.TaskStatus) error {
	var task model.Task
	// Ensure the task belongs to the user
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task not found or unauthorized")
		}
		return fmt.Errorf("failed to find task: %w", err)
	}

	// Validate new status
	switch newStatus {
	case model.TaskStatusTodo, model.TaskStatusInProgress, model.TaskStatusDone:
		task.Status = newStatus
		task.UpdatedAt = time.Now()
	default:
		return fmt.Errorf("invalid task status: %s", newStatus)
	}

	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

// UpdateTask updates a task's fields (title, description, priority, dueDate)
func (s *TaskService) UpdateTask(ctx context.Context, userID, taskID uuid.UUID, updateData map[string]interface{}) error {
	var task model.Task
	// Ensure the task belongs to the user
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task not found or unauthorized")
		}
		return fmt.Errorf("failed to find task: %w", err)
	}

	// Apply updates dynamically
	if title, ok := updateData["title"].(string); ok {
		task.Title = title
	}
	if description, ok := updateData["description"].(string); ok {
		task.Description = description
	}
	if priority, ok := updateData["priority"].(string); ok {
		switch model.TaskPriority(priority) {
		case model.TaskPriorityHigh, model.TaskPriorityMedium, model.TaskPriorityLow:
			task.Priority = model.TaskPriority(priority)
		default:
			return fmt.Errorf("invalid task priority: %s", priority)
		}
	}
	if dueDateStr, ok := updateData["due_date"].(string); ok && dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return fmt.Errorf("invalid due_date format: %w", err)
		}
		task.DueDate = &dueDate
	} else if dueDateStr == "" { // Allow clearing due date
		task.DueDate = nil
	}

	task.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// DeleteTask deletes a task.
func (s *TaskService) DeleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", taskID, userID).Delete(&model.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found or unauthorized")
	}
	return nil
}
