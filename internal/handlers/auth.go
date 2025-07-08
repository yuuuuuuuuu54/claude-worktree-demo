package handlers

import (
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"digeon-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req services.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	response, err := h.userService.Register(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req services.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	response, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")
	if tokenString == auth {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
	}

	newToken, err := utils.RefreshToken(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": newToken,
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, ok := c.Get(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// In a more advanced implementation, you might want to blacklist the token
	// For now, we'll just return a success response
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}