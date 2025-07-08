package middleware

import (
	"digeon-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString == auth {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			claims, err := utils.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Set user information in context
			c.Set(UserIDKey, claims.UserID)
			c.Set(UsernameKey, claims.Username)

			return next(c)
		}
	}
}

func OptionalJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			
			if auth != "" {
				tokenString := strings.TrimPrefix(auth, "Bearer ")
				if tokenString != auth {
					claims, err := utils.ValidateToken(tokenString)
					if err == nil {
						c.Set(UserIDKey, claims.UserID)
						c.Set(UsernameKey, claims.Username)
					}
				}
			}

			return next(c)
		}
	}
}