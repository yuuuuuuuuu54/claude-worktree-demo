package services

import (
	"digeon-backend/internal/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchService struct {
	db *gorm.DB
}

func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{db: db}
}

type SearchResponse struct {
	Users    []models.UserPublic       `json:"users,omitempty"`
	Posts    []models.PostWithDetails  `json:"posts,omitempty"`
	Hashtags []models.Hashtag          `json:"hashtags,omitempty"`
	Total    int64                     `json:"total"`
}

type SearchUsersResponse struct {
	Users  []models.UserPublic `json:"users"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
	Total  int64               `json:"total"`
}

type SearchPostsResponse struct {
	Posts  []models.PostWithDetails `json:"posts"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
	Total  int64                    `json:"total"`
}

type SearchHashtagsResponse struct {
	Hashtags []models.Hashtag `json:"hashtags"`
	Limit    int              `json:"limit"`
	Offset   int              `json:"offset"`
	Total    int64            `json:"total"`
}

// SearchAll searches across users, posts, and hashtags
func (s *SearchService) SearchAll(query string, userID uuid.UUID, limit int) (*SearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchResponse{}, nil
	}

	// Limit for each category (divide total limit by 3)
	categoryLimit := limit / 3
	if categoryLimit < 3 {
		categoryLimit = 3
	}

	response := &SearchResponse{}

	// Search users
	users, err := s.SearchUsers(query, categoryLimit, 0)
	if err == nil {
		response.Users = users.Users
	}

	// Search posts
	posts, err := s.SearchPosts(query, userID, categoryLimit, 0)
	if err == nil {
		response.Posts = posts.Posts
	}

	// Search hashtags
	hashtags, err := s.SearchHashtags(query, categoryLimit, 0)
	if err == nil {
		response.Hashtags = hashtags.Hashtags
	}

	response.Total = int64(len(response.Users) + len(response.Posts) + len(response.Hashtags))

	return response, nil
}

// SearchUsers searches for users by username or display name
func (s *SearchService) SearchUsers(query string, limit, offset int) (*SearchUsersResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchUsersResponse{Users: []models.UserPublic{}, Limit: limit, Offset: offset}, nil
	}

	var users []models.User
	var total int64

	// Use PostgreSQL full-text search and ILIKE for partial matches
	searchPattern := "%" + strings.ToLower(query) + "%"
	
	exactMatch := query
	prefixMatch := strings.ToLower(query) + "%"
	
	dbQuery := s.db.Where(`
		(LOWER(username) LIKE ? OR LOWER(display_name) LIKE ?) AND is_active = true
	`, searchPattern, searchPattern).
		Order(fmt.Sprintf("CASE WHEN username = '%s' THEN 1 WHEN LOWER(username) LIKE '%s' THEN 2 ELSE 3 END, created_at DESC", 
			exactMatch, prefixMatch))

	// Get total count
	if err := dbQuery.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	if err := dbQuery.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
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

	return &SearchUsersResponse{
		Users:  userPublics,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// SearchPosts searches for posts by content
func (s *SearchService) SearchPosts(query string, userID uuid.UUID, limit, offset int) (*SearchPostsResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchPostsResponse{Posts: []models.PostWithDetails{}, Limit: limit, Offset: offset}, nil
	}

	var posts []models.Post
	var total int64

	// Use PostgreSQL full-text search and ILIKE for partial matches
	searchPattern := "%" + strings.ToLower(query) + "%"
	
	dbQuery := s.db.Where(`
		LOWER(content) LIKE ? AND is_draft = false AND is_public = true
	`, searchPattern).
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("created_at DESC")

	// Get total count
	if err := dbQuery.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts with pagination
	if err := dbQuery.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	// Convert to PostWithDetails
	var postsWithDetails []models.PostWithDetails
	for _, post := range posts {
		postWithDetails := models.PostWithDetails{
			Post: post,
		}

		// Check if user liked this post
		if userID != uuid.Nil {
			var like models.Like
			if err := s.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error; err == nil {
				postWithDetails.IsLiked = true
			}

			// Check if user reposted this post
			var repost models.Post
			if err := s.db.Where("author_id = ? AND original_post_id = ? AND type = ?", userID, post.ID, models.PostTypeRepost).First(&repost).Error; err == nil {
				postWithDetails.IsReposted = true
			}
		}

		postsWithDetails = append(postsWithDetails, postWithDetails)
	}

	return &SearchPostsResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// SearchHashtags searches for hashtags
func (s *SearchService) SearchHashtags(query string, limit, offset int) (*SearchHashtagsResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchHashtagsResponse{Hashtags: []models.Hashtag{}, Limit: limit, Offset: offset}, nil
	}

	// Remove # if present
	if strings.HasPrefix(query, "#") {
		query = query[1:]
	}

	var hashtags []models.Hashtag
	var total int64

	searchPattern := "%" + strings.ToLower(query) + "%"
	
	exactMatch := query
	prefixMatch := strings.ToLower(query) + "%"
	
	dbQuery := s.db.Where("LOWER(name) LIKE ?", searchPattern).
		Order(fmt.Sprintf("CASE WHEN name = '%s' THEN 1 WHEN LOWER(name) LIKE '%s' THEN 2 ELSE 3 END, created_at DESC", 
			exactMatch, prefixMatch))

	// Get total count
	if err := dbQuery.Model(&models.Hashtag{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count hashtags: %w", err)
	}

	// Get hashtags with pagination
	if err := dbQuery.Limit(limit).Offset(offset).Find(&hashtags).Error; err != nil {
		return nil, fmt.Errorf("failed to search hashtags: %w", err)
	}

	return &SearchHashtagsResponse{
		Hashtags: hashtags,
		Limit:    limit,
		Offset:   offset,
		Total:    total,
	}, nil
}

// GetHashtagPosts gets posts for a specific hashtag
func (s *SearchService) GetHashtagPosts(hashtag string, userID uuid.UUID, limit, offset int) (*SearchPostsResponse, error) {
	hashtag = strings.TrimSpace(hashtag)
	if hashtag == "" {
		return &SearchPostsResponse{Posts: []models.PostWithDetails{}, Limit: limit, Offset: offset}, nil
	}

	// Remove # if present
	if strings.HasPrefix(hashtag, "#") {
		hashtag = hashtag[1:]
	}

	var posts []models.Post
	var total int64

	// Find hashtag first
	var hashtagRecord models.Hashtag
	if err := s.db.Where("LOWER(name) = LOWER(?)", hashtag).First(&hashtagRecord).Error; err != nil {
		return &SearchPostsResponse{Posts: []models.PostWithDetails{}, Limit: limit, Offset: offset}, nil
	}

	dbQuery := s.db.Joins("JOIN post_hashtags ON posts.id = post_hashtags.post_id").
		Where("post_hashtags.hashtag_id = ? AND posts.is_draft = false AND posts.is_public = true", hashtagRecord.ID).
		Preload("Author").
		Preload("Media").
		Preload("OriginalPost").
		Preload("OriginalPost.Author").
		Preload("ParentPost").
		Preload("ParentPost.Author").
		Order("posts.created_at DESC")

	// Get total count
	if err := dbQuery.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count hashtag posts: %w", err)
	}

	// Get posts with pagination
	if err := dbQuery.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get hashtag posts: %w", err)
	}

	// Convert to PostWithDetails
	var postsWithDetails []models.PostWithDetails
	for _, post := range posts {
		postWithDetails := models.PostWithDetails{
			Post: post,
		}

		// Check if user liked this post
		if userID != uuid.Nil {
			var like models.Like
			if err := s.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error; err == nil {
				postWithDetails.IsLiked = true
			}

			// Check if user reposted this post
			var repost models.Post
			if err := s.db.Where("author_id = ? AND original_post_id = ? AND type = ?", userID, post.ID, models.PostTypeRepost).First(&repost).Error; err == nil {
				postWithDetails.IsReposted = true
			}
		}

		postsWithDetails = append(postsWithDetails, postWithDetails)
	}

	return &SearchPostsResponse{
		Posts:  postsWithDetails,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

// GetTrendingHashtags gets the most popular hashtags
func (s *SearchService) GetTrendingHashtags(limit int) ([]models.Hashtag, error) {
	var hashtags []models.Hashtag

	// Get hashtags with most posts in the last 7 days
	if err := s.db.Select("hashtags.*, COUNT(post_hashtags.hashtag_id) as post_count").
		Joins("JOIN post_hashtags ON hashtags.id = post_hashtags.hashtag_id").
		Joins("JOIN posts ON post_hashtags.post_id = posts.id").
		Where("posts.created_at > NOW() - INTERVAL '7 days' AND posts.is_draft = false AND posts.is_public = true").
		Group("hashtags.id").
		Order("post_count DESC").
		Limit(limit).
		Find(&hashtags).Error; err != nil {
		return nil, fmt.Errorf("failed to get trending hashtags: %w", err)
	}

	return hashtags, nil
}