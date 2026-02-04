package middleware

import (
	"errors"
	"strings"

	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const userIDLocalKey = "userId"

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", "missing Authorization header"))
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", "invalid Authorization header"))
		}

		claims, err := utils.ParseToken(parts[1], "access")
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", "invalid or expired token"))
		}

		c.Locals(userIDLocalKey, claims.UserID)
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (uint, error) {
	v := c.Locals(userIDLocalKey)
	if v == nil {
		return 0, errors.New("missing user context")
	}
	userID, ok := v.(uint)
	if !ok {
		return 0, errors.New("invalid user context")
	}
	return userID, nil
}
