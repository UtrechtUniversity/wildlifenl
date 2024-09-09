package wildlifenl

import (
	"strings"

	"github.com/go-mail/mail"
)

const emailSubject = "Aanmelden bij WildlifeNL"
const emailBody = "Beste {displayName}<br/>De applicatie {appName} wil graag aanmelden bij WildlifeNL met jouw emailadres. Om dit toe te staan, voer onderstaande code in bij deze applicatie.<br/>Code: {code}<br/><br/>Met vriendelijke groet<br/>WildlifeNL<br/><br/>"

func sendCodeByEmail(appName, displayName, email, code string) error {
	if configuration.EmailHost == "no-email" {
		return nil
	}
	body := emailBody
	body = strings.ReplaceAll(body, "{appName}", appName)
	body = strings.ReplaceAll(body, "{displayName}", displayName)
	body = strings.ReplaceAll(body, "{code}", code)
	m := mail.NewMessage()
	m.SetHeader("From", configuration.EmailFrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", body)
	d := mail.NewDialer(configuration.EmailHost, 587, configuration.EmailUser, configuration.EmailPassword)
	return d.DialAndSend(m)
}
