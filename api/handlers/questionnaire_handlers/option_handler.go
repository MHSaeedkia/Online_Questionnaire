package questionnaire_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories/questionnaire_repo"
	"strconv"
)

type OptionHandler struct {
	optionRepo   questionnaire_repo.OptionRepository
	questionRepo questionnaire_repo.QuestionRepository
}

func NewOptionHandler(optionRepo questionnaire_repo.OptionRepository, questionRepo questionnaire_repo.QuestionRepository) *OptionHandler {
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
	questionnaireID := c.Params("questionnaire_id")
	questionID := c.Params("question_id")

	var req CreateOptionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate options
	if len(req.Options) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Options cannot be empty"})
	}

	quesID, err := strconv.ParseUint(questionnaireID, 10, 32) // Base 10, 32-bit unsigned integer
	if err != nil {
		fmt.Println("Error:", err)
	}

	q, err := strconv.ParseUint(questionID, 10, 32) // Base 10, 32-bit unsigned integer
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Cast the uint64 to uint
	qID := uint(quesID)
	ask := uint(q)

	// Retrieve the question to check its type
	question, err := h.questionRepo.GetQuestionByID(qID, ask)
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
