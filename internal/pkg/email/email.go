package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type VerificationMailer interface {
	SendMail(to string, subject string, view string, props map[string]string) error
}

type Email struct {
	SenderName string
	Dialer     *gomail.Dialer
}

func NewMailer() VerificationMailer {
	mailerPort, _ := strconv.Atoi(os.Getenv("MAILER_PORT"))
	mailer := gomail.NewDialer(
		os.Getenv("MAILER_HOST"),
		mailerPort,
		os.Getenv("MAILER_EMAIL"),
		os.Getenv("MAILER_PASSWORD"),
	)
	sender := fmt.Sprintf("%s <%s>", os.Getenv("APP_NAME"), os.Getenv("MAILER_EMAIL"))
	return &Email{sender, mailer}
}

func (e *Email) SendMail(to string, subject string, view string, props map[string]string) error {
	tmpl, err := template.ParseFiles(view)
	if err != nil {
		return err
	}
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, props); err != nil {
		return err
	}
	email := gomail.NewMessage()

	email.SetHeader("Subject", subject)
	email.SetHeader("From", e.SenderName)
	email.SetHeader("To", to)
	email.SetBody("text/html", buff.String())

	if err := e.Dialer.DialAndSend(email); err != nil {
		return err
	}

	return nil
}
