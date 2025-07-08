package services

import (
	"digeon-backend/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaService struct {
	db        *gorm.DB
	uploadDir string
	baseURL   string
}

func NewMediaService(db *gorm.DB) *MediaService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Create upload directory if it doesn't exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(filepath.Join(uploadDir, "images"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "videos"), 0755)

	return &MediaService{
		db:        db,
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

type UploadResponse struct {
	ID           uuid.UUID           `json:"id"`
	URL          string              `json:"url"`
	ThumbnailURL string              `json:"thumbnail_url,omitempty"`
	Type         models.MediaType    `json:"type"`
	FileName     string              `json:"file_name"`
	FileSize     int64               `json:"file_size"`
	Width        int                 `json:"width,omitempty"`
	Height       int                 `json:"height,omitempty"`
	Duration     int                 `json:"duration,omitempty"`
}

func (s *MediaService) UploadFile(file *multipart.FileHeader, userID uuid.UUID) (*UploadResponse, error) {
	// Validate file type
	mediaType, err := s.getMediaType(file.Filename)
	if err != nil {
		return nil, fmt.Errorf("unsupported file type: %w", err)
	}

	// Validate file size
	maxSize := s.getMaxFileSize(mediaType)
	if file.Size > maxSize {
		return nil, fmt.Errorf("file size exceeds limit of %d bytes", maxSize)
	}

	// Generate unique filename
	fileID := uuid.New()
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", fileID.String(), ext)

	// Determine subdirectory based on media type
	var subDir string
	switch mediaType {
	case models.MediaTypeImage, models.MediaTypeGif:
		subDir = "images"
	case models.MediaTypeVideo:
		subDir = "videos"
	default:
		subDir = "other"
	}

	// Create full file path
	filePath := filepath.Join(s.uploadDir, subDir, fileName)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("%s/uploads/%s/%s", s.baseURL, subDir, fileName)

	// Create media record
	media := models.Media{
		ID:       fileID,
		Type:     mediaType,
		URL:      url,
		FileName: file.Filename,
		FileSize: file.Size,
	}

	// Save to database
	if err := s.db.Create(&media).Error; err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	response := &UploadResponse{
		ID:       media.ID,
		URL:      media.URL,
		Type:     media.Type,
		FileName: media.FileName,
		FileSize: media.FileSize,
		Width:    media.Width,
		Height:   media.Height,
		Duration: media.Duration,
	}

	if media.ThumbnailURL != "" {
		response.ThumbnailURL = media.ThumbnailURL
	}

	return response, nil
}

func (s *MediaService) AttachMediaToPost(mediaIDs []string, postID uuid.UUID) error {
	for i, mediaIDStr := range mediaIDs {
		mediaID, err := uuid.Parse(mediaIDStr)
		if err != nil {
			continue
		}

		// Update media with post ID and order
		if err := s.db.Model(&models.Media{}).
			Where("id = ? AND post_id IS NULL", mediaID).
			Updates(map[string]interface{}{
				"post_id": postID,
				"order":   i,
			}).Error; err != nil {
			return fmt.Errorf("failed to attach media to post: %w", err)
		}
	}

	return nil
}

func (s *MediaService) GetMediaByID(mediaID uuid.UUID) (*models.Media, error) {
	var media models.Media
	if err := s.db.First(&media, mediaID).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (s *MediaService) DeleteMedia(mediaID uuid.UUID, userID uuid.UUID) error {
	var media models.Media
	if err := s.db.First(&media, mediaID).Error; err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	// Check if media is attached to a post and if user owns the post
	if media.PostID != uuid.Nil {
		var post models.Post
		if err := s.db.First(&post, media.PostID).Error; err != nil {
			return fmt.Errorf("associated post not found: %w", err)
		}
		if post.AuthorID != userID {
			return fmt.Errorf("unauthorized to delete this media")
		}
	}

	// Delete file from filesystem
	if err := s.deleteFileFromPath(media.URL); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete file from filesystem: %v\n", err)
	}

	// Delete from database
	if err := s.db.Delete(&media).Error; err != nil {
		return fmt.Errorf("failed to delete media record: %w", err)
	}

	return nil
}

func (s *MediaService) GetPostMedia(postID uuid.UUID) ([]models.Media, error) {
	var media []models.Media
	if err := s.db.Where("post_id = ?", postID).Order("\"order\" ASC").Find(&media).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch post media: %w", err)
	}
	return media, nil
}

func (s *MediaService) getMediaType(filename string) (models.MediaType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return models.MediaTypeImage, nil
	case ".gif":
		return models.MediaTypeGif, nil
	case ".mp4", ".mov", ".avi", ".mkv", ".webm":
		return models.MediaTypeVideo, nil
	default:
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}
}

func (s *MediaService) getMaxFileSize(mediaType models.MediaType) int64 {
	switch mediaType {
	case models.MediaTypeImage, models.MediaTypeGif:
		maxSizeStr := os.Getenv("MAX_FILE_SIZE")
		if maxSizeStr == "" {
			return 5 * 1024 * 1024 // 5MB default
		}
		// In a real implementation, you'd parse this properly
		return 5 * 1024 * 1024
	case models.MediaTypeVideo:
		maxSizeStr := os.Getenv("MAX_VIDEO_SIZE")
		if maxSizeStr == "" {
			return 100 * 1024 * 1024 // 100MB default
		}
		// In a real implementation, you'd parse this properly
		return 100 * 1024 * 1024
	default:
		return 5 * 1024 * 1024
	}
}

func (s *MediaService) deleteFileFromPath(url string) error {
	// Extract file path from URL
	// This assumes URL format: http://localhost:8080/uploads/images/filename.jpg
	urlParts := strings.Split(url, "/uploads/")
	if len(urlParts) != 2 {
		return fmt.Errorf("invalid URL format")
	}
	
	relativePath := urlParts[1]
	fullPath := filepath.Join(s.uploadDir, relativePath)
	
	return os.Remove(fullPath)
}

// CleanupOrphanedMedia removes media files that are not attached to any post and are older than 24 hours
func (s *MediaService) CleanupOrphanedMedia() error {
	var orphanedMedia []models.Media
	
	// Find media that's not attached to posts and is older than 24 hours
	if err := s.db.Where("post_id IS NULL AND created_at < ?", time.Now().Add(-24*time.Hour)).Find(&orphanedMedia).Error; err != nil {
		return fmt.Errorf("failed to find orphaned media: %w", err)
	}

	for _, media := range orphanedMedia {
		// Delete file from filesystem
		if err := s.deleteFileFromPath(media.URL); err != nil {
			fmt.Printf("Warning: failed to delete orphaned file %s: %v\n", media.URL, err)
		}

		// Delete from database
		if err := s.db.Delete(&media).Error; err != nil {
			fmt.Printf("Warning: failed to delete orphaned media record %s: %v\n", media.ID, err)
		}
	}

	return nil
}