package permission_handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories/permission_repo"
	"online-questionnaire/internal/repositories/questionnaire_repo"
	"strconv"
)

type GrantPermissionRequest struct {
	Permission uint    `json:"permission_id"`
	ExpiresAt  *string `json:"expires_at"`
}

type PermissionHandler struct {
	questionnaireRepo questionnaire_repo.QuestionnaireRepository
	permissionRepo    permission_repo.PermissionRepository
}

func NewPermissionHandler(qRepo questionnaire_repo.QuestionnaireRepository, pRepo permission_repo.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{
		questionnaireRepo: qRepo,
		permissionRepo:    pRepo,
	}
}

func (h *PermissionHandler) RequestPermission(c *fiber.Ctx) error {
	// Takes the id of requested questionnaire
	questionnaireID, err := strconv.Atoi(c.Params("questionnaireID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
	}

	// Checks if the person requesting is authenticated or not
	requestingUserID, ok := c.Locals("user_id").(uint)
	fmt.Println("requestingUserID=", requestingUserID)
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
	fmt.Println("requestID=", requestID)
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
