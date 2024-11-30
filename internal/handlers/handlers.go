package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

//	@title			Questionnaire API
//	@version		1.0
//	@description	This is a sample API for the questionnaire service.
//	@BasePath		/api

// @Summary		Get version of the API
// @Description	Returns the version of the API
// @Tags			version
// @Accept			json
// @Produce		json
// @Success		200	{string}	string	"ok"
// @Router			/version [get]
func GetVersion(c *fiber.Ctx) error {
	fmt.Println("Getting version!.....")
	version := os.Getenv("PROJECT_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	return c.JSON(fiber.Map{
		"version": version,
	})
}
