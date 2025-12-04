package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile returns a user's public profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		h.logger.Error("Failed to get user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	profile, err := h.userService.GetProfile(c.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get profile", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch profile",
		})
	}

	// Check if current user is following
	if currentUserID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID); ok {
		isFollowing, _ := h.userService.IsFollowing(c.Context(), currentUserID, user.ID)
		profile.IsFollowing = isFollowing
	}

	return c.JSON(profile)
}

// GetCurrentProfile returns the authenticated user's profile
func (h *UserHandler) GetCurrentProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	profile, err := h.userService.GetProfile(c.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get profile", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch profile",
		})
	}

	return c.JSON(profile)
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	CoverURL    *string `json:"coverUrl,omitempty"`
	Location    *string `json:"location,omitempty"`
	Website     *string `json:"website,omitempty"`
	Telegram    *string `json:"telegram,omitempty"`
	Github      *string `json:"github,omitempty"`
}

// UpdateProfile updates the authenticated user's profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.userService.UpdateProfile(c.Context(), userID, service.UpdateProfileInput{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
		CoverURL:    req.CoverURL,
		Location:    req.Location,
		Website:     req.Website,
		Telegram:    req.Telegram,
		Github:      req.Github,
	})
	if err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(user)
}

// Follow follows a user
func (h *UserHandler) Follow(c *fiber.Ctx) error {
	followerID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	username := c.Params("username")
	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if err := h.userService.Follow(c.Context(), followerID, user.ID); err != nil {
		if err.Error() == "Cannot follow yourself" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot follow yourself",
			})
		}
		h.logger.Error("Failed to follow", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to follow user",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Followed successfully",
		"isFollowing": true,
	})
}

// Unfollow unfollows a user
func (h *UserHandler) Unfollow(c *fiber.Ctx) error {
	followerID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	username := c.Params("username")
	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if err := h.userService.Unfollow(c.Context(), followerID, user.ID); err != nil {
		h.logger.Error("Failed to unfollow", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unfollow user",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Unfollowed successfully",
		"isFollowing": false,
	})
}

// GetFollowers returns a user's followers
func (h *UserHandler) GetFollowers(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	followers, total, err := h.userService.GetFollowers(c.Context(), user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get followers", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch followers",
		})
	}

	return c.JSON(fiber.Map{
		"items":   followers,
		"total":   total,
		"hasMore": offset+len(followers) < total,
	})
}

// GetFollowing returns users that a user is following
func (h *UserHandler) GetFollowing(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	following, total, err := h.userService.GetFollowing(c.Context(), user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get following", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch following",
		})
	}

	return c.JSON(fiber.Map{
		"items":   following,
		"total":   total,
		"hasMore": offset+len(following) < total,
	})
}

// GetArticles returns a user's articles
func (h *UserHandler) GetArticles(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.userService.GetByUsername(c.Context(), username)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	articles, total, err := h.userService.GetUserArticles(c.Context(), user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get articles", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch articles",
		})
	}

	return c.JSON(fiber.Map{
		"items":   articles,
		"total":   total,
		"hasMore": offset+len(articles) < total,
	})
}

