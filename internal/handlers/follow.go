package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler(followService *services.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

func (h *FollowHandler) FollowUser(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	followingIDParam := c.Param("user_id")
	followingID, err := uuid.Parse(followingIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if err := h.followService.FollowUser(userID, followingID); err != nil {
		if err.Error() == "cannot follow yourself" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err.Error() == "already following user" {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		if err.Error() == "user not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to follow user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user followed successfully",
	})
}

func (h *FollowHandler) UnfollowUser(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	followingIDParam := c.Param("user_id")
	followingID, err := uuid.Parse(followingIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if err := h.followService.UnfollowUser(userID, followingID); err != nil {
		if err.Error() == "cannot unfollow yourself" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err.Error() == "follow relationship not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unfollow user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user unfollowed successfully",
	})
}

func (h *FollowHandler) GetFollowers(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	limit, offset := h.getPaginationParams(c)

	followers, err := h.followService.GetFollowers(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch followers")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"followers": followers,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *FollowHandler) GetFollowing(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	limit, offset := h.getPaginationParams(c)

	following, err := h.followService.GetFollowing(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch following")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"following": following,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *FollowHandler) CheckFollowStatus(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	followingIDParam := c.Param("user_id")
	followingID, err := uuid.Parse(followingIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	isFollowing, err := h.followService.IsFollowing(userID, followingID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check follow status")
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"is_following": isFollowing,
	})
}

func (h *FollowHandler) GetFollowCounts(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	followersCount, followingCount, err := h.followService.GetFollowCounts(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch follow counts")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"followers_count": followersCount,
		"following_count": followingCount,
	})
}

func (h *FollowHandler) GetSuggestedUsers(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	users, err := h.followService.GetSuggestedUsers(userID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch suggested users")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"limit": limit,
	})
}

func (h *FollowHandler) getPaginationParams(c echo.Context) (int, int) {
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