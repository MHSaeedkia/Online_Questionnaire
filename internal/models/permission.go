package models

type Permission struct {
	ID   uint `gorm:"primaryKey" json:"id"`
	Name Name `gorm:"unique;not null" json:"name"`

	// Relationships
	Users          []User          `gorm:"many2many:user_permissions;" json:"users"`
	Questionnaires []Questionnaire `gorm:"many2many:questionnaire_permissions;" json:"questionnaires"`
}

type Name string

const (
	CanView        Name = "CanView"
	CanViewVote    Name = "CanViewVote"
	CanVote        Name = "CanVote"
	CanEdit        Name = "CanEdit"
	CanChangeRole  Name = "CanChangeRole"
	CanViewReports Name = "CanViewReports"
)
