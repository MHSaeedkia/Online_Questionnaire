package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
	"strconv"
)

func CheckPermission(db *gorm.DB, requiredPermission models.Name) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve user ID from context
		userID, ok := c.Locals("user_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		}

		// Retrieve questionnaire ID from URL parameters
		questionnaireID, err := strconv.Atoi(c.Params("questionnaire_id"))
		fmt.Println("questionnaireID=", questionnaireID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid questionnaire ID"})
		}

		// Check if the user has the required permission for the questionnaire
		var permission models.QuestionnairePermission
		err = db.Joins("JOIN permissions ON permissions.id = questionnaire_permissions.permission_id").
			Where("questionnaire_permissions.questionnaire_id = ? AND questionnaire_permissions.user_id = ? AND permissions.name = ?",
				questionnaireID, userID, requiredPermission).First(&permission).Error

		if err != nil || permission.ID == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not have the necessary permissions"})
		}

		// Proceed to the next handler if the user has permission
		return c.Next()
	}
}

func MockAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mock user ID (normally this comes from token or session)
		c.Locals("user_id", uint(2)) // Assuming user ID 1
		return c.Next()
	}
}
