package routers

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"

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

	api := app.Group("/api/v1")
	questionnaireHandler := questionnareRouter(DB)
	userRepository := userRouter(DB)

	questionnaireRoutes := api.Group("/questionnaires")
	questionnaireRoutes.Post("/questionnaire", questionnaireHandler.CreateQuestionnaire)

	questionnaireRoutes.Post("/question/create", questionnaireHandler.CreateQuestion)
	questionnaireRoutes.Get("/question/get", questionnaireHandler.GetQuestion)
	questionnaireRoutes.Put("/question/update", questionnaireHandler.UpdateQuestion)
	questionnaireRoutes.Delete("/question/delete", questionnaireHandler.DeleteQuestion)

	questionnaireRoutes.Post("/answer/create", questionnaireHandler.CreateAnswer)
	questionnaireRoutes.Get("/answer/get", questionnaireHandler.GetAnswer)
	questionnaireRoutes.Put("/answer/update", questionnaireHandler.UpdateAnswer)
	questionnaireRoutes.Delete("/answer/delete", questionnaireHandler.DeleteAnswer)

	userRouter := api.Group("/user")
	userRouter.Get("/questionnaires", userRepository.Quesionnare)
	userRouter.Put("/questionnaires/edit", userRepository.EditQuestionnare)
	userRouter.Delete("/questionnaires/cancle", userRepository.CancleQuestionnarec)

}
