package mailer

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer    *gomail.Dialer
	fromEmail string
	fromName  string
}

func NewMailer() *Mailer {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587
	}
	return &Mailer{
		dialer: gomail.NewDialer(
			os.Getenv("SMTP_HOST"),
			port,
			os.Getenv("SMTP_USERNAME"),
			os.Getenv("SMTP_PASSWORD"),
		),
		fromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		fromName:  os.Getenv("SMTP_FROM_NAME"),
	}
}

func (m *Mailer) SendEmail(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.fromName+"<"+m.fromEmail+">")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}
