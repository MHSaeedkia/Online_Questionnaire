package handlers

import (
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OptionHandler struct {
	optionRepo   repositories.OptionRepository
	questionRepo repositories.QuestionRepository
}

func NewOptionHandler(optionRepo repositories.OptionRepository, questionRepo repositories.QuestionRepository) *OptionHandler {
	return &OptionHandler{optionRepo, questionRepo}
}

type CreateOptionsRequest struct {
	QuestionID uint            `json:"question_id"`
	Options    []OptionRequest `json:"options"`
}

type OptionRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// CreateOptions handles the creation of options for a specific question
func (h *OptionHandler) CreateOptions(c *fiber.Ctx) error {
	questionnaireID, err := strconv.Atoi(c.Query("questionnaireId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}
	questionID, err := strconv.Atoi(c.Query("questionId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	var req CreateOptionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate options
	if len(req.Options) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Options cannot be empty"})
	}

	// Retrieve the question to check its type
	question, err := h.questionRepo.GetQuestion(uint(questionnaireID), uint(questionID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
	}

	// Ensure the question is a MultipleChoice type
	if question.Type != models.MultipleChoice {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Options can only be added to multiple-choice questions"})
	}

	// Prepare options for insertion
	var options []models.Option
	for _, option := range req.Options {
		options = append(options, models.Option{
			QuestionID: question.ID,
			Text:       option.Text,
			IsCorrect:  option.IsCorrect,
		})
	}

	// Store options in the database
	if err := h.optionRepo.CreateOptions(options); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create options"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Options created successfully"})
}
