package services

import (
	"digeon-backend/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimelineService struct {
	db *gorm.DB
}

func NewTimelineService(db *gorm.DB) *TimelineService {
	return &TimelineService{db: db}
}

type TimelineResponse struct {
	Posts  []models.PostWithDetails `json:"posts"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
	Total  int64                    `json:"total"`
}

// GetHomeTimeline returns posts from users that the current user follows
func (s *TimelineService) GetHomeTimeline(userID uuid.UUID, limit, offset int) (*TimelineResponse, error) {
	var posts []models.Post
	var total int64

	// Get posts from followed users + user's own posts
	query := s.db.Where(`
		author_id IN (
			SELECT following_id FROM follows 
			WHERE follower_id = ? AND deleted_at IS NULL
		) OR author_id = ?
	`, userID, userID).
		Where("is_draft = false AND is_public = true").
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("created_at DESC")

	// Get total count
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts with pagination
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Convert to PostWithDetails
	postsWithDetails, err := s.convertToPostsWithDetails(posts, userID)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetExploreTimeline returns all public posts
func (s *TimelineService) GetExploreTimeline(userID uuid.UUID, limit, offset int) (*TimelineResponse, error) {
	var posts []models.Post
	var total int64

	query := s.db.Where("is_draft = false AND is_public = true").
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("created_at DESC")

	// Get total count
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts with pagination
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Convert to PostWithDetails
	postsWithDetails, err := s.convertToPostsWithDetails(posts, userID)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetTrendingTimeline returns posts sorted by engagement (likes + reposts + comments)
func (s *TimelineService) GetTrendingTimeline(userID uuid.UUID, limit, offset int) (*TimelineResponse, error) {
	var posts []models.Post
	var total int64

	query := s.db.Where("is_draft = false AND is_public = true").
		Where("created_at > NOW() - INTERVAL '7 days'"). // Only posts from last 7 days
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("(likes_count + reposts_count + comments_count) DESC, created_at DESC")

	// Get total count
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts with pagination
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Convert to PostWithDetails
	postsWithDetails, err := s.convertToPostsWithDetails(posts, userID)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetPostReplies returns replies to a specific post
func (s *TimelineService) GetPostReplies(postID, userID uuid.UUID, limit, offset int) (*TimelineResponse, error) {
	var posts []models.Post
	var total int64

	query := s.db.Where("parent_post_id = ? AND is_draft = false AND is_public = true", postID).
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("created_at ASC") // Oldest first for replies

	// Get total count
	if err := query.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts with pagination
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Convert to PostWithDetails
	postsWithDetails, err := s.convertToPostsWithDetails(posts, userID)
	if err != nil {
		return nil, err
	}

	return &TimelineResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// Helper function to convert posts to PostWithDetails
func (s *TimelineService) convertToPostsWithDetails(posts []models.Post, userID uuid.UUID) ([]models.PostWithDetails, error) {
	var result []models.PostWithDetails

	for _, post := range posts {
		postWithDetails := models.PostWithDetails{
			Post: post,
		}

		// Check if user liked this post
		var like models.Like
		if err := s.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error; err == nil {
			postWithDetails.IsLiked = true
		}

		// Check if user reposted this post
		var repost models.Post
		if err := s.db.Where("author_id = ? AND original_post_id = ? AND type = ?", userID, post.ID, models.PostTypeRepost).First(&repost).Error; err == nil {
			postWithDetails.IsReposted = true
		}

		result = append(result, postWithDetails)
	}

	return result, nil
}