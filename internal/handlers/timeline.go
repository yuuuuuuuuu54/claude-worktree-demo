package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TimelineHandler struct {
	timelineService *services.TimelineService
}

func NewTimelineHandler(timelineService *services.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

func (h *TimelineHandler) GetHomeTimeline(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	limit, offset := h.getPaginationParams(c)

	timeline, err := h.timelineService.GetHomeTimeline(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch home timeline")
	}

	return c.JSON(http.StatusOK, timeline)
}

func (h *TimelineHandler) GetExploreTimeline(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		// For anonymous users, we'll use a zero UUID
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	timeline, err := h.timelineService.GetExploreTimeline(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch explore timeline")
	}

	return c.JSON(http.StatusOK, timeline)
}

func (h *TimelineHandler) GetTrendingTimeline(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		// For anonymous users, we'll use a zero UUID
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	timeline, err := h.timelineService.GetTrendingTimeline(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch trending timeline")
	}

	return c.JSON(http.StatusOK, timeline)
}

func (h *TimelineHandler) GetPostReplies(c echo.Context) error {
	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		// For anonymous users, we'll use a zero UUID
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	replies, err := h.timelineService.GetPostReplies(postID, userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch post replies")
	}

	return c.JSON(http.StatusOK, replies)
}

func (h *TimelineHandler) getPaginationParams(c echo.Context) (int, int) {
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