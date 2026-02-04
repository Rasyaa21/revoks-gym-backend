package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MembershipHandler struct {
	service services.MembershipService
}

func NewMembershipHandler(service services.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: service}
}

func (h *MembershipHandler) GetMyMembership(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	resp, err := h.service.GetMyMembership(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get membership", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Membership retrieved successfully", resp))
}

func (h *MembershipHandler) Renew(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized", err.Error()))
	}
	var req dto.RenewMembershipRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}
	resp, err := h.service.Renew(userID, req.Months, req.Plan)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to renew membership", err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Membership renewed successfully", resp))
}
