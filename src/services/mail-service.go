package services

import (
	"crypto/tls"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type MailOptions struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type MailService struct {
	Transporter *gomail.Dialer
}

func NewMailService() (*MailService, error) {
	// Replace "smtp.example.com" with the SMTP server address you want to use for local testing
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	transporter := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	transporter.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &MailService{
		Transporter: transporter,
	}, nil
}

func (ms *MailService) SendMail(mailOptions MailOptions) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", mailOptions.From)
	msg.SetHeader("Subject", mailOptions.Subject)
	msg.SetBody("text/plain", mailOptions.Body)

	for _, to := range mailOptions.To {
		msg.SetHeader("To", to)
		err := ms.Transporter.DialAndSend(msg)
		if err != nil {
			return err
		}

	}
	return nil
}
