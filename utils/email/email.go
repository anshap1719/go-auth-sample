package email

import (
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(subject, toName, toMail, textContent, htmlContent string) error {
	from := mail.NewEmail("Giggle", "anshap1719@gmail.com")
	to := mail.NewEmail(toName, toMail)
	plainTextContent := textContent
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(secrets.SendgridAPIKey)
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func SendMailMulti(subject, name, textContent, htmlContent string, receipients []map[string]string) error {
	m := mail.NewV3Mail()

	from := mail.NewEmail("Giggle <"+name+">", "anshap1719@gmail.com")
	html := mail.NewContent("text/html", htmlContent)
	text := mail.NewContent("text/plain", textContent)

	m.SetFrom(from)
	m.AddContent(html)
	m.AddContent(text)

	p := mail.NewPersonalization()

	var mails []*mail.Email

	for _, r := range receipients {
		mails = append(mails, mail.NewEmail(r["name"], r["email"]))
	}

	p.AddTos(mails...)

	p.Subject = subject

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(secrets.SendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)

	if err != nil {
		return err
	}

	return nil
}
