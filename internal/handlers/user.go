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
// @Summary User Login
// @Description Authenticate a user by providing email or national ID and password to receive a JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login Request"
// @Success 200 {object} models.User "Successfully logged in"
// @Failure 400 {object} models.ErrorResponse "Invalid request payload"
// @Failure 401 {object} models.ErrorResponse "Incorrect credentials"
// @Router /api/user/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest map[string]string
	if err := c.BodyParser(&loginRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	emailOrNationalID, password := loginRequest["email_or_national_id"], loginRequest["password"]
	if emailOrNationalID == "" || password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or National ID and password are required")
	}

	// Call the service layer to perform login
	token, loginErr := h.UserService.Login(emailOrNationalID, password)
	if loginErr != nil {
		return fiber.NewError(fiber.StatusUnauthorized, loginErr.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User logged in successfully",
		"token":   token,
	})
}
