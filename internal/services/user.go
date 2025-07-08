package services

import (
	"digeon-backend/internal/models"
	"digeon-backend/internal/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"max=100"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string             `json:"token"`
	User  models.UserPublic `json:"user"`
}

func (s *UserService) Register(req RegisterRequest) (*LoginResponse, error) {
	// Validate input
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		DisplayName:  req.DisplayName,
		IsActive:     true,
	}

	if user.DisplayName == "" {
		user.DisplayName = req.Username
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  s.toUserPublic(user),
	}, nil
}

func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	var user models.User
	
	// Find user by username or email
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  s.toUserPublic(user),
	}, nil
}

func (s *UserService) GetUserByID(userID uuid.UUID) (*models.UserPublic, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	userPublic := s.toUserPublic(user)
	
	// Get counts
	s.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&userPublic.FollowingCount)
	s.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&userPublic.FollowersCount)
	s.db.Model(&models.Post{}).Where("author_id = ?", userID).Count(&userPublic.PostsCount)

	return &userPublic, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.UserPublic, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	userPublic := s.toUserPublic(user)
	
	// Get counts
	s.db.Model(&models.Follow{}).Where("follower_id = ?", user.ID).Count(&userPublic.FollowingCount)
	s.db.Model(&models.Follow{}).Where("following_id = ?", user.ID).Count(&userPublic.FollowersCount)
	s.db.Model(&models.Post{}).Where("author_id = ?", user.ID).Count(&userPublic.PostsCount)

	return &userPublic, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, updates map[string]interface{}) error {
	allowedFields := []string{"display_name", "bio", "location", "website", "profile_image_url", "cover_image_url"}
	
	filteredUpdates := make(map[string]interface{})
	for _, field := range allowedFields {
		if value, exists := updates[field]; exists {
			filteredUpdates[field] = value
		}
	}

	if len(filteredUpdates) == 0 {
		return errors.New("no valid fields to update")
	}

	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(filteredUpdates).Error
}

func (s *UserService) validateRegisterRequest(req RegisterRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}

	if !strings.Contains(req.Email, "@") {
		return errors.New("invalid email format")
	}

	if !utils.IsValidPassword(req.Password) {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (s *UserService) toUserPublic(user models.User) models.UserPublic {
	return models.UserPublic{
		ID:              user.ID,
		Username:        user.Username,
		DisplayName:     user.DisplayName,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
		CoverImageURL:   user.CoverImageURL,
		Location:        user.Location,
		Website:         user.Website,
		IsVerified:      user.IsVerified,
		CreatedAt:       user.CreatedAt,
	}
}