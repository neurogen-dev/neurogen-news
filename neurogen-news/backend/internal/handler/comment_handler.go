package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
	logger         *zap.Logger
}

func NewCommentHandler(commentService service.CommentService, logger *zap.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}

// GetByArticle returns comments for a specific article
func (h *CommentHandler) GetByArticle(c *fiber.Ctx) error {
	articleIDStr := c.Params("articleId")
	articleID, err := uuid.Parse(articleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid article ID",
		})
	}

	params := service.CommentListParams{
		Sort:     c.Query("sort", "new"),
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("pageSize", 20),
	}

	result, err := h.commentService.GetByArticle(c.Context(), articleID, params)
	if err != nil {
		h.logger.Error("Failed to get comments", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch comments",
		})
	}

	return c.JSON(result)
}

// GetReplies returns replies for a specific comment
func (h *CommentHandler) GetReplies(c *fiber.Ctx) error {
	parentIDStr := c.Params("id")
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	replies, err := h.commentService.GetReplies(c.Context(), parentID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get replies", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch replies",
		})
	}

	return c.JSON(fiber.Map{
		"items": replies,
	})
}

// GetByID returns a single comment by ID
func (h *CommentHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	comment, err := h.commentService.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "comment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		h.logger.Error("Failed to get comment", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch comment",
		})
	}

	return c.JSON(comment)
}

type CreateCommentRequest struct {
	ArticleID string  `json:"articleId" validate:"required"`
	ParentID  *string `json:"parentId,omitempty"`
	Content   string  `json:"content" validate:"required,min=1,max=10000"`
}

// Create creates a new comment
func (h *CommentHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateCommentRequest
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

	var parentID *uuid.UUID
	if req.ParentID != nil {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid parent ID",
			})
		}
		parentID = &pid
	}

	comment, err := h.commentService.Create(c.Context(), userID, service.CreateCommentInput{
		ArticleID: articleID,
		ParentID:  parentID,
		Content:   req.Content,
	})
	if err != nil {
		h.logger.Error("Failed to create comment", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create comment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=10000"`
}

// Update updates an existing comment
func (h *CommentHandler) Update(c *fiber.Ctx) error {
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
			"error": "Invalid comment ID",
		})
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	comment, err := h.commentService.Update(c.Context(), userID, id, req.Content)
	if err != nil {
		if err.Error() == "comment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to edit this comment",
			})
		}
		h.logger.Error("Failed to update comment", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update comment",
		})
	}

	return c.JSON(comment)
}

// Delete deletes a comment
func (h *CommentHandler) Delete(c *fiber.Ctx) error {
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
			"error": "Invalid comment ID",
		})
	}

	if err := h.commentService.Delete(c.Context(), userID, id); err != nil {
		if err.Error() == "comment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		if err.Error() == "FORBIDDEN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to delete this comment",
			})
		}
		h.logger.Error("Failed to delete comment", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete comment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}

type AddReactionRequest struct {
	Emoji string `json:"emoji" validate:"required"`
}

// AddReaction adds a reaction to a comment
func (h *CommentHandler) AddReaction(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	var req AddReactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.commentService.AddReaction(c.Context(), userID, commentID, req.Emoji); err != nil {
		h.logger.Error("Failed to add reaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add reaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reaction added",
	})
}

// RemoveReaction removes a reaction from a comment
func (h *CommentHandler) RemoveReaction(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	idStr := c.Params("id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	if err := h.commentService.RemoveReaction(c.Context(), userID, commentID); err != nil {
		h.logger.Error("Failed to remove reaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove reaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reaction removed",
	})
}

