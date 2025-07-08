package services

import (
	"digeon-backend/internal/models"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	db           *gorm.DB
	mediaService *MediaService
}

func NewPostService(db *gorm.DB, mediaService *MediaService) *PostService {
	return &PostService{
		db:           db,
		mediaService: mediaService,
	}
}

type CreatePostRequest struct {
	Content        string    `json:"content" validate:"max=280"`
	Type           string    `json:"type" validate:"required"`
	OriginalPostID *string   `json:"original_post_id,omitempty"`
	ParentPostID   *string   `json:"parent_post_id,omitempty"`
	MediaURLs      []string  `json:"media_urls,omitempty"`
	IsDraft        bool      `json:"is_draft"`
}

type UpdatePostRequest struct {
	Content string `json:"content" validate:"max=280"`
}

func (s *PostService) CreatePost(userID uuid.UUID, req CreatePostRequest) (*models.PostWithDetails, error) {
	// Validate input
	if err := s.validateCreatePostRequest(req); err != nil {
		return nil, err
	}

	// Create post
	post := models.Post{
		AuthorID: userID,
		Content:  req.Content,
		Type:     models.PostType(req.Type),
		IsDraft:  req.IsDraft,
		IsPublic: true,
	}

	// Handle original post reference (for reposts and quotes)
	if req.OriginalPostID != nil {
		originalID, err := uuid.Parse(*req.OriginalPostID)
		if err != nil {
			return nil, errors.New("invalid original post ID")
		}
		
		// Verify original post exists
		var originalPost models.Post
		if err := s.db.First(&originalPost, originalID).Error; err != nil {
			return nil, errors.New("original post not found")
		}
		
		post.OriginalPostID = &originalID
	}

	// Handle parent post reference (for replies)
	if req.ParentPostID != nil {
		parentID, err := uuid.Parse(*req.ParentPostID)
		if err != nil {
			return nil, errors.New("invalid parent post ID")
		}
		
		// Verify parent post exists
		var parentPost models.Post
		if err := s.db.First(&parentPost, parentID).Error; err != nil {
			return nil, errors.New("parent post not found")
		}
		
		post.ParentPostID = &parentID
	}

	// Save post
	if err := s.db.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Process hashtags
	if err := s.processHashtags(post.ID, req.Content); err != nil {
		return nil, fmt.Errorf("failed to process hashtags: %w", err)
	}

	// Attach media if provided
	if len(req.MediaURLs) > 0 && s.mediaService != nil {
		if err := s.mediaService.AttachMediaToPost(req.MediaURLs, post.ID); err != nil {
			return nil, fmt.Errorf("failed to attach media: %w", err)
		}
	}

	// Update parent post comment count if this is a reply
	if req.ParentPostID != nil {
		parentID, _ := uuid.Parse(*req.ParentPostID)
		s.db.Model(&models.Post{}).Where("id = ?", parentID).Update("comments_count", gorm.Expr("comments_count + 1"))
	}

	// Update original post repost count if this is a repost
	if req.OriginalPostID != nil && req.Type == string(models.PostTypeRepost) {
		originalID, _ := uuid.Parse(*req.OriginalPostID)
		s.db.Model(&models.Post{}).Where("id = ?", originalID).Update("reposts_count", gorm.Expr("reposts_count + 1"))
	}

	// Return post with details
	return s.GetPostWithDetails(post.ID, userID)
}

func (s *PostService) GetPostByID(postID uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := s.db.Preload("Author").Preload("Media").Preload("OriginalPost").Preload("ParentPost").First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetPostWithDetails(postID, userID uuid.UUID) (*models.PostWithDetails, error) {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	postWithDetails := &models.PostWithDetails{
		Post: *post,
	}

	// Check if user liked this post
	var like models.Like
	if err := s.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err == nil {
		postWithDetails.IsLiked = true
	}

	// Check if user reposted this post
	var repost models.Post
	if err := s.db.Where("author_id = ? AND original_post_id = ? AND type = ?", userID, postID, models.PostTypeRepost).First(&repost).Error; err == nil {
		postWithDetails.IsReposted = true
	}

	return postWithDetails, nil
}

func (s *PostService) GetPostsByUserID(userID uuid.UUID, limit, offset int) ([]models.PostWithDetails, error) {
	var posts []models.Post
	if err := s.db.Where("author_id = ? AND is_draft = false", userID).
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("ParentPost").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}

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

func (s *PostService) UpdatePost(postID, userID uuid.UUID, req UpdatePostRequest) error {
	// Verify post exists and belongs to user
	var post models.Post
	if err := s.db.Where("id = ? AND author_id = ?", postID, userID).First(&post).Error; err != nil {
		return errors.New("post not found or unauthorized")
	}

	// Validate content
	if len(req.Content) > 280 {
		return errors.New("content exceeds 280 characters")
	}

	// Update post
	if err := s.db.Model(&post).Update("content", req.Content).Error; err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	// Reprocess hashtags
	if err := s.processHashtags(postID, req.Content); err != nil {
		return fmt.Errorf("failed to process hashtags: %w", err)
	}

	return nil
}

func (s *PostService) DeletePost(postID, userID uuid.UUID) error {
	// Verify post exists and belongs to user
	var post models.Post
	if err := s.db.Where("id = ? AND author_id = ?", postID, userID).First(&post).Error; err != nil {
		return errors.New("post not found or unauthorized")
	}

	// Delete post (soft delete)
	if err := s.db.Delete(&post).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Update parent post comment count if this was a reply
	if post.ParentPostID != nil {
		s.db.Model(&models.Post{}).Where("id = ?", *post.ParentPostID).Update("comments_count", gorm.Expr("comments_count - 1"))
	}

	// Update original post repost count if this was a repost
	if post.OriginalPostID != nil && post.Type == models.PostTypeRepost {
		s.db.Model(&models.Post{}).Where("id = ?", *post.OriginalPostID).Update("reposts_count", gorm.Expr("reposts_count - 1"))
	}

	return nil
}

func (s *PostService) validateCreatePostRequest(req CreatePostRequest) error {
	// Validate content length
	if len(req.Content) > 280 {
		return errors.New("content exceeds 280 characters")
	}

	// Validate post type
	validTypes := []string{
		string(models.PostTypeOriginal),
		string(models.PostTypeRepost),
		string(models.PostTypeQuote),
		string(models.PostTypeReply),
	}
	
	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	
	if !isValidType {
		return errors.New("invalid post type")
	}

	// Validate repost/quote requirements
	if req.Type == string(models.PostTypeRepost) || req.Type == string(models.PostTypeQuote) {
		if req.OriginalPostID == nil {
			return errors.New("original post ID is required for reposts and quotes")
		}
	}

	// Validate reply requirements
	if req.Type == string(models.PostTypeReply) {
		if req.ParentPostID == nil {
			return errors.New("parent post ID is required for replies")
		}
	}

	return nil
}

func (s *PostService) processHashtags(postID uuid.UUID, content string) error {
	// Remove existing hashtag associations
	s.db.Exec("DELETE FROM post_hashtags WHERE post_id = ?", postID)

	// Extract hashtags from content
	hashtags := extractHashtags(content)

	for _, tag := range hashtags {
		var hashtag models.Hashtag
		
		// Find or create hashtag
		if err := s.db.Where("name = ?", tag).First(&hashtag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				hashtag = models.Hashtag{Name: tag}
				if err := s.db.Create(&hashtag).Error; err != nil {
					continue
				}
			} else {
				continue
			}
		}

		// Create association
		s.db.Exec("INSERT INTO post_hashtags (post_id, hashtag_id) VALUES (?, ?) ON CONFLICT DO NOTHING", postID, hashtag.ID)
	}

	return nil
}

func extractHashtags(content string) []string {
	re := regexp.MustCompile(`#(\w+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	var hashtags []string
	for _, match := range matches {
		if len(match) > 1 {
			hashtags = append(hashtags, strings.ToLower(match[1]))
		}
	}
	
	return hashtags
}