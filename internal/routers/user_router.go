package routers

import (
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/repositories"

	"gorm.io/gorm"
)

func userRouter(db *gorm.DB) *handlers.UserHandler {
	userRepo := repositories.NewUserRepository(db)
	return handlers.NewUserHandler(userRepo)
}
