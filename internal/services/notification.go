package services

import (
	"digeon-backend/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

type NotificationResponse struct {
	ID       uuid.UUID                `json:"id"`
	Type     models.NotificationType  `json:"type"`
	Message  string                   `json:"message"`
	IsRead   bool                     `json:"is_read"`
	Actor    models.UserPublic        `json:"actor"`
	Post     *models.Post             `json:"post,omitempty"`
	CreatedAt string                  `json:"created_at"`
}

// CreateFollowNotification creates a notification when someone follows a user
func (s *NotificationService) CreateFollowNotification(actorID, userID uuid.UUID) error {
	// Don't create notification if user follows themselves
	if actorID == userID {
		return nil
	}

	// Check if notification already exists
	var existingNotification models.Notification
	if err := s.db.Where("user_id = ? AND actor_id = ? AND type = ?", 
		userID, actorID, models.NotificationTypeFollow).First(&existingNotification).Error; err == nil {
		return nil // Notification already exists
	}

	notification := models.Notification{
		UserID:  userID,
		ActorID: actorID,
		Type:    models.NotificationTypeFollow,
		Message: "started following you",
		IsRead:  false,
	}

	return s.db.Create(&notification).Error
}

// CreateLikeNotification creates a notification when someone likes a post
func (s *NotificationService) CreateLikeNotification(actorID, postID uuid.UUID) error {
	// Get post owner
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		return err
	}

	// Don't create notification if user likes their own post
	if actorID == post.AuthorID {
		return nil
	}

	// Check if notification already exists
	var existingNotification models.Notification
	if err := s.db.Where("user_id = ? AND actor_id = ? AND post_id = ? AND type = ?", 
		post.AuthorID, actorID, postID, models.NotificationTypeLike).First(&existingNotification).Error; err == nil {
		return nil // Notification already exists
	}

	notification := models.Notification{
		UserID:  post.AuthorID,
		ActorID: actorID,
		Type:    models.NotificationTypeLike,
		PostID:  &postID,
		Message: "liked your post",
		IsRead:  false,
	}

	return s.db.Create(&notification).Error
}

// CreateCommentNotification creates a notification when someone comments on a post
func (s *NotificationService) CreateCommentNotification(actorID, postID uuid.UUID) error {
	// Get post owner
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		return err
	}

	// Don't create notification if user comments on their own post
	if actorID == post.AuthorID {
		return nil
	}

	notification := models.Notification{
		UserID:  post.AuthorID,
		ActorID: actorID,
		Type:    models.NotificationTypeComment,
		PostID:  &postID,
		Message: "commented on your post",
		IsRead:  false,
	}

	return s.db.Create(&notification).Error
}

// CreateRepostNotification creates a notification when someone reposts a post
func (s *NotificationService) CreateRepostNotification(actorID, postID uuid.UUID) error {
	// Get post owner
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		return err
	}

	// Don't create notification if user reposts their own post
	if actorID == post.AuthorID {
		return nil
	}

	notification := models.Notification{
		UserID:  post.AuthorID,
		ActorID: actorID,
		Type:    models.NotificationTypeRepost,
		PostID:  &postID,
		Message: "reposted your post",
		IsRead:  false,
	}

	return s.db.Create(&notification).Error
}

// CreateQuoteNotification creates a notification when someone quotes a post
func (s *NotificationService) CreateQuoteNotification(actorID, postID uuid.UUID) error {
	// Get post owner
	var post models.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		return err
	}

	// Don't create notification if user quotes their own post
	if actorID == post.AuthorID {
		return nil
	}

	notification := models.Notification{
		UserID:  post.AuthorID,
		ActorID: actorID,
		Type:    models.NotificationTypeQuote,
		PostID:  &postID,
		Message: "quoted your post",
		IsRead:  false,
	}

	return s.db.Create(&notification).Error
}

// GetNotifications gets notifications for a user
func (s *NotificationService) GetNotifications(userID uuid.UUID, limit, offset int) ([]NotificationResponse, error) {
	var notifications []models.Notification
	if err := s.db.Where("user_id = ?", userID).
		Preload("Actor").
		Preload("Post").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	var responses []NotificationResponse
	for _, notification := range notifications {
		response := NotificationResponse{
			ID:      notification.ID,
			Type:    notification.Type,
			Message: notification.Message,
			IsRead:  notification.IsRead,
			Actor: models.UserPublic{
				ID:              notification.Actor.ID,
				Username:        notification.Actor.Username,
				DisplayName:     notification.Actor.DisplayName,
				Bio:             notification.Actor.Bio,
				ProfileImageURL: notification.Actor.ProfileImageURL,
				CoverImageURL:   notification.Actor.CoverImageURL,
				Location:        notification.Actor.Location,
				Website:         notification.Actor.Website,
				IsVerified:      notification.Actor.IsVerified,
				CreatedAt:       notification.Actor.CreatedAt,
			},
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if notification.Post != nil {
			response.Post = notification.Post
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetUnreadNotificationsCount gets the count of unread notifications for a user
func (s *NotificationService) GetUnreadNotificationsCount(userID uuid.UUID) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *NotificationService) MarkNotificationAsRead(notificationID, userID uuid.UUID) error {
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	
	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}

	return nil
}

// MarkAllNotificationsAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(notificationID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{})
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}

	return nil
}

// DeleteAllNotifications deletes all notifications for a user
func (s *NotificationService) DeleteAllNotifications(userID uuid.UUID) error {
	if err := s.db.Where("user_id = ?", userID).Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf("failed to delete all notifications: %w", err)
	}
	return nil
}