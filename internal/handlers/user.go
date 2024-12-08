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

// Login godoc
//
//	@Summary		User Login
//	@Description	User Login using national ID and password to obtain a JWT token.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		map[string]string	true	"Login credentials (national ID and password)"
//	@Success		200			{object}	map[string]interface{}	"User logged in successfully with a JWT token"
//	@Failure		400			{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		401			{object}	map[string]interface{}	"Invalid credentials"
//	@Router			/api/user/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest map[string]string
	if err := c.BodyParser(&loginRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	nationalID, password := loginRequest["national_id"], loginRequest["password"]
	if nationalID == "" || password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "National ID and password are required")
	}

	// Call the service layer to perform login
	token, loginErr := h.UserService.Login(nationalID, password)
	if loginErr != nil {
		return fiber.NewError(fiber.StatusUnauthorized, loginErr.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User logged in successfully",
		"token":   token,
	})
}
