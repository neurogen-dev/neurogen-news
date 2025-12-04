package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type AchievementHandler struct {
	achievementService service.AchievementService
	logger             *zap.Logger
}

func NewAchievementHandler(achievementService service.AchievementService, logger *zap.Logger) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
		logger:             logger,
	}
}

// List returns all achievements
func (h *AchievementHandler) List(c *fiber.Ctx) error {
	achievements, err := h.achievementService.GetAll(c.Context())
	if err != nil {
		h.logger.Error("Failed to get achievements", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch achievements",
		})
	}

	return c.JSON(fiber.Map{
		"items": achievements,
	})
}

// GetByID returns an achievement by ID
func (h *AchievementHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid achievement ID",
		})
	}

	achievement, err := h.achievementService.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "achievement not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Achievement not found",
			})
		}
		h.logger.Error("Failed to get achievement", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch achievement",
		})
	}

	return c.JSON(achievement)
}

// GetUserAchievements returns achievements for a user
func (h *AchievementHandler) GetUserAchievements(c *fiber.Ctx) error {
	username := c.Params("username")

	// TODO: Get user by username and then get their achievements
	// For now, just return empty if no user found

	userIDStr := c.Query("userId")
	var userID uuid.UUID

	if userIDStr != "" {
		var err error
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}
	} else if username != "" {
		// Would need to look up user by username
		return c.JSON(fiber.Map{
			"items": []interface{}{},
		})
	}

	achievements, err := h.achievementService.GetUserAchievements(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user achievements", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch achievements",
		})
	}

	return c.JSON(fiber.Map{
		"items": achievements,
	})
}

// GetMyAchievements returns the authenticated user's achievements
func (h *AchievementHandler) GetMyAchievements(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	achievements, err := h.achievementService.GetUserAchievements(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get achievements", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch achievements",
		})
	}

	return c.JSON(fiber.Map{
		"items": achievements,
	})
}

// GetProgress returns the authenticated user's achievement progress
func (h *AchievementHandler) GetProgress(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	progress, err := h.achievementService.GetProgress(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get progress", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch progress",
		})
	}

	return c.JSON(progress)
}

// CheckAchievements checks and awards new achievements for the authenticated user
func (h *AchievementHandler) CheckAchievements(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	newAchievements, err := h.achievementService.CheckAndAward(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to check achievements", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check achievements",
		})
	}

	return c.JSON(fiber.Map{
		"newAchievements": newAchievements,
	})
}

