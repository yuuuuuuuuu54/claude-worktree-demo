package database

import (
	"digeon-backend/internal/models"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")
	
	// Create database schema
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS digeon").Error; err != nil {
		return err
	}

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Hashtag{},
		&models.Like{},
		&models.Follow{},
		&models.Comment{},
		&models.Notification{},
		&models.Media{},
	)
	
	if err != nil {
		return err
	}
	
	// Create indexes for better performance
	if err := createIndexes(db); err != nil {
		return err
	}
	
	// Create unique constraints
	if err := createConstraints(db); err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

func createIndexes(db *gorm.DB) error {
	// Posts indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_author_created ON posts(author_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(type)")
	
	// Likes indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_post_like ON likes(user_id, post_id) WHERE deleted_at IS NULL")
	
	// Follows indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_follower_following ON follows(follower_id, following_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id)")
	
	// Comments indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_created ON comments(post_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_comment_id)")
	
	// Notifications indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(user_id, is_read, created_at DESC)")
	
	// Media indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_media_post_order ON media(post_id, \"order\")")
	
	// Full-text search indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_content_gin ON posts USING gin(to_tsvector('english', content))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username_gin ON users USING gin(to_tsvector('english', username))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_display_name_gin ON users USING gin(to_tsvector('english', display_name))")
	
	return nil
}

func createConstraints(db *gorm.DB) error {
	// User constraints
	db.Exec("ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_username_length CHECK (length(username) >= 3)")
	db.Exec("ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_email_format CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}$')")
	
	// Post constraints
	db.Exec("ALTER TABLE posts ADD CONSTRAINT IF NOT EXISTS chk_content_length CHECK (length(content) <= 280)")
	db.Exec("ALTER TABLE posts ADD CONSTRAINT IF NOT EXISTS chk_no_self_reply CHECK (id != parent_post_id)")
	
	// Comment constraints
	db.Exec("ALTER TABLE comments ADD CONSTRAINT IF NOT EXISTS chk_comment_length CHECK (length(content) <= 280)")
	db.Exec("ALTER TABLE comments ADD CONSTRAINT IF NOT EXISTS chk_no_self_reply CHECK (id != parent_comment_id)")
	
	// Follow constraints
	db.Exec("ALTER TABLE follows ADD CONSTRAINT IF NOT EXISTS chk_no_self_follow CHECK (follower_id != following_id)")
	
	// Media constraints
	db.Exec("ALTER TABLE media ADD CONSTRAINT IF NOT EXISTS chk_file_size CHECK (file_size >= 0)")
	db.Exec("ALTER TABLE media ADD CONSTRAINT IF NOT EXISTS chk_dimensions CHECK (width >= 0 AND height >= 0)")
	
	return nil
}