package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

type QRHandler struct {
	service services.QRService
}

func NewQRHandler(service services.QRService) *QRHandler {
	return &QRHandler{service: service}
}

func (h *QRHandler) GetMyQRCode(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	resp, err := h.service.GetMyQRCode(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get QR", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("QR generated successfully", resp))
}

// Scan is intended for gate/QR scanner device
func (h *QRHandler) Scan(c *fiber.Ctx) error {
	var req dto.ScanQRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}
	resp, err := h.service.Scan(req.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to scan QR", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Scan processed", resp))
}
