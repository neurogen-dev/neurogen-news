package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type UploadHandler struct {
	uploadService service.UploadService
	logger        *zap.Logger
}

func NewUploadHandler(uploadService service.UploadService, logger *zap.Logger) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		logger:        logger,
	}
}

// UploadImage uploads an image
func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided",
		})
	}

	uploadType := c.FormValue("type", "article")
	if uploadType != "article" && uploadType != "avatar" && uploadType != "cover" {
		uploadType = "article"
	}

	result, err := h.uploadService.UploadImage(c.Context(), userID, file, uploadType)
	if err != nil {
		h.logger.Error("Upload failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

// DeleteFile deletes an uploaded file
func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req struct {
		URL string `json:"url" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.uploadService.DeleteFile(c.Context(), userID, req.URL); err != nil {
		h.logger.Error("Delete failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"message": "File deleted",
	})
}

