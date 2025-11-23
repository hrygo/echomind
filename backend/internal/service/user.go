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
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService handles user-related business logic.
type UserService struct {
	db         *gorm.DB
	jwtCfg     configs.JWTConfig
	orgService *OrganizationService
}

// NewUserService creates a new UserService.
func NewUserService(db *gorm.DB, jwtCfg configs.JWTConfig, orgService *OrganizationService) *UserService {
	return &UserService{
		db:         db,
		jwtCfg:     jwtCfg,
		orgService: orgService,
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

	// Use Transaction to ensure User and Org are created together
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Create default organization
		// Note: We use the *service method* but passing the transaction context would be ideal.
		// However, s.orgService.CreatePersonalOrganization starts its own transaction.
		// Nested transactions are supported by GORM (SavePoints).
		// Nested transactions are supported by GORM (SavePoints).
		_, err := s.orgService.CreatePersonalOrganization(ctx, user, tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user and generates a JWT token.
func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, *model.User, bool, error) {
	var user model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, false, ErrInvalidCredentials
		}
		return "", nil, false, err
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, false, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, s.jwtCfg.Secret, s.jwtCfg.ExpirationHours)
	if err != nil {
		return "", nil, false, err
	}

	// Check if user has a connected email account
	var hasAccount bool
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.EmailAccount{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
		// Log error but don't fail login
		_ = err
	}
	hasAccount = count > 0

	return token, &user, hasAccount, nil
}

// UpdateUserRole updates the role of a user.
func (s *UserService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) error {
	return s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("role", role).Error
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
