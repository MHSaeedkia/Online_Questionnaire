package models

import "time"

type QuestionnairePermission struct {
	ID              uint          `gorm:"primaryKey"`
	QuestionnaireID uint          `json:"questionnaire_id"`
	Questionnaire   Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	UserID          uint          `json:"user_id"`
	User            User          `gorm:"foreignKey:UserID"`
	PermissionID    uint          `json:"permission_id"`
	Permission      Permission    `gorm:"foreignKey:PermissionID"`
	ExpiresAt       *time.Time    `json:"valid_until"`
	Status          string        `gorm:"default:'Pending'" json:"status"`
}
