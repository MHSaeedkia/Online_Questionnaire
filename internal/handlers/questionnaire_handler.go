package handlers

import (
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
)

type QuestionnaireHandler struct {
	repo repositories.QuestionnaireRepository
}

func NewQuestionnaireHandler(repo repositories.QuestionnaireRepository) *QuestionnaireHandler {
	return &QuestionnaireHandler{repo}
}

func (h *QuestionnaireHandler) CreateQuestionnaire(c *fiber.Ctx) error {
	var req models.Questionnaire
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Set additional fields
	req.CreationTime = time.Now()

	// Call repository to create
	if err := h.repo.CreateQuestionnaire(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create questionnaire"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Questionnaire created successfully", "data": req})
}
