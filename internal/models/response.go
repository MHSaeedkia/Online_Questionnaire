package models

import "time"

type Response struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	QuestionnaireID uint      `json:"questionnaire_id"`
	QuestionID      uint      `json:"question_id"`
	UserID          uint      `json:"user_id"`
	Answer          string    `gorm:"not null" json:"answer"` // User's answer (encrypted)
	Secret          string    `json:"secret"`                 // Unique secret for each vote
	SubmittedAt     time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	IsWithdrawn     bool      `gorm:"default:false" json:"is_withdrawn"`
}
