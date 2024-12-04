package handlers

import (
	"github.com/gofiber/fiber/v2"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
)

type ConditionalLogicHandler struct {
	repo         repositories.ConditionalLogicRepository
	questionRepo repositories.QuestionRepository
	optionRepo   repositories.OptionRepository
}

func NewConditionalLogicHandler(repo repositories.ConditionalLogicRepository, questionRepo repositories.QuestionRepository, optionRepo repositories.OptionRepository) *ConditionalLogicHandler {
	return &ConditionalLogicHandler{repo, questionRepo, optionRepo}
}

type CreateConditionalLogicRequest struct {
	QuestionnaireID  uint `json:"questionnaire_id"`
	QuestionID       uint `json:"question_id"`
	OptionID         uint `json:"option_id"`
	TargetQuestionID uint `json:"target_question_id"`
}

// // CreateConditionalLogic handles the creation of conditional logic
//
//	func (h *ConditionalLogicHandler) CreateConditionalLogic(c *fiber.Ctx) error {
//		var req CreateConditionalLogicRequest
//
//		if err := c.BodyParser(&req); err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
//		}
//
//		// Validate if Questionnaire exists
//		_, err := h.questionRepo.GetQuestionnaireByID(
//			req.QuestionnaireID,
//		)
//		if err != nil {
//			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Questionnaire not found"})
//		}
//
//		// Validate if Question exists
//		_, e := h.questionRepo.GetQuestionByID(
//			req.QuestionnaireID,
//			req.QuestionID,
//		)
//		if e != nil {
//			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
//		}
//
//		// Validate if the Option exists
//		option, err := h.optionRepo.GetOptionByID(req.OptionID)
//		if err != nil || option.QuestionID != req.QuestionID {
//			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Option not found or does not belong to the specified question"})
//		}
//
//		// Validate if the Target Question exists
//		_, err = h.questionRepo.GetQuestionByID(
//			req.QuestionnaireID,
//			req.TargetQuestionID,
//		)
//		if err != nil {
//			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Target question not found"})
//		}
//
//		// Create Conditional Logic
//		logic := &models.ConditionalLogic{
//			QuestionID:       req.QuestionID,
//			OptionID:         req.OptionID,
//			TargetQuestionID: req.TargetQuestionID,
//		}
//
//		if err := h.repo.CreateConditionalLogic(logic); err != nil {
//			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create conditional logic"})
//		}
//
//		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Conditional logic created successfully"})
//	}
func (h *ConditionalLogicHandler) CreateConditionalLogic(c *fiber.Ctx) error {
	// Extract `questionnaire_id` and `question_id` from the URL
	questionnaireID, err := c.ParamsInt("questionnaire_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	questionID, err := c.ParamsInt("question_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid question ID"})
	}

	var req CreateConditionalLogicRequest
	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Use extracted path parameters instead of relying on body fields
	req.QuestionnaireID = uint(questionnaireID)
	req.QuestionID = uint(questionID)

	// Validate if Questionnaire exists
	_, err = h.questionRepo.GetQuestionnaireByID(req.QuestionnaireID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Questionnaire not found"})
	}

	// Validate if Question exists
	_, err = h.questionRepo.GetQuestionByID(req.QuestionnaireID, req.QuestionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Question not found"})
	}

	// Validate if Option exists
	option, err := h.optionRepo.GetOptionByID(req.OptionID)
	if err != nil || option.QuestionID != req.QuestionID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Option not found or does not belong to the specified question"})
	}

	// Validate Target Question
	_, err = h.questionRepo.GetQuestionByID(req.QuestionnaireID, req.TargetQuestionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Target question not found"})
	}

	// Create the Conditional Logic
	logic := &models.ConditionalLogic{
		QuestionID:       req.QuestionID,
		OptionID:         req.OptionID,
		TargetQuestionID: req.TargetQuestionID,
	}

	if err := h.repo.CreateConditionalLogic(logic); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create conditional logic"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Conditional logic created successfully"})
}
