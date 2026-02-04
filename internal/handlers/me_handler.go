package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MeHandler struct {
	userService services.UserService
}

func NewMeHandler(userService services.UserService) *MeHandler {
	return &MeHandler{userService: userService}
}

func (h *MeHandler) GetMe(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("User not found", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Profile retrieved successfully", user))
}
