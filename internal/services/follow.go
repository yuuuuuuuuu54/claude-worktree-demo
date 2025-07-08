package services

import (
	"digeon-backend/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

func NewFollowService(db *gorm.DB, notificationService *NotificationService) *FollowService {
	return &FollowService{
		db:                  db,
		notificationService: notificationService,
	}
}

func (s *FollowService) FollowUser(followerID, followingID uuid.UUID) error {
	// Check if trying to follow themselves
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Check if following user exists
	var followingUser models.User
	if err := s.db.First(&followingUser, followingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Check if already following
	var existingFollow models.Follow
	if err := s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollow).Error; err == nil {
		return errors.New("already following user")
	}

	// Create follow relationship
	follow := models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := s.db.Create(&follow).Error; err != nil {
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	// Create notification
	if s.notificationService != nil {
		s.notificationService.CreateFollowNotification(followerID, followingID)
	}

	return nil
}

func (s *FollowService) UnfollowUser(followerID, followingID uuid.UUID) error {
	// Check if trying to unfollow themselves
	if followerID == followingID {
		return errors.New("cannot unfollow yourself")
	}

	// Find and delete follow relationship
	var follow models.Follow
	if err := s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("follow relationship not found")
		}
		return fmt.Errorf("failed to find follow relationship: %w", err)
	}

	// Delete follow relationship
	if err := s.db.Delete(&follow).Error; err != nil {
		return fmt.Errorf("failed to delete follow relationship: %w", err)
	}

	return nil
}

func (s *FollowService) GetFollowers(userID uuid.UUID, limit, offset int) ([]models.UserPublic, error) {
	var follows []models.Follow
	if err := s.db.Where("following_id = ?", userID).
		Preload("Follower").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}

	var users []models.UserPublic
	for _, follow := range follows {
		userPublic := models.UserPublic{
			ID:              follow.Follower.ID,
			Username:        follow.Follower.Username,
			DisplayName:     follow.Follower.DisplayName,
			Bio:             follow.Follower.Bio,
			ProfileImageURL: follow.Follower.ProfileImageURL,
			CoverImageURL:   follow.Follower.CoverImageURL,
			Location:        follow.Follower.Location,
			Website:         follow.Follower.Website,
			IsVerified:      follow.Follower.IsVerified,
			CreatedAt:       follow.Follower.CreatedAt,
		}
		users = append(users, userPublic)
	}

	return users, nil
}

func (s *FollowService) GetFollowing(userID uuid.UUID, limit, offset int) ([]models.UserPublic, error) {
	var follows []models.Follow
	if err := s.db.Where("follower_id = ?", userID).
		Preload("Following").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch following: %w", err)
	}

	var users []models.UserPublic
	for _, follow := range follows {
		userPublic := models.UserPublic{
			ID:              follow.Following.ID,
			Username:        follow.Following.Username,
			DisplayName:     follow.Following.DisplayName,
			Bio:             follow.Following.Bio,
			ProfileImageURL: follow.Following.ProfileImageURL,
			CoverImageURL:   follow.Following.CoverImageURL,
			Location:        follow.Following.Location,
			Website:         follow.Following.Website,
			IsVerified:      follow.Following.IsVerified,
			CreatedAt:       follow.Following.CreatedAt,
		}
		users = append(users, userPublic)
	}

	return users, nil
}

func (s *FollowService) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var follow models.Follow
	if err := s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return true, nil
}

func (s *FollowService) GetFollowCounts(userID uuid.UUID) (followersCount, followingCount int64, err error) {
	// Get followers count
	if err := s.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&followersCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	// Get following count
	if err := s.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&followingCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count following: %w", err)
	}

	return followersCount, followingCount, nil
}

func (s *FollowService) GetSuggestedUsers(userID uuid.UUID, limit int) ([]models.UserPublic, error) {
	// Get users that the current user is not following
	// and exclude the current user
	var users []models.User
	if err := s.db.Where(`
		id != ? AND
		id NOT IN (
			SELECT following_id FROM follows 
			WHERE follower_id = ? AND deleted_at IS NULL
		)
	`, userID, userID).
		Where("is_active = true").
		Order("created_at DESC").
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch suggested users: %w", err)
	}

	var userPublics []models.UserPublic
	for _, user := range users {
		userPublic := models.UserPublic{
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
		userPublics = append(userPublics, userPublic)
	}

	return userPublics, nil
}