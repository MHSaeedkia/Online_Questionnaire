package routers

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// get config .
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

	questionnaireRepo := repositories.NewQuestionnaireRepository(DB)
	questionRepo := repositories.NewQuestionRepository(DB)
	optionRepo := repositories.NewOptionRepository(DB)
	conditionalLogicRepo := repositories.NewConditionalLogicRepository(DB)
	permissionRepo := repositories.NewPermissionRepository(DB)
	responseRepo := repositories.NewResponseRepository(DB)
	userRepository := repositories.NewUserRepository(DB)

	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireRepo)
	questionHandler := handlers.NewQuestionHandler(questionRepo)
	optionHandler := handlers.NewOptionHandler(optionRepo, questionRepo)
	conditionalLogicHandler := handlers.NewConditionalLogicHandler(conditionalLogicRepo, questionRepo, optionRepo)
	permissionHandler := handlers.NewPermissionHandler(questionnaireRepo, permissionRepo)
	responseHandler := handlers.NewResponseHandler(responseRepo)
	userHandler := handlers.NewUserHandler(userRepository)

	api := app.Group("/api/v1")

	questionnaireRoutes := api.Group("/questionnaires")
	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire)

	questionRoutes := api.Group("/questions")
	questionRoutes.Post("/", questionHandler.CreateQuestion)
	questionRoutes.Get("/", questionHandler.GetQuestion)
	questionRoutes.Put("/", questionHandler.UpdateQuestion)
	questionRoutes.Delete("/", questionHandler.DeleteQuestion)

	optionRoutes := api.Group("/options")
	optionRoutes.Post("/question", optionHandler.CreateOptions)

	conditionalLogicRoutes := api.Group("/conditional-logic")
	conditionalLogicRoutes.Post("/question", conditionalLogicHandler.CreateConditionalLogic)

	permitionRoutes := api.Group("/permition")
	permitionRoutes.Post("/", permissionHandler.RequestPermission)
	permitionRoutes.Put("/status", permissionHandler.ApproveOrDenyPermissionRequest)

	answerRoutes := api.Group("/response")
	answerRoutes.Post("/fill", middleware.CheckPermission(DB, models.CanViewVote), responseHandler.FillQuestionnaire)
	answerRoutes.Put("/edit", middleware.CheckPermission(DB, models.CanEdit), responseHandler.EditResponse)
	answerRoutes.Post("/", responseHandler.CreateResponse)
	answerRoutes.Get("/", responseHandler.GetResponse)
	answerRoutes.Put("/", responseHandler.UpdateResponse)
	answerRoutes.Delete("/", responseHandler.DeleteResponse)

	userRouter := api.Group("/user")
	userRouter.Get("/questionnaires", userHandler.Quesionnare)
	userRouter.Put("/questionnaires/edit", userHandler.EditQuestionnare)
	userRouter.Delete("/questionnaires/cancle", userHandler.CancleQuestionnarec)

}
