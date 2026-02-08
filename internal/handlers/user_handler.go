package handlers

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to create user", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewSuccessResponse("User created successfully", user))
}

// GetAllUsers handles GET /api/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	users, pagination, err := h.service.GetAllUsers(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Failed to get users", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponseWithPagination("Users retrieved successfully", users, pagination))
}

// GetUserByID handles GET /api/users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid user ID", err.Error()))
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("User not found", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("User retrieved successfully", user))
}

// UpdateUser handles PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid user ID", err.Error()))
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to update user", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("User updated successfully", user))
}

// DeleteUser handles DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid user ID", err.Error()))
	}

	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("Failed to delete user", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("User deleted successfully", nil))
}

// ChangePassword handles PUT /api/users/:id/password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid user ID", err.Error()))
	}

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.service.ChangePassword(id, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Failed to change password", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse("Password changed successfully", nil))
}
