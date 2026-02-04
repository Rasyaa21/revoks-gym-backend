package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(service services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) MyList(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	items, err := h.service.MyList(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get notifications", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Notifications retrieved", items))
}

func (h *NotificationHandler) Detail(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid notification ID", err.Error()))
	}
	item, err := h.service.Detail(userID, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("Notification not found", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Notification detail retrieved", item))
}

func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid notification ID", err.Error()))
	}
	if err := h.service.MarkRead(userID, uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to mark read", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Notification marked as read", nil))
}
