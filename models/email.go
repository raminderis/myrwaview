package models

import (
	"fmt"
	"os"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@live.com"
)

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es
}

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

func (service *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case service.DefaultSender != "":
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func (service *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	service.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
	msg.WriteTo(os.Stdout)
	//err := service.dialer.DialAndSend(msg)
	// if err != nil {
	// 	return fmt.Errorf("send: %w", err)
	// }
	return nil
}

func (service *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password link is: " + resetURL,
		HTML:      `<p>To reset your password link is: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}
	err := service.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password: %w", err)
	}
	return nil
}
