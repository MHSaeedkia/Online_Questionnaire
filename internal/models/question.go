package models

type Question struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	QuestionnaireID  uint    `json:"questionnaire_id"`
	Text             string  `gorm:"not null" json:"text"`
	Type             Type    `json:"type"`
	Order            int     `gorm:"not null" json:"order"` // Position in sequential order
	HasCorrectOption bool    `gorm:"default:false" json:"has_correct_option"`
	CorrectOptionID  *uint   `json:"correct_option_id"` // If it has a correct option, store its ID
	Media            *string `json:"media"`             // URL or path of image, video, audio, or file (optional)
	// One-to-One or One-to-Many relationship with ConditionalLogic
	ConditionalLogic []ConditionalLogic `gorm:"foreignKey:QuestionID" json:"conditional_logic"`

	Options []Option `gorm:"foreignKey:QuestionID" json:"options"`
}

type Type string

const (
	MultipleChoice Type = "MultipleChoice"
	Descriptive    Type = "Descriptive"
)
