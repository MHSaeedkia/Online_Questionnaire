package handlers

import (
	"fmt"
	"online-questionnaire/internal/models"
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Questionnaire fetch successfully", "data": quesionnare})
}

func (h *UserHandler) EditQuestionnare(c *fiber.Ctx) error {
	ownerId, err := strconv.Atoi(c.Query("ownerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}
	quesionnareId, err := strconv.Atoi(c.Query("quesionnareId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	var request models.Questionnaire
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err = h.repo.EditQuestionnare(uint(ownerId), uint(quesionnareId), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch questionnaires by id : %v from owner id : %v", quesionnareId, ownerId)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Quesionnare by id : %v has updated", quesionnareId)})
}

func (h *UserHandler) CancleQuestionnarec(c *fiber.Ctx) error {
	quesionnareId, err := strconv.Atoi(c.Query("quesionnareId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	err = h.repo.CancleQuestionnare(uint(quesionnareId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch questionnaires by id : %v", quesionnareId)})

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Quesionnare by id : %v has deleted", quesionnareId)})

}
