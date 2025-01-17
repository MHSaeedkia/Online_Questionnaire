package db

import (
	"fmt"
	"online-questionnaire/configs"
	"online-questionnaire/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *configs.DatabaseConfig) (*gorm.DB, error) {
	//dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
	//	cfg.Host,
	//	cfg.User,
	//	cfg.Password,
	//	cfg.DBName,
	//	cfg.Port,
	//	cfg.SSLMode,
	//)
	dsn := "postgres://root:root@localhost:5432/OnlineQuestionnaire"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Database connection failed with DSN: %s\n", dsn)
		return nil, err
	}

	if err := migrate(db); err != nil {
		fmt.Printf("migrations failed: %v\n", err.Error())
		return nil, err
	}
	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Questionnaire{},
		&models.Permission{},
		&models.QuestionnairePermission{},
		&models.Question{},
		&models.Option{},
		&models.Response{},
		&models.ConditionalLogic{},
		&models.Notification{},
	)
}
