package services

import (
	"net/smtp"

	"github.com/Cinnamoon-dev/blue-gopher/pkg/config"
)

type MailService struct{}

func NewMailService() MailService {
	return MailService{}
}

func (s *MailService) SendEmail(to string, subject string, body string) error {
	env := config.NewEnv()
	from := env.MailUsername
	pass := env.MailPassword

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
