package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type WorkoutHandler struct {
	service services.WorkoutService
}

func NewWorkoutHandler(service services.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{service: service}
}

func (h *WorkoutHandler) Create(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	var req dto.CreateWorkoutProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Validation error", "title is required"))
	}
	resp, err := h.service.Create(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to create workout progress", err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(dto.NewSuccessResponse("Workout progress created", resp))
}

func (h *WorkoutHandler) MyHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	limit, _ := strconv.Atoi(c.Query("limit", "30"))
	resp, err := h.service.MyHistory(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get workout history", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Workout progress history retrieved", resp))
}
