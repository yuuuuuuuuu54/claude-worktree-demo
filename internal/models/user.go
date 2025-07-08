package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username        string    `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email           string    `gorm:"uniqueIndex;not null;size:100" json:"email"`
	PasswordHash    string    `gorm:"not null;size:255" json:"-"`
	DisplayName     string    `gorm:"size:100" json:"display_name"`
	Bio             string    `gorm:"size:280" json:"bio"`
	ProfileImageURL string    `gorm:"size:500" json:"profile_image_url"`
	CoverImageURL   string    `gorm:"size:500" json:"cover_image_url"`
	Location        string    `gorm:"size:100" json:"location"`
	Website         string    `gorm:"size:200" json:"website"`
	IsVerified      bool      `gorm:"default:false" json:"is_verified"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Posts         []Post         `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Likes         []Like         `gorm:"foreignKey:UserID" json:"likes,omitempty"`
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	
	// Following relationships
	Following []Follow `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	Followers []Follow `gorm:"foreignKey:FollowingID" json:"followers,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserPublic struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	DisplayName     string    `json:"display_name"`
	Bio             string    `json:"bio"`
	ProfileImageURL string    `json:"profile_image_url"`
	CoverImageURL   string    `json:"cover_image_url"`
	Location        string    `json:"location"`
	Website         string    `json:"website"`
	IsVerified      bool      `json:"is_verified"`
	CreatedAt       time.Time `json:"created_at"`
	FollowersCount  int64     `json:"followers_count"`
	FollowingCount  int64     `json:"following_count"`
	PostsCount      int64     `json:"posts_count"`
}