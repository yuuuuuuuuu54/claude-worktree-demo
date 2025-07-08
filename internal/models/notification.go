package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeFollow    NotificationType = "follow"
	NotificationTypeLike      NotificationType = "like"
	NotificationTypeComment   NotificationType = "comment"
	NotificationTypeRepost    NotificationType = "repost"
	NotificationTypeQuote     NotificationType = "quote"
	NotificationTypeMention   NotificationType = "mention"
)

type Notification struct {
	ID       uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID   uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	ActorID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"actor_id"`
	Type     NotificationType `gorm:"not null" json:"type"`
	PostID   *uuid.UUID       `gorm:"type:uuid;index" json:"post_id,omitempty"`
	Message  string           `gorm:"size:500" json:"message"`
	IsRead   bool             `gorm:"default:false" json:"is_read"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user"`
	Actor User  `gorm:"foreignKey:ActorID" json:"actor"`
	Post  *Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}