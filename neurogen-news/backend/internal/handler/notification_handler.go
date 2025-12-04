package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type NotificationHandler struct {
	notificationService service.NotificationService
	logger              *zap.Logger
}

func NewNotificationHandler(notificationService service.NotificationService, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger,
	}
}

// List returns user's notifications
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	result, err := h.notificationService.GetByUser(c.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get notifications", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	return c.JSON(result)
}

// GetUnreadCount returns the count of unread notifications
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	count, err := h.notificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get unread count", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
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
			"error": "Invalid notification ID",
		})
	}

	if err := h.notificationService.MarkAsRead(c.Context(), userID, id); err != nil {
		h.logger.Error("Failed to mark as read", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := h.notificationService.MarkAllAsRead(c.Context(), userID); err != nil {
		h.logger.Error("Failed to mark all as read", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark all notifications as read",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}

// Delete deletes a notification
func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
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
			"error": "Invalid notification ID",
		})
	}

	if err := h.notificationService.Delete(c.Context(), userID, id); err != nil {
		h.logger.Error("Failed to delete notification", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notification",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification deleted",
	})
}

// DeleteAll deletes all notifications
func (h *NotificationHandler) DeleteAll(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := h.notificationService.DeleteAll(c.Context(), userID); err != nil {
		h.logger.Error("Failed to delete all notifications", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notifications",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications deleted",
	})
}

