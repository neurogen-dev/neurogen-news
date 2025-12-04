package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/service"
)

type ArticleHandler struct {
	articleService service.ArticleService
	logger         *zap.Logger
}

func NewArticleHandler(articleService service.ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		logger:         logger,
	}
}

func (h *ArticleHandler) List(c *fiber.Ctx) error {
	params := service.ArticleListParams{
		Sort:        c.Query("sort", "popular"),
		Level:       c.Query("level"),
		ContentType: c.Query("contentType"),
		CategoryID:  c.Query("categoryId"),
		TagID:       c.Query("tagId"),
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

func (h *ArticleHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	article, err := h.articleService.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found",
			})
		}
		h.logger.Error("Failed to get article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch article",
		})
	}

	// Record view
	userIP := c.IP()
	_ = h.articleService.RecordView(c.Context(), id, userIP)

	return c.JSON(article)
}

func (h *ArticleHandler) GetBySlug(c *fiber.Ctx) error {
	categorySlug := c.Params("category")
	articleSlug := c.Params("slug")

	article, err := h.articleService.GetBySlug(c.Context(), categorySlug, articleSlug)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found",
			})
		}
		h.logger.Error("Failed to get article by slug", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch article",
		})
	}

	// Record view
	userIP := c.IP()
	_ = h.articleService.RecordView(c.Context(), article.ID, userIP)

	return c.JSON(article)
}

type CreateArticleRequest struct {
	Title           string             `json:"title" validate:"required,min=5,max=200"`
	Content         string             `json:"content" validate:"required,min=100"`
	Lead            *string            `json:"lead,omitempty"`
	CoverImageURL   *string            `json:"coverImageUrl,omitempty"`
	Level           model.ArticleLevel `json:"level" validate:"required"`
	ContentType     model.ContentType  `json:"contentType" validate:"required"`
	CategoryID      string             `json:"categoryId" validate:"required"`
	Tags            []string           `json:"tags,omitempty"`
	IsNSFW          bool               `json:"isNsfw"`
	CommentsEnabled bool               `json:"commentsEnabled"`
	Status          model.ArticleStatus `json:"status"`
	MetaTitle       *string            `json:"metaTitle,omitempty"`
	MetaDescription *string            `json:"metaDescription,omitempty"`
}

func (h *ArticleHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Set default status
	if req.Status == "" {
		req.Status = model.StatusDraft
	}

	article, err := h.articleService.Create(c.Context(), userID, service.CreateArticleInput{
		Title:           req.Title,
		Content:         req.Content,
		Lead:            req.Lead,
		CoverImageURL:   req.CoverImageURL,
		Level:           req.Level,
		ContentType:     req.ContentType,
		CategoryID:      categoryID,
		Tags:            req.Tags,
		IsNSFW:          req.IsNSFW,
		CommentsEnabled: req.CommentsEnabled,
		Status:          req.Status,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	})

	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create article",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(article)
}

func (h *ArticleHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	var req service.UpdateArticleInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	article, err := h.articleService.Update(c.Context(), userID, id, req)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to edit this article",
			})
		}
		h.logger.Error("Failed to update article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update article",
		})
	}

	return c.JSON(article)
}

func (h *ArticleHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	if err := h.articleService.Delete(c.Context(), userID, id); err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Article not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to delete this article",
			})
		}
		h.logger.Error("Failed to delete article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete article",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Article deleted successfully",
	})
}

type AddReactionRequest struct {
	Emoji string `json:"emoji" validate:"required"`
}

func (h *ArticleHandler) AddReaction(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	articleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	var req AddReactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.articleService.AddReaction(c.Context(), userID, articleID, req.Emoji); err != nil {
		h.logger.Error("Failed to add reaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add reaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reaction added",
	})
}

func (h *ArticleHandler) RemoveReaction(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	articleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	if err := h.articleService.RemoveReaction(c.Context(), userID, articleID); err != nil {
		h.logger.Error("Failed to remove reaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove reaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reaction removed",
	})
}

func (h *ArticleHandler) Bookmark(c *fiber.Ctx) error {
	// TODO: Implement bookmark
	return c.JSON(fiber.Map{
		"message": "Bookmarked",
	})
}

func (h *ArticleHandler) RemoveBookmark(c *fiber.Ctx) error {
	// TODO: Implement remove bookmark
	return c.JSON(fiber.Map{
		"message": "Bookmark removed",
	})
}

