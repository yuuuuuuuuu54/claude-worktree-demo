package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	var req services.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	comment, err := h.commentService.CreateComment(userID, postID, req)
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) CreateReply(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	commentIDParam := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment ID")
	}

	var req services.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	reply, err := h.commentService.CreateReply(userID, commentID, req)
	if err != nil {
		if err.Error() == "comment not found" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, reply)
}

func (h *CommentHandler) GetComments(c echo.Context) error {
	postIDParam := c.Param("post_id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid post ID")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	comments, err := h.commentService.GetComments(postID, userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch comments")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"comments": comments,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *CommentHandler) GetReplies(c echo.Context) error {
	commentIDParam := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment ID")
	}

	// Get user ID from context (optional)
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}

	limit, offset := h.getPaginationParams(c)

	replies, err := h.commentService.GetReplies(commentID, userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch replies")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"replies": replies,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	commentIDParam := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment ID")
	}

	var req services.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.commentService.UpdateComment(commentID, userID, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "comment updated successfully",
	})
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	commentIDParam := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid comment ID")
	}

	if err := h.commentService.DeleteComment(commentID, userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "comment deleted successfully",
	})
}

func (h *CommentHandler) getPaginationParams(c echo.Context) (int, int) {
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