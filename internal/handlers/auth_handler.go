package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Login failed", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Login success", resp))
}
