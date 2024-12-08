package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func FixDateOfBirth(c *fiber.Ctx) error {
	// Parse the request body to get the raw JSON
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	// Check if the date_of_birth key exists
	if dateOfBirth, exists := requestBody["date_of_birth"].(string); exists {
		// Check if the date format is YYYY-MM-DD
		if len(dateOfBirth) == 10 && strings.Contains(dateOfBirth, "-") {
			// Fix the date format by appending "T00:00:00Z"
			updatedDate := dateOfBirth + "T00:00:00Z"

			// Update the date_of_birth field in the request body
			requestBody["date_of_birth"] = updatedDate

			// Log the updated date for debugging
			log.Println("Updated date_of_birth:", updatedDate)

			// Set the updated request body back into the context
			c.Locals("updatedBody", requestBody)
		}
	}

	// Call the next handler in the chain
	return c.Next()
}
