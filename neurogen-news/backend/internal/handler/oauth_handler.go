package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/neurogen-news/backend/internal/middleware"
	"github.com/neurogen-news/backend/internal/service"
)

type OAuthHandler struct {
	oauthService service.OAuthService
	logger       *zap.Logger
}

func NewOAuthHandler(oauthService service.OAuthService, logger *zap.Logger) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		logger:       logger,
	}
}

// GetAuthURL returns the OAuth authorization URL for a provider
func (h *OAuthHandler) GetAuthURL(c *fiber.Ctx) error {
	provider := c.Params("provider")
	if provider == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provider is required",
		})
	}

	// Generate state for CSRF protection
	state := uuid.New().String()

	// Store state in session/cookie for verification
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   300, // 5 minutes
	})

	url, err := h.oauthService.GetAuthURL(service.OAuthProvider(provider), state)
	if err != nil {
		h.logger.Error("Failed to get auth URL", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url":   url,
		"state": state,
	})
}

// Callback handles OAuth callback from provider
func (h *OAuthHandler) Callback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	// Verify state
	storedState := c.Cookies("oauth_state")
	if storedState == "" || storedState != state {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid state parameter",
		})
	}

	// Clear state cookie
	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
	})

	user, tokens, err := h.oauthService.HandleCallback(c.Context(), service.OAuthProvider(provider), code)
	if err != nil {
		h.logger.Error("OAuth callback failed",
			zap.String("provider", provider),
			zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	return c.JSON(fiber.Map{
		"user":         user,
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expiresIn":    tokens.ExpiresIn,
	})
}

// CallbackRedirect handles OAuth callback with redirect to frontend
func (h *OAuthHandler) CallbackRedirect(c *fiber.Ctx) error {
	provider := c.Params("provider")
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Handle error from provider
	if errorParam != "" {
		frontendURL := c.Query("redirect_uri", "/login")
		return c.Redirect(frontendURL + "?error=" + errorParam)
	}

	if code == "" {
		return c.Redirect("/login?error=no_code")
	}

	// Verify state
	storedState := c.Cookies("oauth_state")
	if storedState == "" || storedState != state {
		return c.Redirect("/login?error=invalid_state")
	}

	// Clear state cookie
	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
	})

	_, tokens, err := h.oauthService.HandleCallback(c.Context(), service.OAuthProvider(provider), code)
	if err != nil {
		h.logger.Error("OAuth callback failed",
			zap.String("provider", provider),
			zap.Error(err))
		return c.Redirect("/login?error=auth_failed")
	}

	// Redirect to frontend with tokens
	frontendURL := "/auth/callback"
	return c.Redirect(frontendURL + "?access_token=" + tokens.AccessToken + "&refresh_token=" + tokens.RefreshToken)
}

type TelegramAuthRequest struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  string `json:"auth_date"`
	Hash      string `json:"hash"`
}

// TelegramAuth handles Telegram authentication
func (h *OAuthHandler) TelegramAuth(c *fiber.Ctx) error {
	var req TelegramAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Convert to map for verification
	data := map[string]string{
		"id":        req.ID,
		"auth_date": req.AuthDate,
		"hash":      req.Hash,
	}
	if req.FirstName != "" {
		data["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		data["last_name"] = req.LastName
	}
	if req.Username != "" {
		data["username"] = req.Username
	}
	if req.PhotoURL != "" {
		data["photo_url"] = req.PhotoURL
	}

	oauthUser, err := h.oauthService.VerifyTelegramAuth(data)
	if err != nil {
		h.logger.Error("Telegram auth verification failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Telegram authentication",
		})
	}

	// Handle callback with verified user
	user, tokens, err := h.oauthService.HandleCallback(c.Context(), service.ProviderTelegram, oauthUser.ProviderID)
	if err != nil {
		h.logger.Error("Failed to handle Telegram auth", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	return c.JSON(fiber.Map{
		"user":         user,
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expiresIn":    tokens.ExpiresIn,
	})
}

// LinkProvider links OAuth provider to current user
func (h *OAuthHandler) LinkProvider(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	provider := c.Params("provider")
	code := c.Query("code")

	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	if err := h.oauthService.LinkProvider(c.Context(), userID, service.OAuthProvider(provider), code); err != nil {
		h.logger.Error("Failed to link provider", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link provider",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Provider linked successfully",
	})
}

// UnlinkProvider unlinks OAuth provider from current user
func (h *OAuthHandler) UnlinkProvider(c *fiber.Ctx) error {
	userID, ok := c.Locals(string(middleware.UserIDKey)).(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	provider := c.Params("provider")

	if err := h.oauthService.UnlinkProvider(c.Context(), userID, service.OAuthProvider(provider)); err != nil {
		h.logger.Error("Failed to unlink provider", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unlink provider",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Provider unlinked successfully",
	})
}


