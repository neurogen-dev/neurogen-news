package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
	articleService  service.ArticleService
	logger          *zap.Logger
}

func NewCategoryHandler(
	categoryService service.CategoryService,
	articleService service.ArticleService,
	logger *zap.Logger,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		articleService:  articleService,
		logger:          logger,
	}
}

// List returns all categories
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAll(c.Context())
	if err != nil {
		h.logger.Error("Failed to get categories", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(fiber.Map{
		"items": categories,
	})
}

// GetBySlug returns a category by slug with its articles
func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		h.logger.Error("Failed to get category", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category",
		})
	}

	// Check if user is subscribed
	var isSubscribed bool
	if userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID); ok {
		isSubscribed, _ = h.categoryService.IsSubscribed(c.Context(), userID, category.ID)
	}

	return c.JSON(fiber.Map{
		"category":     category,
		"isSubscribed": isSubscribed,
	})
}

// GetArticles returns articles for a category
func (h *CategoryHandler) GetArticles(c *fiber.Ctx) error {
	slug := c.Params("slug")

	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		h.logger.Error("Failed to get category", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category",
		})
	}

	params := service.ArticleListParams{
		Sort:        c.Query("sort", "popular"),
		Level:       c.Query("level"),
		ContentType: c.Query("contentType"),
		CategoryID:  category.ID.String(),
		TimeRange:   c.Query("timeRange", "all"),
		Page:        c.QueryInt("page", 1),
		PageSize:    c.QueryInt("pageSize", 20),
	}

	result, err := h.articleService.List(c.Context(), params)
	if err != nil {
		h.logger.Error("Failed to list articles", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch articles",
		})
	}

	return c.JSON(result)
}

// Subscribe subscribes user to a category
func (h *CategoryHandler) Subscribe(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	slug := c.Params("slug")
	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category",
		})
	}

	if err := h.categoryService.Subscribe(c.Context(), userID, category.ID); err != nil {
		h.logger.Error("Failed to subscribe", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to subscribe",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Subscribed successfully",
		"isSubscribed": true,
	})
}

// Unsubscribe unsubscribes user from a category
func (h *CategoryHandler) Unsubscribe(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	slug := c.Params("slug")
	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category",
		})
	}

	if err := h.categoryService.Unsubscribe(c.Context(), userID, category.ID); err != nil {
		h.logger.Error("Failed to unsubscribe", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unsubscribe",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Unsubscribed successfully",
		"isSubscribed": false,
	})
}

// GetUserSubscriptions returns user's category subscriptions
func (h *CategoryHandler) GetUserSubscriptions(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	categories, err := h.categoryService.GetUserSubscriptions(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get subscriptions", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subscriptions",
		})
	}

	return c.JSON(fiber.Map{
		"items": categories,
	})
}
