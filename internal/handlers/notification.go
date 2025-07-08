package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetNotifications(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	limit, offset := h.getPaginationParams(c)

	notifications, err := h.notificationService.GetNotifications(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch notifications")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": notifications,
		"limit":         limit,
		"offset":        offset,
	})
}

func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	count, err := h.notificationService.GetUnreadNotificationsCount(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch unread count")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"unread_count": count,
	})
}

func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	notificationIDParam := c.Param("notification_id")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid notification ID")
	}

	if err := h.notificationService.MarkNotificationAsRead(notificationID, userID); err != nil {
		if err.Error() == "notification not found or unauthorized" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to mark notification as read")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "notification marked as read",
	})
}

func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	if err := h.notificationService.MarkAllNotificationsAsRead(userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to mark all notifications as read")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "all notifications marked as read",
	})
}

func (h *NotificationHandler) DeleteNotification(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	notificationIDParam := c.Param("notification_id")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid notification ID")
	}

	if err := h.notificationService.DeleteNotification(notificationID, userID); err != nil {
		if err.Error() == "notification not found or unauthorized" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete notification")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "notification deleted",
	})
}

func (h *NotificationHandler) DeleteAllNotifications(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	if err := h.notificationService.DeleteAllNotifications(userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete all notifications")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "all notifications deleted",
	})
}

func (h *NotificationHandler) getPaginationParams(c echo.Context) (int, int) {
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