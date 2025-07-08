package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_following" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_following" json:"following_id"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower"`
	Following User `gorm:"foreignKey:FollowingID" json:"following"`
}

func (f *Follow) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// Ensure unique constraint on follower_id and following_id
func (Follow) TableName() string {
	return "follows"
}