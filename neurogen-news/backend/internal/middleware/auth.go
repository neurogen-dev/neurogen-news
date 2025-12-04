package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neurogen-news/backend/internal/model"
	"github.com/neurogen-news/backend/internal/service"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserKey     contextKey = "user"
	UserRoleKey contextKey = "userRole"
)

func Auth(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(c.Context(), token)
		if err != nil {
			if err == service.ErrTokenExpired {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token expired",
					"code":  "TOKEN_EXPIRED",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Set user info in context
		c.Locals(string(UserIDKey), claims.UserID)
		c.Locals(string(UserRoleKey), claims.Role)

		return c.Next()
	}
}

func OptionalAuth(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := authService.ValidateToken(c.Context(), token)
		if err == nil {
			c.Locals(string(UserIDKey), claims.UserID)
			c.Locals(string(UserRoleKey), claims.Role)
		}

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals(string(UserRoleKey)).(model.UserRole)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		// Check if user's role is in the allowed roles
		for _, r := range roles {
			if string(role) == r {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}
}

// Error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

