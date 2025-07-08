package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostType string

const (
	PostTypeOriginal PostType = "original"
	PostTypeRepost   PostType = "repost"
	PostTypeQuote    PostType = "quote"
	PostTypeReply    PostType = "reply"
)

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"author_id"`
	Content     string    `gorm:"size:280" json:"content"`
	Type        PostType  `gorm:"default:'original'" json:"type"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	IsDraft     bool      `gorm:"default:false" json:"is_draft"`
	
	// For reposts and quotes
	OriginalPostID *uuid.UUID `gorm:"type:uuid;index" json:"original_post_id,omitempty"`
	
	// For replies
	ParentPostID *uuid.UUID `gorm:"type:uuid;index" json:"parent_post_id,omitempty"`
	
	// Metrics
	LikesCount    int `gorm:"default:0" json:"likes_count"`
	RepostsCount  int `gorm:"default:0" json:"reposts_count"`
	CommentsCount int `gorm:"default:0" json:"comments_count"`
	ViewsCount    int `gorm:"default:0" json:"views_count"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Author       User      `gorm:"foreignKey:AuthorID" json:"author"`
	OriginalPost *Post     `gorm:"foreignKey:OriginalPostID" json:"original_post,omitempty"`
	ParentPost   *Post     `gorm:"foreignKey:ParentPostID" json:"parent_post,omitempty"`
	Media        []Media   `gorm:"foreignKey:PostID" json:"media,omitempty"`
	Likes        []Like    `gorm:"foreignKey:PostID" json:"likes,omitempty"`
	Comments     []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Hashtags     []Hashtag `gorm:"many2many:post_hashtags;" json:"hashtags,omitempty"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PostWithDetails struct {
	Post
	IsLiked      bool `json:"is_liked"`
	IsReposted   bool `json:"is_reposted"`
	IsBookmarked bool `json:"is_bookmarked"`
}

type Hashtag struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Posts []Post `gorm:"many2many:post_hashtags;" json:"posts,omitempty"`
}

func (h *Hashtag) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}