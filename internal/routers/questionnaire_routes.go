package routers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	handlers2 "online-questionnaire/api/handlers/permission_handler"
	"online-questionnaire/api/handlers/questionnaire_handlers"
	"online-questionnaire/api/handlers/response_handler"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories/permission_repo"
	"online-questionnaire/internal/repositories/questionnaire_repo"
	"online-questionnaire/internal/repositories/response_repo"
	"online-questionnaire/pkg/middleware"
)

func SetupRoutes(app *fiber.App) {
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatal(err)
	}

	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	api := app.Group("/api")
	questionnaireRoutes := api.Group("/questionnaires")

	questionnaireRepo := questionnaire_repo.NewQuestionnaireRepository(DB)
	questionRepo := questionnaire_repo.NewQuestionRepository(DB)
	optionRepo := questionnaire_repo.NewOptionRepository(DB)
	conditionalLogicRepo := questionnaire_repo.NewConditionalLogicRepository(DB)
	permissionRepo := permission_repo.NewPermissionRepository(DB)
	responseRepo := response_repo.NewResponseRepository(DB)

	questionnaireHandler := questionnaire_handlers.NewQuestionnaireHandler(questionnaireRepo)
	questionHandler := questionnaire_handlers.NewQuestionHandler(questionRepo)
	optionHandler := questionnaire_handlers.NewOptionHandler(optionRepo, questionRepo)
	conditionalLogicHandler := questionnaire_handlers.NewConditionalLogicHandler(conditionalLogicRepo, questionRepo, optionRepo)
	permissionHandler := handlers2.NewPermissionHandler(questionnaireRepo, permissionRepo)
	responseHandler := response_handler.NewResponseHandler(responseRepo)

	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire)

	questionnaireRoutes.Post("/:questionnaire_id/questions", questionHandler.CreateQuestion)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/options", optionHandler.CreateOptions)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/conditional-logic", conditionalLogicHandler.CreateConditionalLogic)

	//questionnaireRoutes.Post("/:questionnaireID/permissions", permissionHandler.GrantPermissionToUser)

	questionnaireRoutes.Post("/:questionnaireID/permissions/request", permissionHandler.RequestPermission)

	questionnaireRoutes.Put("/permissions/:requestID", permissionHandler.ApproveOrDenyPermissionRequest)

	questionnaireRoutes.Post("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanViewVote), responseHandler.FillQuestionnaire)

	questionnaireRoutes.Put("/:questionnaire_id/responses", middleware.CheckPermission(DB, models.CanEdit), responseHandler.EditResponse)
}
