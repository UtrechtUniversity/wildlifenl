package wildlifenl

import (
	"log"
	"strings"

	"github.com/go-mail/mail"
)

type Mailer struct {
	config *Configuration
}

func newMailer(config *Configuration) *Mailer {
	return &Mailer{config: config}
}

func (e *Mailer) Ping() error {
	// TODO #36: test connection to email provider.
	return nil
}

func (e *Mailer) SendCode(appName, displayName, email, code string) error {
	if e.config.EmailHost == "no-email" {
		log.Println("Code for", email, "is:", code)
		return nil
	}
	body := emailBody
	body = strings.ReplaceAll(body, "{appName}", appName)
	body = strings.ReplaceAll(body, "{displayName}", displayName)
	body = strings.ReplaceAll(body, "{code}", code)
	m := mail.NewMessage()
	m.SetHeader("From", e.config.EmailFrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", body)
	d := mail.NewDialer(e.config.EmailHost, 587, e.config.EmailUser, e.config.EmailPassword)
	return d.DialAndSend(m)
}
