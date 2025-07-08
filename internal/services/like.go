package services

import (
	"digeon-backend/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

func NewLikeService(db *gorm.DB, notificationService *NotificationService) *LikeService {
	return &LikeService{
		db:                  db,
		notificationService: notificationService,
	}
}

func (s *LikeService) LikePost(userID, postID uuid.UUID) error {
	// Check if post exists
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return fmt.Errorf("failed to find post: %w", err)
	}

	// Check if already liked
	var existingLike models.Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error; err == nil {
		return errors.New("post already liked")
	}

	// Create like
	like := models.Like{
		UserID: userID,
		PostID: postID,
	}

	if err := s.db.Create(&like).Error; err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	// Update post likes count
	if err := s.db.Model(&post).Update("likes_count", gorm.Expr("likes_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to update likes count: %w", err)
	}

	// Create notification
	if s.notificationService != nil {
		s.notificationService.CreateLikeNotification(userID, postID)
	}

	return nil
}

func (s *LikeService) UnlikePost(userID, postID uuid.UUID) error {
	// Check if post exists
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return fmt.Errorf("failed to find post: %w", err)
	}

	// Find and delete like
	var like models.Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("like not found")
		}
		return fmt.Errorf("failed to find like: %w", err)
	}

	// Delete like
	if err := s.db.Delete(&like).Error; err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	// Update post likes count
	if err := s.db.Model(&post).Update("likes_count", gorm.Expr("likes_count - 1")).Error; err != nil {
		return fmt.Errorf("failed to update likes count: %w", err)
	}

	return nil
}

func (s *LikeService) GetPostLikes(postID uuid.UUID, limit, offset int) ([]models.UserPublic, error) {
	var likes []models.Like
	if err := s.db.Where("post_id = ?", postID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&likes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch likes: %w", err)
	}

	var users []models.UserPublic
	for _, like := range likes {
		userPublic := models.UserPublic{
			ID:              like.User.ID,
			Username:        like.User.Username,
			DisplayName:     like.User.DisplayName,
			Bio:             like.User.Bio,
			ProfileImageURL: like.User.ProfileImageURL,
			CoverImageURL:   like.User.CoverImageURL,
			Location:        like.User.Location,
			Website:         like.User.Website,
			IsVerified:      like.User.IsVerified,
			CreatedAt:       like.User.CreatedAt,
		}
		users = append(users, userPublic)
	}

	return users, nil
}

func (s *LikeService) GetUserLikes(userID uuid.UUID, limit, offset int) ([]models.PostWithDetails, error) {
	var likes []models.Like
	if err := s.db.Where("user_id = ?", userID).
		Preload("Post").
		Preload("Post.Author").
		Preload("Post.Media").
		Preload("Post.OriginalPost").
		Preload("Post.OriginalPost.Author").
		Preload("Post.ParentPost").
		Preload("Post.ParentPost.Author").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&likes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user likes: %w", err)
	}

	var posts []models.PostWithDetails
	for _, like := range likes {
		postWithDetails := models.PostWithDetails{
			Post:     like.Post,
			IsLiked:  true, // Obviously true since we're getting user's liked posts
		}

		// Check if user reposted this post
		var repost models.Post
		if err := s.db.Where("author_id = ? AND original_post_id = ? AND type = ?", userID, like.Post.ID, models.PostTypeRepost).First(&repost).Error; err == nil {
			postWithDetails.IsReposted = true
		}

		posts = append(posts, postWithDetails)
	}

	return posts, nil
}

func (s *LikeService) IsPostLikedByUser(userID, postID uuid.UUID) (bool, error) {
	var like models.Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check like status: %w", err)
	}
	return true, nil
}