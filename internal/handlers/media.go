package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MediaHandler struct {
	mediaService *services.MediaService
}

func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

func (h *MediaHandler) UploadFile(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no file uploaded")
	}

	// Upload file
	response, err := h.mediaService.UploadFile(file, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *MediaHandler) UploadMultipleFiles(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse multipart form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no files uploaded")
	}

	// Limit number of files
	if len(files) > 4 {
		return echo.NewHTTPError(http.StatusBadRequest, "maximum 4 files allowed")
	}

	var responses []services.UploadResponse
	for _, file := range files {
		response, err := h.mediaService.UploadFile(file, userID)
		if err != nil {
			// If one file fails, continue with others but report the error
			continue
		}
		responses = append(responses, *response)
	}

	if len(responses) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "all file uploads failed")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uploaded_files": responses,
		"total_uploaded": len(responses),
		"total_files":    len(files),
	})
}

func (h *MediaHandler) GetMedia(c echo.Context) error {
	mediaIDParam := c.Param("media_id")
	mediaID, err := uuid.Parse(mediaIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid media ID")
	}

	media, err := h.mediaService.GetMediaByID(mediaID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "media not found")
	}

	return c.JSON(http.StatusOK, media)
}

func (h *MediaHandler) DeleteMedia(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	mediaIDParam := c.Param("media_id")
	mediaID, err := uuid.Parse(mediaIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid media ID")
	}

	if err := h.mediaService.DeleteMedia(mediaID, userID); err != nil {
		if err.Error() == "unauthorized to delete this media" {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "media deleted successfully",
	})
}

func (h *MediaHandler) GetPostMedia(c echo.Context) error {
	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	media, err := h.mediaService.GetPostMedia(postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch post media")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"media": media,
	})
}