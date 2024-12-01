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
	Responses      []*Response      `gorm:"foreignKey:QuestionnaireID"`
}

type Role string

const (
	SuperAdmin Role = "SuperAdmin"
	Guest      Role = "User"
)

type Questionnaire struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Title            string         `gorm:"not null" json:"title"`
	CreationTime     time.Time      `gorm:"autoCreateTime" json:"creation_time"`
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	OrderType        OrderType      `gorm:"not null;default:'Sequential'" json:"order_type"`
	AllowReturn      bool           `gorm:"default:true" json:"allow_return"`
	MaxParticipation int            `gorm:"default:0" json:"max_participation"`
	ResponseTime     int            `gorm:"default:0" json:"response_time"`
	AnonymityLevel   AnonymityLevel `gorm:"not null;default:'Invisible'"`

	OwnerID uint `json:"owner_id"`
	Owner   User `gorm:"foreignKey:OwnerID" json:"owner"`

	// Relationships
	Permissions []Permission `gorm:"many2many:questionnaire_permissions;" json:"permissions"`
	Questions   []*Question  `gorm:"foreignKey:QuestionnaireID"`

	AgeRestriction      *int    `json:"age_restriction"`
	LocationRestriction *string `json:"location_restriction"`
	GenderRestriction   *string `json:"gender_restriction"`

	WithdrawalDeadline time.Time `json:"withdrawal_deadline"`
}

type Response struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	QuestionnaireID uint      `json:"questionnaire_id"`
	QuestionID      uint      `json:"question_id"`
	UserID          uint      `json:"user_id"`                       // Who submitted the response
	Answer          string    `gorm:"not null" json:"answer"`        // User's answer (encrypted)
	Secret          string    `gorm:"not null;unique" json:"secret"` // Unique secret for each vote
	SubmittedAt     time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	IsWithdrawn     bool      `gorm:"default:false" json:"is_withdrawn"`
}

type OrderType string
type AnonymityLevel string

const (
	Sequential OrderType = "Sequential"
	Random     OrderType = "Random"
)
const (
	VisibleToAll   AnonymityLevel = "All"
	CreatorOrAdmin AnonymityLevel = "CreatorOrAdmin"
	Invisible      AnonymityLevel = "Invisible"
)

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

type QuestionnairePermission struct {
	ID              uint          `gorm:"primaryKey"`
	QuestionnaireID uint          `json:"questionnaire_id"`
	Questionnaire   Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	UserID          uint          `json:"user_id"`
	User            User          `gorm:"foreignKey:UserID"`
	PermissionID    uint          `json:"permission_id"`
	Permission      Permission    `gorm:"foreignKey:PermissionID"`
}

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

type Option struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	QuestionID uint   `json:"question_id"`
	Text       string `gorm:"not null;unique" json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

type ConditionalLogic struct {
	ID               uint `gorm:"primaryKey"`
	QuestionID       uint `json:"question_id"`        // The question this condition belongs to
	OptionID         uint `json:"option_id"`          // Option that triggers this condition
	TargetQuestionID uint `json:"target_question_id"` // Question to be shown if the condition is met
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`   // Notification title
	Message   string    `gorm:"not null" json:"message"` // Detailed message
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
}
