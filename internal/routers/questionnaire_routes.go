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
	"online-questionnaire/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(cfg configs.Config, app *fiber.App) {
	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	api := app.Group("/api")
	userRoutes := api.Group("/users")
	questionnaireRoutes := api.Group("/questionnaires")

	fmt.Println("cfg.JWT.Secret=", cfg.JWT.Secret)
	jwtMiddleware := middleware.JWTMiddleware(cfg.JWT.Secret)

	userRepo := user_repo.NewUserRepository(DB)
	userService := services.NewUserService(userRepo, cfg)

	questionnaireRepo := questionnaire_repo.NewQuestionnaireRepository(DB)
	questionRepo := questionnaire_repo.NewQuestionRepository(DB)
	optionRepo := questionnaire_repo.NewOptionRepository(DB)
	conditionalLogicRepo := questionnaire_repo.NewConditionalLogicRepository(DB)
	permissionRepo := permission_repo.NewPermissionRepository(DB)
	responseRepo := response_repo.NewResponseRepository(DB)

	userHandler := user_handler.NewUserHandler(userService)

	questionnaireHandler := questionnaire_handlers.NewQuestionnaireHandler(questionnaireRepo)
	questionHandler := questionnaire_handlers.NewQuestionHandler(questionRepo)
	optionHandler := questionnaire_handlers.NewOptionHandler(optionRepo, questionRepo)
	conditionalLogicHandler := questionnaire_handlers.NewConditionalLogicHandler(conditionalLogicRepo, questionRepo, optionRepo)
	permissionHandler := handlers2.NewPermissionHandler(questionnaireRepo, permissionRepo)
	responseHandler := response_handler.NewResponseHandler(responseRepo)

	userRoutes.Post("/signup", middlewares.FixDateOfBirth, userHandler.Signup)
	userRoutes.Post("/login", userHandler.Login)

	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire)

	questionnaireRoutes.Post("/:questionnaire_id/questions", questionHandler.CreateQuestion)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/options", optionHandler.CreateOptions)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/conditional-logic", conditionalLogicHandler.CreateConditionalLogic)

	questionnaireRoutes.Post("/:questionnaireID/permissions/request", jwtMiddleware, permissionHandler.RequestPermission)

	questionnaireRoutes.Put("/permissions/:requestID", permissionHandler.ApproveOrDenyPermissionRequest)

	questionnaireRoutes.Post("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanViewVote), responseHandler.FillQuestionnaire)

	questionnaireRoutes.Put("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanEdit), responseHandler.EditResponse)
}
