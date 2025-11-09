package wildlifenl

import (
	"bytes"
	"log"
	"text/template"
	"time"

	_ "embed"

	"github.com/go-mail/mail"
)

const (
	emailSubject = "Aanmelden bij WildlifeNL"
)

var (
	//go:embed templates/email.go.tmpl
	emailTemplateFS string
	emailTemplate   *template.Template
)

func init() {
	emailTemplate = template.Must(template.New("email").Parse(emailTemplateFS))
}

type mailData struct {
	AppName     string
	DisplayName string
	Code        string
	Year        int
}

type Mailer struct {
	config *Configuration
}

func newMailer(config *Configuration) *Mailer {
	return &Mailer{config: config}
}

func (e *Mailer) Ping() error {
	if e.config.EmailHost == "no-email" {
		return nil
	}
	s, err := e.dailer().Dial()
	if err != nil {
		return err
	}
	return s.Close()
}

func (e *Mailer) SendCode(appName, displayName, email, code string) error {
	if e.config.EmailHost == "no-email" {
		log.Println("Code for", email, "is:", code)
		return nil
	}
	var bodyBuffer bytes.Buffer
	err := emailTemplate.Execute(&bodyBuffer, mailData{
		AppName:     appName,
		DisplayName: displayName,
		Code:        code,
		Year:        time.Now().Year(),
	})
	if err != nil {
		return err
	}
	body := bodyBuffer.String()
	m := mail.NewMessage()
	m.SetHeader("From", e.config.EmailFrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", body)
	return e.dailer().DialAndSend(m)
}

func (e *Mailer) dailer() *mail.Dialer {
	return mail.NewDialer(e.config.EmailHost, 587, e.config.EmailUser, e.config.EmailPassword)
}
