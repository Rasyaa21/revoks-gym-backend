package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	service services.AttendanceService
}

func NewAttendanceHandler(service services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) MyHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "30"))
	resp, err := h.service.MyHistory(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get attendance history", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Attendance history retrieved", resp))
}
