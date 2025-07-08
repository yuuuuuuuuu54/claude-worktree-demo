package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"message": "Digeon Backend is running",
		})
	})

	// API ルート
	api := e.Group("/api")
	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Digeon API",
			"version": "1.0.0",
		})
	})

	// サーバー起動
	log.Println("Server starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}