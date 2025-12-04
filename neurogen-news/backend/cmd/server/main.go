package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/config"
	"github.com/neurogen-news/backend/internal/handler"
	appmiddleware "github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/repository"
	"github.com/neurogen-news/backend/internal/service"
	"github.com/neurogen-news/backend/internal/websocket"
	"github.com/neurogen-news/backend/pkg/logger"
)

//go:embed web/dist/*
var staticFiles embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger, err := logger.New(cfg.LogLevel, cfg.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Initialize database
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis
	redis, err := repository.NewRedisClient(cfg.RedisURL)
	if err != nil {
		zapLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redis.Close()

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(service.Deps{
		Repos:     repos,
		Redis:     redis,
		JWTSecret: cfg.JWTSecret,
		Logger:    zapLogger,
	})

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Neurogen.News API",
		ServerHeader:          "Neurogen.News",
		DisableStartupMessage: cfg.Environment == "production",
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           120 * time.Second,
		BodyLimit:             10 * 1024 * 1024, // 10MB
		ErrorHandler:          middleware.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(etag.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	app.Use(helmet.New())

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	// Initialize WebSocket hub
	wsHub := websocket.NewHub(redis, zapLogger)
	go wsHub.Run(context.Background())

	// Initialize WebSocket handler
	wsHandler := websocket.NewHandler(wsHub, services.Auth, zapLogger)

	// Initialize handlers
	handlers := handler.NewHandlers(services, zapLogger)

	// Setup routes
	setupRoutes(app, handlers, services, wsHandler, cfg)

	// Serve static files (Vue SPA)
	if cfg.ServeStatic {
		staticFS, err := fs.Sub(staticFiles, "web/dist")
		if err != nil {
			zapLogger.Fatal("Failed to load static files", zap.Error(err))
		}

		app.Use("/", filesystem.New(filesystem.Config{
			Root:         http.FS(staticFS),
			Browse:       false,
			Index:        "index.html",
			NotFoundFile: "index.html", // SPA fallback
		}))
	}

	// Start server
	go func() {
		addr := ":" + cfg.Port
		zapLogger.Info("Starting server", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			zapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited properly")
}

func setupRoutes(app *fiber.App, h *handler.Handlers, s *service.Services, ws *websocket.Handler, cfg *config.Config) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// WebSocket
	app.Use("/ws", websocket.UpgradeCheck())
	app.Get("/ws", appmiddleware.OptionalAuth(s.Auth), ws.Upgrade())

	// API v1
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)
	auth.Post("/logout", appmiddleware.Auth(s.Auth), h.Auth.Logout)
	auth.Post("/forgot-password", h.Auth.ForgotPassword)
	auth.Post("/reset-password", h.Auth.ResetPassword)

	// User routes
	users := api.Group("/users")
	users.Get("/me", appmiddleware.Auth(s.Auth), h.User.GetCurrentProfile)
	users.Put("/me", appmiddleware.Auth(s.Auth), h.User.UpdateProfile)
	users.Get("/:username", appmiddleware.OptionalAuth(s.Auth), h.User.GetProfile)
	users.Get("/:username/articles", h.User.GetArticles)
	users.Get("/:username/followers", h.User.GetFollowers)
	users.Get("/:username/following", h.User.GetFollowing)
	users.Post("/:username/follow", appmiddleware.Auth(s.Auth), h.User.Follow)
	users.Delete("/:username/follow", appmiddleware.Auth(s.Auth), h.User.Unfollow)

	// Article routes
	articles := api.Group("/articles")
	articles.Get("/", h.Article.List)
	articles.Get("/:id", h.Article.GetByID)
	articles.Get("/slug/:category/:slug", h.Article.GetBySlug)
	articles.Post("/", appmiddleware.Auth(s.Auth), h.Article.Create)
	articles.Put("/:id", appmiddleware.Auth(s.Auth), h.Article.Update)
	articles.Delete("/:id", appmiddleware.Auth(s.Auth), h.Article.Delete)
	articles.Post("/:id/reactions", appmiddleware.Auth(s.Auth), h.Article.AddReaction)
	articles.Delete("/:id/reactions", appmiddleware.Auth(s.Auth), h.Article.RemoveReaction)
	articles.Post("/:id/bookmark", appmiddleware.Auth(s.Auth), h.Article.Bookmark)
	articles.Delete("/:id/bookmark", appmiddleware.Auth(s.Auth), h.Article.RemoveBookmark)

	// Comment routes
	comments := api.Group("/comments")
	comments.Get("/article/:articleId", h.Comment.GetByArticle)
	comments.Get("/:id", h.Comment.GetByID)
	comments.Get("/:id/replies", h.Comment.GetReplies)
	comments.Post("/", appmiddleware.Auth(s.Auth), h.Comment.Create)
	comments.Put("/:id", appmiddleware.Auth(s.Auth), h.Comment.Update)
	comments.Delete("/:id", appmiddleware.Auth(s.Auth), h.Comment.Delete)
	comments.Post("/:id/reactions", appmiddleware.Auth(s.Auth), h.Comment.AddReaction)
	comments.Delete("/:id/reactions", appmiddleware.Auth(s.Auth), h.Comment.RemoveReaction)

	// Category routes
	categories := api.Group("/categories")
	categories.Get("/", h.Category.List)
	categories.Get("/subscriptions", appmiddleware.Auth(s.Auth), h.Category.GetUserSubscriptions)
	categories.Get("/:slug", appmiddleware.OptionalAuth(s.Auth), h.Category.GetBySlug)
	categories.Get("/:slug/articles", h.Category.GetArticles)
	categories.Post("/:slug/subscribe", appmiddleware.Auth(s.Auth), h.Category.Subscribe)
	categories.Delete("/:slug/subscribe", appmiddleware.Auth(s.Auth), h.Category.Unsubscribe)

	// Tag routes
	tags := api.Group("/tags")
	tags.Get("/", h.Tag.List)
	tags.Get("/popular", h.Tag.GetPopular)
	tags.Get("/search", h.Tag.Search)
	tags.Get("/:slug", h.Tag.GetBySlug)
	tags.Get("/:slug/articles", h.Tag.GetArticles)

	// Search routes
	search := api.Group("/search")
	search.Get("/", h.Search.Search)
	search.Get("/suggestions", h.Search.GetSuggestions)

	// Notification routes
	notifications := api.Group("/notifications")
	notifications.Use(appmiddleware.Auth(s.Auth))
	notifications.Get("/", h.Notification.List)
	notifications.Get("/unread", h.Notification.GetUnreadCount)
	notifications.Put("/:id/read", h.Notification.MarkAsRead)
	notifications.Put("/read-all", h.Notification.MarkAllAsRead)
	notifications.Delete("/:id", h.Notification.Delete)
	notifications.Delete("/", h.Notification.DeleteAll)

	// Bookmark routes
	bookmarks := api.Group("/bookmarks")
	bookmarks.Use(appmiddleware.Auth(s.Auth))
	bookmarks.Get("/", h.Bookmark.List)
	bookmarks.Post("/", h.Bookmark.Add)
	bookmarks.Delete("/:articleId", h.Bookmark.Remove)
	bookmarks.Get("/:articleId/check", h.Bookmark.Check)
	bookmarks.Get("/folders", h.Bookmark.ListFolders)
	bookmarks.Post("/folders", h.Bookmark.CreateFolder)
	bookmarks.Put("/folders/:id", h.Bookmark.UpdateFolder)
	bookmarks.Delete("/folders/:id", h.Bookmark.DeleteFolder)
	bookmarks.Put("/move", h.Bookmark.MoveToFolder)

	// Draft routes
	drafts := api.Group("/drafts")
	drafts.Use(appmiddleware.Auth(s.Auth))
	drafts.Get("/", h.Draft.List)
	drafts.Get("/autosave", h.Draft.GetAutoSave)
	drafts.Post("/autosave", h.Draft.AutoSave)
	drafts.Get("/:id", h.Draft.GetByID)
	drafts.Post("/", h.Draft.Create)
	drafts.Put("/:id", h.Draft.Update)
	drafts.Delete("/:id", h.Draft.Delete)

	// Achievement routes
	achievements := api.Group("/achievements")
	achievements.Get("/", h.Achievement.List)
	achievements.Get("/my", appmiddleware.Auth(s.Auth), h.Achievement.GetMyAchievements)
	achievements.Get("/progress", appmiddleware.Auth(s.Auth), h.Achievement.GetProgress)
	achievements.Post("/check", appmiddleware.Auth(s.Auth), h.Achievement.CheckAchievements)
	achievements.Get("/:id", h.Achievement.GetByID)
	achievements.Get("/user/:username", h.Achievement.GetUserAchievements)

	// Upload routes
	uploads := api.Group("/uploads")
	uploads.Use(appmiddleware.Auth(s.Auth))
	uploads.Post("/image", h.Upload.UploadImage)
	uploads.Delete("/", h.Upload.DeleteFile)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(appmiddleware.Auth(s.Auth))
	admin.Use(appmiddleware.RequireRole("ADMIN", "EDITOR", "MODERATOR"))
	admin.Get("/dashboard", h.Admin.GetDashboard)
	admin.Get("/users", h.Admin.GetUsers)
	admin.Post("/users/:id/ban", h.Admin.BanUser)
	admin.Delete("/users/:id/ban", h.Admin.UnbanUser)
	admin.Get("/reports", h.Admin.GetReports)
	admin.Put("/reports/:id", h.Admin.ResolveReport)
	admin.Delete("/articles/:id", h.Admin.DeleteArticle)
	admin.Delete("/comments/:id", h.Admin.DeleteComment)
	admin.Get("/settings", h.Admin.GetSettings)
	admin.Put("/settings", h.Admin.UpdateSettings)
}

