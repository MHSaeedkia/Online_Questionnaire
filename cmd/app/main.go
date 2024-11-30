package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/routers"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

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

	routers.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))

}
