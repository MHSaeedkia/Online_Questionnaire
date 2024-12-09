package utils

import (
	"fmt"
	"net/smtp"
)

// SendEmail struct to manage email sending configuration
type SendEmail struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// SendEmail sends an email via SMTP
func (s *SendEmail) SendEmail(to, subject, body string) error {
	from := s.SMTPUsername
	auth := smtp.PlainAuth("", from, s.SMTPPassword, s.SMTPHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort), auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
