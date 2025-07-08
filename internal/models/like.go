package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_post_like" json:"user_id"`
	PostID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_post_like" json:"post_id"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user"`
	Post Post `gorm:"foreignKey:PostID" json:"post"`
}

func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// Ensure unique constraint on user_id and post_id
func (Like) TableName() string {
	return "likes"
}