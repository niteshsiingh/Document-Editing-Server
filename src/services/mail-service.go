package services

import (
	"fmt"
	"os"

	"github.com/sendinblue/APIv3-go-library/v2/lib"
)

type MailOptions struct {
	From    string
	To      []string
	Subject string
	Body    string
}

type MailService struct {
	client *lib.APIClient
}

func NewMailService() (*MailService, error) {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BREVO_API_KEY environment variable is not set")
	}

	cfg := lib.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := lib.NewAPIClient(cfg)

	return &MailService{
		client: client,
	}, nil
}

func (ms *MailService) SendMail(mailOptions MailOptions) error {
	// Convert recipients to Brevo format
	var toList []lib.SendSmtpEmailTo
	for _, recipient := range mailOptions.To {
		toList = append(toList, lib.SendSmtpEmailTo{
			Email: recipient,
		})
	}

	// Create email data
	emailData := lib.SendSmtpEmail{
		Sender: &lib.SendSmtpEmailSender{
			Email: os.Getenv("BREVO_SENDER_EMAIL"),
		},
		To:          toList,
		Subject:     mailOptions.Subject,
		TextContent: mailOptions.Body,
	}

	// Send email
	_, _, err := ms.client.TransactionalEmailsApi.SendTransacEmail(nil, emailData)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (ms *MailService) SendHTMLMail(mailOptions MailOptions, htmlContent string) error {
	var toList []lib.SendSmtpEmailTo
	for _, recipient := range mailOptions.To {
		toList = append(toList, lib.SendSmtpEmailTo{
			Email: recipient,
		})
	}

	emailData := lib.SendSmtpEmail{
		Sender: &lib.SendSmtpEmailSender{
			Email: mailOptions.From,
		},
		To:          toList,
		Subject:     mailOptions.Subject,
		TextContent: mailOptions.Body,    // Fallback plain text
		HtmlContent: htmlContent,         // HTML version
	}

	_, _, err := ms.client.TransactionalEmailsApi.SendTransacEmail(nil, emailData)
	if err != nil {
		return fmt.Errorf("failed to send HTML email: %v", err)
	}

	return nil
}