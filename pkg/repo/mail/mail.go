package mail

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"github.com/sendgrid/sendgrid-go"
	mail2 "github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

const (
	SENDGRIDKEY = "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"
	UserName    = "Utk"
	SenderMail  = "utkarsh.wizwack@gmail.com"
)

type mail struct {
	client   *sendgrid.Client
	fromMail *mail2.Email
}

func NewMail() svc.Mail {
	client := sendgrid.NewSendClient(SENDGRIDKEY)
	from := mail2.NewEmail(UserName, SenderMail)
	return &mail{client: client, fromMail: from}
}

func (m *mail) SendMail(user model.User, otp string) error {
	subject := "forget password otp"
	to := mail2.NewEmail(user.Username, user.Email)
	textContent := "Your otp to reset password is : " + otp
	htmlContent := "<p> You otp to reset password is : " + otp + "</p>"
	email := mail2.NewSingleEmail(m.fromMail, subject, to, textContent, htmlContent)

	res, err := m.client.Send(email)
	if err != nil {
		log.Println("error in sending mail")
		return err
	}
	log.Println("status : ", res.StatusCode)
	log.Println("body : ", res.Body)
	log.Println("headers : ", res.Headers)
	if res.StatusCode != 200 && res.StatusCode != 202 {
		return svc.ErrUnexpected
	}
	return nil
}
