package handlers

import (
	"github.com/gofiber/fiber/v2"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"strconv"
)

type QuestionHandler struct {
	repo repositories.QuestionRepository
}

func NewQuestionHandler(repo repositories.QuestionRepository) *QuestionHandler {
	return &QuestionHandler{repo}
}

func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	questionnaireID, err := strconv.Atoi(c.Params("questionnaire_id"))
	if err != nil || questionnaireID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	// Validate that the questionnaire exists
	questionnaire, err := h.repo.GetQuestionnaireByID(uint(questionnaireID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Questionnaire not found"})
	}

	// Parse the question from the request body
	var req models.Question
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Assign the questionnaire ID and other default values
	req.QuestionnaireID = questionnaire.ID

	// Save the question to the database
	if err := h.repo.CreateQuestion(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create question"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Question created successfully",
		"data":    req,
	})
}
