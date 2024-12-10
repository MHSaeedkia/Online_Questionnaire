package models

type ConditionalLogic struct {
	ID               uint `gorm:"primaryKey"`
	QuestionID       uint `json:"question_id"`        // The question this condition belongs to
	OptionID         uint `json:"option_id"`          // Option that triggers this condition
	TargetQuestionID uint `json:"target_question_id"` // Question to be shown if the condition is met
}
