package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/service"
)

type TagHandler struct {
	tagService     service.TagService
	articleService service.ArticleService
	logger         *zap.Logger
}

func NewTagHandler(tagService service.TagService, articleService service.ArticleService, logger *zap.Logger) *TagHandler {
	return &TagHandler{
		tagService:     tagService,
		articleService: articleService,
		logger:         logger,
	}
}

// List returns all tags
func (h *TagHandler) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)

	tags, err := h.tagService.GetAll(c.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get tags", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags",
		})
	}

	return c.JSON(fiber.Map{
		"items": tags,
	})
}

// GetPopular returns popular tags
func (h *TagHandler) GetPopular(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)

	tags, err := h.tagService.GetPopular(c.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get popular tags", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags",
		})
	}

	return c.JSON(fiber.Map{
		"items": tags,
	})
}

// GetBySlug returns a tag by slug with its articles
func (h *TagHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	tag, err := h.tagService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "tag not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tag not found",
			})
		}
		h.logger.Error("Failed to get tag", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tag",
		})
	}

	return c.JSON(tag)
}

// GetArticles returns articles for a tag
func (h *TagHandler) GetArticles(c *fiber.Ctx) error {
	slug := c.Params("slug")

	tag, err := h.tagService.GetBySlug(c.Context(), slug)
	if err != nil {
		if err.Error() == "tag not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tag not found",
			})
		}
		h.logger.Error("Failed to get tag", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tag",
		})
	}

	params := service.ArticleListParams{
		Sort:      c.Query("sort", "popular"),
		Level:     c.Query("level"),
		TagID:     tag.ID.String(),
		TimeRange: c.Query("timeRange", "all"),
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("pageSize", 20),
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

// Search searches for tags
func (h *TagHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.JSON(fiber.Map{
			"items": []interface{}{},
		})
	}

	limit := c.QueryInt("limit", 10)

	tags, err := h.tagService.Search(c.Context(), query, limit)
	if err != nil {
		h.logger.Error("Failed to search tags", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search tags",
		})
	}

	return c.JSON(fiber.Map{
		"items": tags,
	})
}

// GetByID returns a tag by ID
func (h *TagHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tag ID",
		})
	}

	tag, err := h.tagService.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "tag not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tag not found",
			})
		}
		h.logger.Error("Failed to get tag", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tag",
		})
	}

	return c.JSON(tag)
}

