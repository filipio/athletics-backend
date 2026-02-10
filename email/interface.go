package email

import (
	"context"

	"github.com/resend/resend-go/v3"
)

//go:generate mockgen -destination=../mocks/mock_email_sender.go -package=mocks github.com/filipio/athletics-backend/email EmailSender

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, params VerificationEmailParams) error
}

type ResendEmailSender struct {
	client *resend.Client
}

func NewResendEmailSender(client *resend.Client) *ResendEmailSender {
	return &ResendEmailSender{client: client}
}

func (s *ResendEmailSender) SendVerificationEmail(ctx context.Context, params VerificationEmailParams) error {
	emailParams := resend.SendEmailRequest{
		To: []string{params.To},
		Template: &resend.EmailTemplate{
			Id: EmailVerificationTemplateID,
			Variables: map[string]any{
				"token": params.VerificationToken,
			},
		},
	}
	_, err := s.client.Emails.SendWithContext(ctx, &emailParams)
	return err
}

func GetDefaultEmailSender() EmailSender {
	return NewResendEmailSender(GetClient())
}
