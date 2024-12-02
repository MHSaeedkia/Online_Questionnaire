package handlers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/services"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

// Signup godoc
//
//	@Summary		User Signup
//	@Description	Register a new user by providing email, national ID, password, and other optional details.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User				true	"User Signup Request"
//	@Success		201		{object}	map[string]interface{}	"User created successfully with a JWT token"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		409		{object}	map[string]interface{}	"User already exists with this email or national ID"
//	@Router			/api/user/signup [post]
func (h *UserHandler) Signup(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, registerErr := h.UserService.SignUp(&user)
	if registerErr != nil {
		return fiber.NewError(fiber.StatusConflict, registerErr.Error())
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"token":   token,
	})
}
