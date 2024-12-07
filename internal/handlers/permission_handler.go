package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"strconv"
)

type GrantPermissionRequest struct {
	Permission uint    `json:"permission_id"`
	ExpiresAt  *string `json:"expires_at"`
}

type PermissionHandler struct {
	questionnaireRepo repositories.QuestionnaireRepository
	permissionRepo    repositories.PermissionRepository
}

func NewPermissionHandler(qRepo repositories.QuestionnaireRepository, pRepo repositories.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{
		questionnaireRepo: qRepo,
		permissionRepo:    pRepo,
	}
}

func (h *PermissionHandler) RequestPermission(c *fiber.Ctx) error {
	questionnaireID, err := strconv.Atoi(c.Params("questionnaireID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	requestingUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	var req GrantPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	permission := models.QuestionnairePermission{
		QuestionnaireID: uint(questionnaireID),
		UserID:          requestingUserID,
		PermissionID:    req.Permission,
		Status:          "Pending",
	}
	err = h.permissionRepo.CreateQuestionnairePermission(&permission)
	if err != nil {
		log.Println("Error creating permission request:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to request permission"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Permission request submitted"})
}

func (h *PermissionHandler) ApproveOrDenyPermissionRequest(c *fiber.Ctx) error {
	requestID, err := strconv.Atoi(c.Params("requestID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request ID"})
	}

	requestingUserID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	permission, err := h.permissionRepo.GetQuestionnairePermission(uint(requestID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Request not found"})
	}

	questionnaire, err := h.questionnaireRepo.GetByID(permission.QuestionnaireID)
	if err != nil || questionnaire.OwnerID != requestingUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized"})
	}

	permission.Status = req.Status
	if req.Status == "Approved" {
		permission.ExpiresAt = nil // or set expiration
	}
	err = h.permissionRepo.UpdateQuestionnairePermission(&permission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update permission"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Request processed"})
}

//package handlers
//
//import (
//	"github.com/gofiber/fiber/v2"
//	"log"
//	"online-questionnaire/internal/models"
//	"online-questionnaire/internal/repositories"
//	"strconv"
//)
//
//type GrantPermissionRequest struct {
//	UserID     uint        `json:"user_id"`
//	Permission models.Name `json:"permission"`
//	ExpiresAt  *string     `json:"expires_at"`
//}
//
//type PermissionHandler struct {
//	questionnaireRepo repositories.QuestionnaireRepository
//	permissionRepo    repositories.PermissionRepository
//}
//
//func NewPermissionHandler(qRepo repositories.QuestionnaireRepository, pRepo repositories.PermissionRepository) *PermissionHandler {
//	return &PermissionHandler{
//		questionnaireRepo: qRepo,
//		permissionRepo:    pRepo,
//	}
//}
//
//func (h *PermissionHandler) GrantPermissionToUser(c *fiber.Ctx) error {
//	questionnaireID, err := strconv.Atoi(c.Params("questionnaireID"))
//	if err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
//	}
//
//	var req GrantPermissionRequest
//	if err := c.BodyParser(&req); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
//	}
//
//	// Validate if the requesting user is the owner of the questionnaire
//	requestingUserID, ok := c.Locals("user_id").(uint) // Assuming user_id is set in middleware
//	if !ok {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
//	}
//
//	questionnaire, err := h.questionnaireRepo.GetByID(uint(questionnaireID))
//	if err != nil {
//		log.Println("Error fetching questionnaire:", err)
//		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Questionnaire not found"})
//	}
//
//	if questionnaire.OwnerID != requestingUserID {
//		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to grant permissions for this questionnaire"})
//	}
//
//	// Add the permission to the database
//	err = h.permissionRepo.GrantPermission(uint(questionnaireID), req.UserID, req.Permission, req.ExpiresAt)
//	if err != nil {
//		log.Println("Error granting permission:", err)
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to grant permission"})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Permission granted successfully"})
//}
