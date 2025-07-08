package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LikeHandler struct {
	likeService *services.LikeService
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

func (h *LikeHandler) LikePost(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	if err := h.likeService.LikePost(userID, postID); err != nil {
		if err.Error() == "post already liked" {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if err.Error() == "post not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to like post")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "post liked successfully",
	})
}

func (h *LikeHandler) UnlikePost(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	if err := h.likeService.UnlikePost(userID, postID); err != nil {
		if err.Error() == "like not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err.Error() == "post not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unlike post")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "post unliked successfully",
	})
}

func (h *LikeHandler) GetPostLikes(c echo.Context) error {
	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	limit, offset := h.getPaginationParams(c)

	users, err := h.likeService.GetPostLikes(postID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch post likes")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *LikeHandler) GetUserLikes(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	limit, offset := h.getPaginationParams(c)

	posts, err := h.likeService.GetUserLikes(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch user likes")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *LikeHandler) CheckLikeStatus(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	isLiked, err := h.likeService.IsPostLikedByUser(userID, postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check like status")
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"is_liked": isLiked,
	})
}

func (h *LikeHandler) getPaginationParams(c echo.Context) (int, int) {
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

	return limit, offset
}