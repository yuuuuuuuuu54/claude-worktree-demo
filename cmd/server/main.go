package main

import (
	"digeon-backend/internal/config"
	"digeon-backend/internal/database"
	"digeon-backend/internal/handlers"
	"digeon-backend/internal/middleware"
	"digeon-backend/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database connection
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(db)
	mediaService := services.NewMediaService(db)
	postService := services.NewPostService(db, mediaService)
	timelineService := services.NewTimelineService(db)
	notificationService := services.NewNotificationService(db)
	likeService := services.NewLikeService(db, notificationService)
	followService := services.NewFollowService(db, notificationService)
	commentService := services.NewCommentService(db, postService, notificationService)
	searchService := services.NewSearchService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	timelineHandler := handlers.NewTimelineHandler(timelineService)
	likeHandler := handlers.NewLikeHandler(likeService)
	followHandler := handlers.NewFollowHandler(followService)
	commentHandler := handlers.NewCommentHandler(commentService)
	searchHandler := handlers.NewSearchHandler(searchService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	// Initialize Echo
	e := echo.New()

	// ミドルウェア
	e.Use(echo_middleware.Logger())
	e.Use(echo_middleware.Recover())
	
	// CORS設定
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}
	e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins: []string{corsOrigins},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"message": "Digeon Backend is running",
			"database": "connected",
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

	// 認証ルート
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/me", authHandler.Me, middleware.JWTMiddleware())

	// ユーザールート
	users := api.Group("/users")
	users.GET("/:id", userHandler.GetUserByID)
	users.GET("/username/:username", userHandler.GetUserByUsername)
	users.PUT("/profile", userHandler.UpdateProfile, middleware.JWTMiddleware())
	users.GET("/:user_id/posts", postHandler.GetUserPosts, middleware.OptionalJWTMiddleware())
	users.GET("/:user_id/likes", likeHandler.GetUserLikes, middleware.OptionalJWTMiddleware())
	users.POST("/:user_id/follow", followHandler.FollowUser, middleware.JWTMiddleware())
	users.DELETE("/:user_id/follow", followHandler.UnfollowUser, middleware.JWTMiddleware())
	users.GET("/:user_id/followers", followHandler.GetFollowers, middleware.OptionalJWTMiddleware())
	users.GET("/:user_id/following", followHandler.GetFollowing, middleware.OptionalJWTMiddleware())
	users.GET("/:user_id/follow-status", followHandler.CheckFollowStatus, middleware.JWTMiddleware())
	users.GET("/:user_id/follow-counts", followHandler.GetFollowCounts, middleware.OptionalJWTMiddleware())
	users.GET("/suggested", followHandler.GetSuggestedUsers, middleware.JWTMiddleware())

	// 投稿ルート
	posts := api.Group("/posts")
	posts.POST("", postHandler.CreatePost, middleware.JWTMiddleware())
	posts.GET("/:id", postHandler.GetPostByID, middleware.OptionalJWTMiddleware())
	posts.PUT("/:id", postHandler.UpdatePost, middleware.JWTMiddleware())
	posts.DELETE("/:id", postHandler.DeletePost, middleware.JWTMiddleware())
	posts.GET("/:post_id/replies", timelineHandler.GetPostReplies, middleware.OptionalJWTMiddleware())
	posts.POST("/:post_id/like", likeHandler.LikePost, middleware.JWTMiddleware())
	posts.DELETE("/:post_id/like", likeHandler.UnlikePost, middleware.JWTMiddleware())
	posts.GET("/:post_id/likes", likeHandler.GetPostLikes, middleware.OptionalJWTMiddleware())
	posts.GET("/:post_id/like-status", likeHandler.CheckLikeStatus, middleware.JWTMiddleware())
	posts.POST("/:post_id/comments", commentHandler.CreateComment, middleware.JWTMiddleware())
	posts.GET("/:post_id/comments", commentHandler.GetComments, middleware.OptionalJWTMiddleware())

	// コメントルート
	comments := api.Group("/comments")
	comments.PUT("/:comment_id", commentHandler.UpdateComment, middleware.JWTMiddleware())
	comments.DELETE("/:comment_id", commentHandler.DeleteComment, middleware.JWTMiddleware())
	comments.POST("/:comment_id/replies", commentHandler.CreateReply, middleware.JWTMiddleware())
	comments.GET("/:comment_id/replies", commentHandler.GetReplies, middleware.OptionalJWTMiddleware())

	// タイムラインルート
	timeline := api.Group("/timeline")
	timeline.GET("/home", timelineHandler.GetHomeTimeline, middleware.JWTMiddleware())
	timeline.GET("/explore", timelineHandler.GetExploreTimeline, middleware.OptionalJWTMiddleware())
	timeline.GET("/trending", timelineHandler.GetTrendingTimeline, middleware.OptionalJWTMiddleware())

	// 検索ルート
	search := api.Group("/search")
	search.GET("", searchHandler.SearchAll, middleware.OptionalJWTMiddleware())
	search.GET("/users", searchHandler.SearchUsers, middleware.OptionalJWTMiddleware())
	search.GET("/posts", searchHandler.SearchPosts, middleware.OptionalJWTMiddleware())
	search.GET("/hashtags", searchHandler.SearchHashtags, middleware.OptionalJWTMiddleware())
	search.GET("/hashtags/:hashtag/posts", searchHandler.GetHashtagPosts, middleware.OptionalJWTMiddleware())
	search.GET("/trending-hashtags", searchHandler.GetTrendingHashtags, middleware.OptionalJWTMiddleware())

	// 通知ルート
	notifications := api.Group("/notifications")
	notifications.GET("", notificationHandler.GetNotifications, middleware.JWTMiddleware())
	notifications.GET("/unread-count", notificationHandler.GetUnreadCount, middleware.JWTMiddleware())
	notifications.PUT("/:notification_id/read", notificationHandler.MarkAsRead, middleware.JWTMiddleware())
	notifications.PUT("/read-all", notificationHandler.MarkAllAsRead, middleware.JWTMiddleware())
	notifications.DELETE("/:notification_id", notificationHandler.DeleteNotification, middleware.JWTMiddleware())
	notifications.DELETE("/all", notificationHandler.DeleteAllNotifications, middleware.JWTMiddleware())

	// メディアルート
	media := api.Group("/media")
	media.POST("/upload", mediaHandler.UploadFile, middleware.JWTMiddleware())
	media.POST("/upload-multiple", mediaHandler.UploadMultipleFiles, middleware.JWTMiddleware())
	media.GET("/:media_id", mediaHandler.GetMedia, middleware.OptionalJWTMiddleware())
	media.DELETE("/:media_id", mediaHandler.DeleteMedia, middleware.JWTMiddleware())
	media.GET("/post/:post_id", mediaHandler.GetPostMedia, middleware.OptionalJWTMiddleware())

	// 静的ファイル配信
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	e.Static("/uploads", uploadDir)

	// サーバー起動
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on port %s", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}