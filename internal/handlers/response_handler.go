package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"strconv"
)

type ResponseHandler struct {
	responseRepo repositories.ResponseRepository
}

func NewResponseHandler(responseRepo repositories.ResponseRepository) *ResponseHandler {
	return &ResponseHandler{responseRepo: responseRepo}
}

// FillQuestionnaire handles user responses to a questionnaire.
func (h *ResponseHandler) FillQuestionnaire(c *fiber.Ctx) error {
	// Parse questionnaire ID from the URL
	questionnaireID, err := strconv.Atoi(c.Params("questionnaire_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	//// Parse user ID from the request body/JWT/session
	//userID, err := strconv.Atoi(c.FormValue("user_id"))
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	//}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	// Parse responses from the request body
	var responses []models.Response
	if err := c.BodyParser(&responses); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Attach the user and questionnaire IDs to each response
	for i := range responses {
		responses[i].UserID = uint(userID)
		responses[i].QuestionnaireID = uint(questionnaireID)
	}

	// Save responses
	for _, response := range responses {
		if err := h.responseRepo.CreateResponse(&response); err != nil {
			log.Printf("Error saving response: %v, Response: %+v", err, response)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save responses"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Responses submitted successfully"})
}

// EditResponse allows users to edit their responses to a questionnaire.
func (h *ResponseHandler) EditResponse(c *fiber.Ctx) error {
	// Parse questionnaire ID from the URL
	questionnaireID, err := strconv.Atoi(c.Params("questionnaire_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	// Parse user ID from the request body (or extract from JWT/session as needed)
	userID, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Parse the updated responses from the request body
	var updatedResponses []models.Response
	if err := c.BodyParser(&updatedResponses); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Loop through and update each response
	for _, updatedResponse := range updatedResponses {
		// Validate ownership of the response
		existingResponse := &models.Response{}
		if err := h.responseRepo.GetResponseByID(updatedResponse.ID, existingResponse); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Response not found"})
		}

		if existingResponse.UserID != uint(userID) || existingResponse.QuestionnaireID != uint(questionnaireID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to edit this response"})
		}

		// Update the answer
		existingResponse.Answer = updatedResponse.Answer
		if err := h.responseRepo.UpdateResponse(existingResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update response"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Responses updated successfully"})
}
