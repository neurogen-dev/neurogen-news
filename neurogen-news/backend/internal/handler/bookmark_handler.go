package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type BookmarkHandler struct {
	bookmarkService service.BookmarkService
	logger          *zap.Logger
}

func NewBookmarkHandler(bookmarkService service.BookmarkService, logger *zap.Logger) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkService: bookmarkService,
		logger:          logger,
	}
}

// List returns user's bookmarks
func (h *BookmarkHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var folderID *uuid.UUID
	if folderIDStr := c.Query("folderId"); folderIDStr != "" {
		fid, err := uuid.Parse(folderIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid folder ID",
			})
		}
		folderID = &fid
	}

	params := service.BookmarkListParams{
		FolderID: folderID,
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("pageSize", 20),
	}

	result, err := h.bookmarkService.GetByUser(c.Context(), userID, params)
	if err != nil {
		h.logger.Error("Failed to get bookmarks", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bookmarks",
		})
	}

	return c.JSON(result)
}

type AddBookmarkRequest struct {
	ArticleID string  `json:"articleId" validate:"required"`
	FolderID  *string `json:"folderId,omitempty"`
}

// Add adds an article to bookmarks
func (h *BookmarkHandler) Add(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req AddBookmarkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	articleID, err := uuid.Parse(req.ArticleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	var folderID *uuid.UUID
	if req.FolderID != nil {
		fid, err := uuid.Parse(*req.FolderID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid folder ID",
			})
		}
		folderID = &fid
	}

	if err := h.bookmarkService.Add(c.Context(), userID, articleID, folderID); err != nil {
		h.logger.Error("Failed to add bookmark", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add bookmark",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Bookmark added",
		"isBookmarked": true,
	})
}

// Remove removes an article from bookmarks
func (h *BookmarkHandler) Remove(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	articleIDStr := c.Params("articleId")
	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	if err := h.bookmarkService.Remove(c.Context(), userID, articleID); err != nil {
		h.logger.Error("Failed to remove bookmark", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove bookmark",
		})
	}

	return c.JSON(fiber.Map{
		"message":      "Bookmark removed",
		"isBookmarked": false,
	})
}

// Check checks if an article is bookmarked
func (h *BookmarkHandler) Check(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	articleIDStr := c.Params("articleId")
	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	isBookmarked, err := h.bookmarkService.IsBookmarked(c.Context(), userID, articleID)
	if err != nil {
		h.logger.Error("Failed to check bookmark", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check bookmark",
		})
	}

	return c.JSON(fiber.Map{
		"isBookmarked": isBookmarked,
	})
}

// ListFolders returns user's bookmark folders
func (h *BookmarkHandler) ListFolders(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	folders, err := h.bookmarkService.GetFolders(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get folders", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch folders",
		})
	}

	return c.JSON(fiber.Map{
		"items": folders,
	})
}

type CreateFolderRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// CreateFolder creates a new bookmark folder
func (h *BookmarkHandler) CreateFolder(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	folder, err := h.bookmarkService.CreateFolder(c.Context(), userID, req.Name)
	if err != nil {
		h.logger.Error("Failed to create folder", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create folder",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(folder)
}

type UpdateFolderRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// UpdateFolder updates a bookmark folder
func (h *BookmarkHandler) UpdateFolder(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid folder ID",
		})
	}

	var req UpdateFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.bookmarkService.UpdateFolder(c.Context(), userID, folderID, req.Name); err != nil {
		h.logger.Error("Failed to update folder", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update folder",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Folder updated",
	})
}

// DeleteFolder deletes a bookmark folder
func (h *BookmarkHandler) DeleteFolder(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	folderIDStr := c.Params("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid folder ID",
		})
	}

	if err := h.bookmarkService.DeleteFolder(c.Context(), userID, folderID); err != nil {
		h.logger.Error("Failed to delete folder", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete folder",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Folder deleted",
	})
}

type MoveBookmarkRequest struct {
	ArticleID string  `json:"articleId" validate:"required"`
	FolderID  *string `json:"folderId"`
}

// MoveToFolder moves a bookmark to a folder
func (h *BookmarkHandler) MoveToFolder(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req MoveBookmarkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	articleID, err := uuid.Parse(req.ArticleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	var folderID *uuid.UUID
	if req.FolderID != nil && *req.FolderID != "" {
		fid, err := uuid.Parse(*req.FolderID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid folder ID",
			})
		}
		folderID = &fid
	}

	if err := h.bookmarkService.MoveToFolder(c.Context(), userID, articleID, folderID); err != nil {
		h.logger.Error("Failed to move bookmark", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to move bookmark",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Bookmark moved",
	})
}

