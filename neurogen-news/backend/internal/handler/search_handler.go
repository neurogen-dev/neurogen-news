package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/service"
)

type SearchHandler struct {
	searchService service.SearchService
	logger        *zap.Logger
}

func NewSearchHandler(searchService service.SearchService, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		logger:        logger,
	}
}

// Search performs a search
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	params := service.SearchParams{
		Type:     c.Query("type", "all"),
		Sort:     c.Query("sort", "relevance"),
		Level:    c.Query("level"),
		Category: c.Query("category"),
		Limit:    c.QueryInt("limit", 20),
		Offset:   c.QueryInt("offset", 0),
	}

	result, err := h.searchService.Search(c.Context(), query, params)
	if err != nil {
		h.logger.Error("Search failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search failed",
		})
	}

	return c.JSON(result)
}

// GetSuggestions returns search suggestions
func (h *SearchHandler) GetSuggestions(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.JSON(fiber.Map{
			"items": []interface{}{},
		})
	}

	limit := c.QueryInt("limit", 5)

	suggestions, err := h.searchService.GetSuggestions(c.Context(), query, limit)
	if err != nil {
		h.logger.Error("Failed to get suggestions", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get suggestions",
		})
	}

	return c.JSON(fiber.Map{
		"items": suggestions,
	})
}

