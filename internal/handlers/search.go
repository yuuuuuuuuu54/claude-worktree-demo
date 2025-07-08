package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) SearchAll(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query is required")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.searchService.SearchAll(query, userID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "search failed")
	}

	return c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query is required")
	}

	limit, offset := h.getPaginationParams(c)

	results, err := h.searchService.SearchUsers(query, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "user search failed")
	}

	return c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) SearchPosts(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query is required")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	results, err := h.searchService.SearchPosts(query, userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "post search failed")
	}

	return c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) SearchHashtags(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query is required")
	}

	limit, offset := h.getPaginationParams(c)

	results, err := h.searchService.SearchHashtags(query, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "hashtag search failed")
	}

	return c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) GetHashtagPosts(c echo.Context) error {
	hashtag := c.Param("hashtag")
	if hashtag == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "hashtag is required")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	results, err := h.searchService.GetHashtagPosts(hashtag, userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch hashtag posts")
	}

	return c.JSON(http.StatusOK, results)
}

func (h *SearchHandler) GetTrendingHashtags(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	hashtags, err := h.searchService.GetTrendingHashtags(limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch trending hashtags")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"hashtags": hashtags,
		"limit":    limit,
	})
}

func (h *SearchHandler) getPaginationParams(c echo.Context) (int, int) {
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