package models

type Option struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	QuestionID uint   `json:"question_id"`
	Text       string `gorm:"not null" json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}
