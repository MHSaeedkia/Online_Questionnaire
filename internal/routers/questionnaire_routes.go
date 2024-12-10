package routers

import (
	"fmt"
	"log"
	handlers2 "online-questionnaire/api/handlers/permission_handler"
	"online-questionnaire/api/handlers/questionnaire_handlers"
	"online-questionnaire/api/handlers/response_handler"
	"online-questionnaire/api/handlers/user_handler"
	"online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/middlewares"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories/permission_repo"
	"online-questionnaire/internal/repositories/questionnaire_repo"
	"online-questionnaire/internal/repositories/response_repo"
	"online-questionnaire/internal/repositories/user_repo"
	"online-questionnaire/internal/services"
	"online-questionnaire/internal/utils"
	"online-questionnaire/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(cfg configs.Config, app *fiber.App) {
	// Connect to the database
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	// Connect to Redis
	redisClient := utils.NewRedisClient()

	// Configure email sender
	emailSender := &utils.SendEmail{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUsername: "golang.project2@gmail.com",
		SMTPPassword: "xdgrmaefztthqowj",
	}

	// API route group
	api := app.Group("/api")
	api.Use(middlewares.RateLimiter(redisClient.Client, 1)) // Limit to 1 request per second
	userRoutes := api.Group("/users")
	questionnaireRoutes := api.Group("/questionnaires")

	// Setup JWT middleware
	fmt.Println("cfg.JWT.Secret=", cfg.JWT.Secret)
	jwtMiddleware := middleware.JWTMiddleware(cfg.JWT.Secret)

	// Initialize repositories
	userRepo := user_repo.NewUserRepository(DB)
	questionnaireRepo := questionnaire_repo.NewQuestionnaireRepository(DB)
	questionRepo := questionnaire_repo.NewQuestionRepository(DB)
	optionRepo := questionnaire_repo.NewOptionRepository(DB)
	conditionalLogicRepo := questionnaire_repo.NewConditionalLogicRepository(DB)
	permissionRepo := permission_repo.NewPermissionRepository(DB)
	responseRepo := response_repo.NewResponseRepository(DB)

	// Initialize services
	userService := services.NewUserService(userRepo, cfg)
	oauthService := services.NewOAuthService(cfg, userRepo)
	verificationService := services.NewVerificationService(redisClient, emailSender)

	// Initialize handlers
	userHandler := user_handler.NewUserHandler(userService)
	oauthHandler := user_handler.NewOAuthHandler(oauthService)
	verificationHandler := user_handler.NewVerificationHandler(verificationService)

	questionnaireHandler := questionnaire_handlers.NewQuestionnaireHandler(questionnaireRepo)
	questionHandler := questionnaire_handlers.NewQuestionHandler(questionRepo)
	optionHandler := questionnaire_handlers.NewOptionHandler(optionRepo, questionRepo)
	conditionalLogicHandler := questionnaire_handlers.NewConditionalLogicHandler(conditionalLogicRepo, questionRepo, optionRepo)
	permissionHandler := handlers2.NewPermissionHandler(questionnaireRepo, permissionRepo)
	responseHandler := response_handler.NewResponseHandler(responseRepo)
	voteHandler := questionnaire_handlers.NewVoteHandler(questionnaireRepo, questionRepo, responseRepo)

	// Route for Google OAuth login
	api.Post("/user/oauth", oauthHandler.GoogleLogin)

	// Verification routes for email verification
	api.Post("/verification/send", verificationHandler.SendVerificationCode)
	api.Post("/verification/validate", verificationHandler.ValidateVerificationCode)

	// User-related routes
	userRoutes.Post("/signup", middlewares.FixDateOfBirth, userHandler.Signup) // User signup
	userRoutes.Post("/login", userHandler.Login) // User login

	// Questionnaire CRUD routes
	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire) // Create a new questionnaire
	questionnaireRoutes.Post("/:questionnaire_id/questions", questionHandler.CreateQuestion) // Add a question to a questionnaire
	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/options", optionHandler.CreateOptions) // Add options to a question
	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/conditional-logic", conditionalLogicHandler.CreateConditionalLogic) // Add conditional logic to a question

	// Permission-related routes
	questionnaireRoutes.Post("/:questionnaireID/permissions/request", jwtMiddleware, permissionHandler.RequestPermission) // Request permissions for a questionnaire
	questionnaireRoutes.Put("/permissions/:requestID", permissionHandler.ApproveOrDenyPermissionRequest) // Approve or deny permission requests

	// Voting route for questionnaires
	questionnaireRoutes.Post("/:id/vote", voteHandler.VoteOnQuestionnaire) // Submit a vote for a questionnaire

	// Response-related routes
	questionnaireRoutes.Post("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanViewVote), responseHandler.FillQuestionnaire) // Submit a response to a questionnaire
	questionnaireRoutes.Put("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanEdit), responseHandler.EditResponse) // Edit an existing response
}

