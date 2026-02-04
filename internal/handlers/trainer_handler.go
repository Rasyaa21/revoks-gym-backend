package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TrainerHandler struct {
	service services.TrainerService
}

func NewTrainerHandler(service services.TrainerService) *TrainerHandler {
	return &TrainerHandler{service: service}
}

func (h *TrainerHandler) List(c *fiber.Ctx) error {
	items, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get PT list", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("PT list retrieved", items))
}

func (h *TrainerHandler) Detail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid PT ID", err.Error()))
	}
	item, err := h.service.Detail(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("PT not found", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("PT detail retrieved", item))
}

func (h *TrainerHandler) Schedule(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid PT ID", err.Error()))
	}
	items, err := h.service.Schedule(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get PT schedule", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("PT schedule retrieved", items))
}
