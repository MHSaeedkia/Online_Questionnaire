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

func (h *QuestionnaireHandler) CreateQuestion(c *fiber.Ctx) error {
	var req models.Question
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Call repository to create
	if err := h.repo.CreateQuestion(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create question"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Question created successfully", "data": req})
}

func (h *QuestionnaireHandler) CreateAnswer(c *fiber.Ctx) error {
	var req models.Response
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Call repository to create
	if err := h.repo.CreateAnswer(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create answer"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Answer created successfully", "data": req})
}
