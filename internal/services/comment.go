package services

import (
	"digeon-backend/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	db                  *gorm.DB
	postService         *PostService
	notificationService *NotificationService
}

func NewCommentService(db *gorm.DB, postService *PostService, notificationService *NotificationService) *CommentService {
	return &CommentService{
		db:                  db,
		postService:         postService,
		notificationService: notificationService,
	}
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=280"`
}

type CommentResponse struct {
	ID             uuid.UUID         `json:"id"`
	UserID         uuid.UUID         `json:"user_id"`
	PostID         uuid.UUID         `json:"post_id"`
	Content        string            `json:"content"`
	ParentID       *uuid.UUID        `json:"parent_id,omitempty"`
	LikesCount     int               `json:"likes_count"`
	RepliesCount   int               `json:"replies_count"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	User           models.UserPublic `json:"user"`
	IsLiked        bool              `json:"is_liked"`
}

func (s *CommentService) CreateComment(userID, postID uuid.UUID, req CreateCommentRequest) (*CommentResponse, error) {
	// Validate input
	if len(req.Content) == 0 {
		return nil, errors.New("comment content is required")
	}
	if len(req.Content) > 280 {
		return nil, errors.New("comment content exceeds 280 characters")
	}

	// Check if post exists
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	// Create comment as a reply post
	postIDStr := postID.String()
	createPostReq := CreatePostRequest{
		Content:      req.Content,
		Type:         string(models.PostTypeReply),
		ParentPostID: &postIDStr,
		IsDraft:      false,
	}

	postWithDetails, err := s.postService.CreatePost(userID, createPostReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Convert to comment response
	commentResponse := &CommentResponse{
		ID:           postWithDetails.ID,
		UserID:       postWithDetails.AuthorID,
		PostID:       postID,
		Content:      postWithDetails.Content,
		LikesCount:   postWithDetails.LikesCount,
		RepliesCount: postWithDetails.CommentsCount,
		CreatedAt:    postWithDetails.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    postWithDetails.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		User: models.UserPublic{
			ID:              postWithDetails.Author.ID,
			Username:        postWithDetails.Author.Username,
			DisplayName:     postWithDetails.Author.DisplayName,
			Bio:             postWithDetails.Author.Bio,
			ProfileImageURL: postWithDetails.Author.ProfileImageURL,
			CoverImageURL:   postWithDetails.Author.CoverImageURL,
			Location:        postWithDetails.Author.Location,
			Website:         postWithDetails.Author.Website,
			IsVerified:      postWithDetails.Author.IsVerified,
			CreatedAt:       postWithDetails.Author.CreatedAt,
		},
		IsLiked: postWithDetails.IsLiked,
	}

	if postWithDetails.ParentPostID != nil {
		commentResponse.ParentID = postWithDetails.ParentPostID
	}

	// Create notification
	if s.notificationService != nil {
		s.notificationService.CreateCommentNotification(userID, postID)
	}

	return commentResponse, nil
}

func (s *CommentService) CreateReply(userID, commentID uuid.UUID, req CreateCommentRequest) (*CommentResponse, error) {
	// Validate input
	if len(req.Content) == 0 {
		return nil, errors.New("reply content is required")
	}
	if len(req.Content) > 280 {
		return nil, errors.New("reply content exceeds 280 characters")
	}

	// Check if comment exists
	var comment models.Post
	if err := s.db.Where("id = ? AND type = ?", commentID, models.PostTypeReply).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	// Create reply as a nested comment
	commentIDStr := commentID.String()
	createPostReq := CreatePostRequest{
		Content:      req.Content,
		Type:         string(models.PostTypeReply),
		ParentPostID: &commentIDStr,
		IsDraft:      false,
	}

	postWithDetails, err := s.postService.CreatePost(userID, createPostReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create reply: %w", err)
	}

	// Convert to comment response
	commentResponse := &CommentResponse{
		ID:           postWithDetails.ID,
		UserID:       postWithDetails.AuthorID,
		PostID:       *comment.ParentPostID, // The original post ID
		Content:      postWithDetails.Content,
		ParentID:     &commentID,             // The comment being replied to
		LikesCount:   postWithDetails.LikesCount,
		RepliesCount: postWithDetails.CommentsCount,
		CreatedAt:    postWithDetails.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    postWithDetails.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		User: models.UserPublic{
			ID:              postWithDetails.Author.ID,
			Username:        postWithDetails.Author.Username,
			DisplayName:     postWithDetails.Author.DisplayName,
			Bio:             postWithDetails.Author.Bio,
			ProfileImageURL: postWithDetails.Author.ProfileImageURL,
			CoverImageURL:   postWithDetails.Author.CoverImageURL,
			Location:        postWithDetails.Author.Location,
			Website:         postWithDetails.Author.Website,
			IsVerified:      postWithDetails.Author.IsVerified,
			CreatedAt:       postWithDetails.Author.CreatedAt,
		},
		IsLiked: postWithDetails.IsLiked,
	}

	return commentResponse, nil
}

func (s *CommentService) GetComments(postID uuid.UUID, userID uuid.UUID, limit, offset int) ([]CommentResponse, error) {
	var posts []models.Post
	query := s.db.Where("parent_post_id = ? AND type = ? AND is_draft = false", postID, models.PostTypeReply).
		Preload("Author").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	var comments []CommentResponse
	for _, post := range posts {
		// Check if user liked this comment
		isLiked := false
		if userID != uuid.Nil {
			var like models.Like
			if err := s.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error; err == nil {
				isLiked = true
			}
		}

		comment := CommentResponse{
			ID:           post.ID,
			UserID:       post.AuthorID,
			PostID:       postID,
			Content:      post.Content,
			LikesCount:   post.LikesCount,
			RepliesCount: post.CommentsCount,
			CreatedAt:    post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			User: models.UserPublic{
				ID:              post.Author.ID,
				Username:        post.Author.Username,
				DisplayName:     post.Author.DisplayName,
				Bio:             post.Author.Bio,
				ProfileImageURL: post.Author.ProfileImageURL,
				CoverImageURL:   post.Author.CoverImageURL,
				Location:        post.Author.Location,
				Website:         post.Author.Website,
				IsVerified:      post.Author.IsVerified,
				CreatedAt:       post.Author.CreatedAt,
			},
			IsLiked: isLiked,
		}

		if post.ParentPostID != nil {
			comment.ParentID = post.ParentPostID
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *CommentService) GetReplies(commentID uuid.UUID, userID uuid.UUID, limit, offset int) ([]CommentResponse, error) {
	var posts []models.Post
	query := s.db.Where("parent_post_id = ? AND type = ? AND is_draft = false", commentID, models.PostTypeReply).
		Preload("Author").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch replies: %w", err)
	}

	// Get the original comment to find the post ID
	var parentComment models.Post
	if err := s.db.First(&parentComment, commentID).Error; err != nil {
		return nil, fmt.Errorf("failed to find parent comment: %w", err)
	}

	var replies []CommentResponse
	for _, post := range posts {
		// Check if user liked this reply
		isLiked := false
		if userID != uuid.Nil {
			var like models.Like
			if err := s.db.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&like).Error; err == nil {
				isLiked = true
			}
		}

		reply := CommentResponse{
			ID:           post.ID,
			UserID:       post.AuthorID,
			PostID:       *parentComment.ParentPostID, // The original post ID
			Content:      post.Content,
			ParentID:     &commentID,
			LikesCount:   post.LikesCount,
			RepliesCount: post.CommentsCount,
			CreatedAt:    post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			User: models.UserPublic{
				ID:              post.Author.ID,
				Username:        post.Author.Username,
				DisplayName:     post.Author.DisplayName,
				Bio:             post.Author.Bio,
				ProfileImageURL: post.Author.ProfileImageURL,
				CoverImageURL:   post.Author.CoverImageURL,
				Location:        post.Author.Location,
				Website:         post.Author.Website,
				IsVerified:      post.Author.IsVerified,
				CreatedAt:       post.Author.CreatedAt,
			},
			IsLiked: isLiked,
		}

		replies = append(replies, reply)
	}

	return replies, nil
}

func (s *CommentService) DeleteComment(commentID, userID uuid.UUID) error {
	// Use the existing post service to delete the comment
	return s.postService.DeletePost(commentID, userID)
}

func (s *CommentService) UpdateComment(commentID, userID uuid.UUID, req CreateCommentRequest) error {
	// Validate input
	if len(req.Content) == 0 {
		return errors.New("comment content is required")
	}
	if len(req.Content) > 280 {
		return errors.New("comment content exceeds 280 characters")
	}

	// Use the existing post service to update the comment
	updateReq := UpdatePostRequest{
		Content: req.Content,
	}
	return s.postService.UpdatePost(commentID, userID, updateReq)
}