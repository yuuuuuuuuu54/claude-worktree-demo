package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeGif   MediaType = "gif"
)

type Media struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PostID      uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	Type        MediaType `gorm:"not null" json:"type"`
	URL         string    `gorm:"not null;size:500" json:"url"`
	ThumbnailURL string   `gorm:"size:500" json:"thumbnail_url"`
	FileName    string    `gorm:"size:255" json:"file_name"`
	FileSize    int64     `gorm:"default:0" json:"file_size"`
	Width       int       `gorm:"default:0" json:"width"`
	Height      int       `gorm:"default:0" json:"height"`
	Duration    int       `gorm:"default:0" json:"duration"` // For videos in seconds
	AltText     string    `gorm:"size:500" json:"alt_text"`
	Order       int       `gorm:"default:0" json:"order"` // For ordering multiple media in a post
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Post Post `gorm:"foreignKey:PostID" json:"post"`
}

func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}