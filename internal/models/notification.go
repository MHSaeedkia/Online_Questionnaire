package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`   // Notification title
	Message   string    `gorm:"not null" json:"message"` // Detailed message
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
}
