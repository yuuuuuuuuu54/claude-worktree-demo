package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PostID  uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	Content string    `gorm:"size:280;not null" json:"content"`
	
	// For nested comments (replies to comments)
	ParentCommentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_comment_id,omitempty"`
	
	// Metrics
	LikesCount int `gorm:"default:0" json:"likes_count"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Post          Post      `gorm:"foreignKey:PostID" json:"post"`
	ParentComment *Comment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []Comment `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CommentWithDetails struct {
	Comment
	IsLiked bool `json:"is_liked"`
}