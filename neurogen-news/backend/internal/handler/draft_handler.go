package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/service"
)

type DraftHandler struct {
	draftService service.DraftService
	logger       *zap.Logger
}

func NewDraftHandler(draftService service.DraftService, logger *zap.Logger) *DraftHandler {
	return &DraftHandler{
		draftService: draftService,
		logger:       logger,
	}
}

// List returns user's drafts
func (h *DraftHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	result, err := h.draftService.GetByUser(c.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get drafts", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch drafts",
		})
	}

	return c.JSON(result)
}

// GetByID returns a specific draft
func (h *DraftHandler) GetByID(c *fiber.Ctx) error {
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
			"error": "Invalid draft ID",
		})
	}

	draft, err := h.draftService.GetByID(c.Context(), userID, id)
	if err != nil {
		if err.Error() == "draft not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Draft not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to view this draft",
			})
		}
		h.logger.Error("Failed to get draft", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch draft",
		})
	}

	return c.JSON(draft)
}

type CreateDraftRequest struct {
	Title         string            `json:"title" validate:"required,min=1,max=200"`
	Content       string            `json:"content"`
	CoverImageURL *string           `json:"coverImageUrl,omitempty"`
	ContentType   model.ContentType `json:"contentType" validate:"required"`
	CategoryID    *string           `json:"categoryId,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
}

// Create creates a new draft
func (h *DraftHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		cid, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		categoryID = &cid
	}

	draft, err := h.draftService.Create(c.Context(), userID, service.CreateDraftInput{
		Title:         req.Title,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		ContentType:   req.ContentType,
		CategoryID:    categoryID,
		Tags:          req.Tags,
	})
	if err != nil {
		h.logger.Error("Failed to create draft", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create draft",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(draft)
}

type UpdateDraftRequest struct {
	Title         *string            `json:"title,omitempty"`
	Content       *string            `json:"content,omitempty"`
	CoverImageURL *string            `json:"coverImageUrl,omitempty"`
	ContentType   *model.ContentType `json:"contentType,omitempty"`
	CategoryID    *string            `json:"categoryId,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

// Update updates a draft
func (h *DraftHandler) Update(c *fiber.Ctx) error {
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
			"error": "Invalid draft ID",
		})
	}

	var req UpdateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		cid, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		categoryID = &cid
	}

	draft, err := h.draftService.Update(c.Context(), userID, id, service.UpdateDraftInput{
		Title:         req.Title,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		ContentType:   req.ContentType,
		CategoryID:    categoryID,
		Tags:          req.Tags,
	})
	if err != nil {
		if err.Error() == "draft not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Draft not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to edit this draft",
			})
		}
		h.logger.Error("Failed to update draft", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update draft",
		})
	}

	return c.JSON(draft)
}

// Delete deletes a draft
func (h *DraftHandler) Delete(c *fiber.Ctx) error {
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
			"error": "Invalid draft ID",
		})
	}

	if err := h.draftService.Delete(c.Context(), userID, id); err != nil {
		if err.Error() == "draft not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Draft not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to delete this draft",
			})
		}
		h.logger.Error("Failed to delete draft", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete draft",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Draft deleted successfully",
	})
}

type AutoSaveRequest struct {
	ArticleID     *string           `json:"articleId,omitempty"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	CoverImageURL *string           `json:"coverImageUrl,omitempty"`
	ContentType   model.ContentType `json:"contentType"`
	CategoryID    *string           `json:"categoryId,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
}

// AutoSave auto-saves a draft
func (h *DraftHandler) AutoSave(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req AutoSaveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var articleID *uuid.UUID
	if req.ArticleID != nil {
		aid, err := uuid.Parse(*req.ArticleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid article ID",
			})
		}
		articleID = &aid
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		cid, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		categoryID = &cid
	}

	if err := h.draftService.AutoSave(c.Context(), userID, service.AutoSaveInput{
		ArticleID:     articleID,
		Title:         req.Title,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		ContentType:   req.ContentType,
		CategoryID:    categoryID,
		Tags:          req.Tags,
	}); err != nil {
		h.logger.Error("Failed to auto-save", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to auto-save",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Auto-saved successfully",
	})
}

// GetAutoSave returns the latest auto-save
func (h *DraftHandler) GetAutoSave(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var articleID *uuid.UUID
	if articleIDStr := c.Query("articleId"); articleIDStr != "" {
		aid, err := uuid.Parse(articleIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid article ID",
			})
		}
		articleID = &aid
	}

	draft, err := h.draftService.GetLatestAutoSave(c.Context(), userID, articleID)
	if err != nil {
		if err.Error() == "draft not found" {
			return c.JSON(fiber.Map{
				"draft": nil,
			})
		}
		h.logger.Error("Failed to get auto-save", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get auto-save",
		})
	}

	return c.JSON(fiber.Map{
		"draft": draft,
	})
}

