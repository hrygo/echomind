package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/utils"
)

var ( // Define custom errors
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

// UserService handles user-related business logic.
type UserService struct {
	db     *gorm.DB
	jwtCfg configs.JWTConfig
}

// NewUserService creates a new UserService.
func NewUserService(db *gorm.DB, jwtCfg configs.JWTConfig) *UserService {
	return &UserService{
		db:     db,
		jwtCfg: jwtCfg,
	}
}

// RegisterUser creates a new user in the database.
func (s *UserService) RegisterUser(ctx context.Context, email, password, name string) (*model.User, error) {
	// Check if user already exists
	var existingUser model.User
	if s.db.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error == nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and generates a JWT token.
func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, *model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, s.jwtCfg.Secret, s.jwtCfg.ExpirationHours)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}
