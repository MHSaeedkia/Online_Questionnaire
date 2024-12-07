package handlers

import (
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"strconv"
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

func (h *QuestionnaireHandler) GetQuestion(c *fiber.Ctx) error {
	quetionId, err := strconv.Atoi(c.Query("quetionId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	resp, err := h.repo.GetQuestion(uint(quetionId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get question"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "", "data": resp})
}

func (h *QuestionnaireHandler) UpdateQuestion(c *fiber.Ctx) error {
	var req models.Question
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.repo.UpdateQuestion(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update question"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Question updated successfully", "data": req})
}

func (h *QuestionnaireHandler) DeleteQuestion(c *fiber.Ctx) error {
	quesionId, err := strconv.Atoi(c.Query("quesionId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.repo.DeleteQuestion(uint(quesionId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete question"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Question deleted successfully"})
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

func (h *QuestionnaireHandler) GetAnswer(c *fiber.Ctx) error {
	answerId, err := strconv.Atoi(c.Query("answerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad query param"})
	}

	resp, err := h.repo.GetAnswer(uint(answerId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get answer"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "", "data": resp})
}

func (h *QuestionnaireHandler) UpdateAnswer(c *fiber.Ctx) error {
	var req models.Response
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.repo.UpdateAnswer(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update answer"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Answer updated successfully", "data": req})
}

func (h *QuestionnaireHandler) DeleteAnswer(c *fiber.Ctx) error {
	answerId, err := strconv.Atoi(c.Query("answerId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.repo.DeleteAnswer(uint(answerId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete answer"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Answer deleted successfully"})
}
