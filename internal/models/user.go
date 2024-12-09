package models

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NationalID    string    `gorm:"unique;not null" json:"national_id"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Gender        string    `json:"gender"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	City          string    `json:"city"`
	WalletBalance float64   `gorm:"default:0.0" json:"wallet_balance"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Role          Role      `gorm:"default:User" json:"role"`

	// Relationships
	Questionnaires []*Questionnaire `gorm:"foreignKey:OwnerID" json:"questionnaires"`
	Permissions    []Permission     `gorm:"many2many:user_permissions;" json:"permissions"`
	Notifications  []Notification   `gorm:"foreignKey:UserID" json:"notifications"`
	Responses      []*Response      `json:"responses"`
}

type Role string

const (
	SuperAdmin Role = "SuperAdmin"
	Guest      Role = "User"
)
