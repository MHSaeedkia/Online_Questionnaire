package handlers

import (
	"fmt"
	"online-questionnaire/internal/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *UserHandler {
	return &UserHandler{repo}
}

func (h *UserHandler) Quesionnare(c *fiber.Ctx) error {
	ownerId, err := strconv.Atoi(c.Query("ownerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	quesionnare, err := h.repo.Quesionnare(uint(ownerId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch questionnaires by owner id : %v", ownerId)})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Questionnaire fetch successfully", "data": quesionnare})
}
