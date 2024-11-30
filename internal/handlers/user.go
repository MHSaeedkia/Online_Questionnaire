package handlers

import (
	"github.com/gofiber/fiber/v2"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req authRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	user := models.User{
		Email:    req.Email,
		Password: utils.GeneratePassword(req.Password),
	}
	res := repositories.UserRepository{}.SaveUser(*user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "user created",
	})
}
