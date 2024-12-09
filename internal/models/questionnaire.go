package models

import "time"

type Questionnaire struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Title               string         `gorm:"not null" json:"title"`
	CreationTime        time.Time      `gorm:"autoCreateTime" json:"creation_time"`
	StartTime           time.Time      `json:"start_time"`
	EndTime             time.Time      `json:"end_time"`
	OrderType           OrderType      `gorm:"not null;default:'Sequential'" json:"order_type"`
	AllowReturn         bool           `gorm:"default:true" json:"allow_return"`
	MaxParticipation    int            `gorm:"default:0" json:"max_participation"`
	ResponseTime        int            `gorm:"default:0" json:"response_time"`
	AnonymityLevel      AnonymityLevel `gorm:"not null;default:'Invisible'"`
	AgeRestriction      *int           `json:"age_restriction"`
	LocationRestriction *string        `json:"location_restriction"`
	GenderRestriction   *string        `json:"gender_restriction"`
	WithdrawalDeadline  time.Time      `json:"withdrawal_deadline"`

	OwnerID uint `json:"owner_id"`
	Owner   User `gorm:"foreignKey:OwnerID" json:"owner"`

	// Relationships
	Permissions []Permission `gorm:"many2many:questionnaire_permissions;" json:"permissions"`
	Questions   []*Question  `gorm:"foreignKey:QuestionnaireID"`
	Responses   []*Response  `gorm:"foreignKey:QuestionnaireID" json:"responses"`
}

type OrderType string

const (
	Sequential OrderType = "Sequential"
	Random     OrderType = "Random"
)

type AnonymityLevel string

const (
	VisibleToAll   AnonymityLevel = "All"
	CreatorOrAdmin AnonymityLevel = "CreatorOrAdmin"
	Invisible      AnonymityLevel = "Invisible"
)
