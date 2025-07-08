package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req services.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	post, err := h.postService.CreatePost(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostByID(c echo.Context) error {
	postIDParam := c.Param("id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		// If user is not authenticated, just return the post without user-specific details
		post, err := h.postService.GetPostByID(postID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "post not found")
		}
		return c.JSON(http.StatusOK, post)
	}

	post, err := h.postService.GetPostWithDetails(postID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "post not found")
	}

	return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetUserPosts(c echo.Context) error {
	userIDParam := c.Param("user_id")
	targetUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	// Get pagination parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := h.postService.GetPostsByUserID(targetUserID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch posts")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	var req services.UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.postService.UpdatePost(postID, userID, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "post updated successfully",
	})
}

func (h *PostHandler) DeletePost(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	if err := h.postService.DeletePost(postID, userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "post deleted successfully",
	})
}