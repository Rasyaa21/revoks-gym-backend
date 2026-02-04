package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TemplateHandler struct {
	service services.TemplateService
}

func NewTemplateHandler(service services.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (h *TemplateHandler) MyFollowed(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	resp, err := h.service.MyFollowed(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get templates", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Followed templates retrieved", resp))
}

func (h *TemplateHandler) Detail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid template ID", err.Error()))
	}
	resp, err := h.service.Detail(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("Template not found", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Template detail retrieved", resp))
}

func (h *TemplateHandler) Follow(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	var req dto.FollowTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}
	if req.TemplateID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Validation error", "template_id is required"))
	}
	if err := h.service.Follow(userID, req.TemplateID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to follow template", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Template followed", nil))
}
