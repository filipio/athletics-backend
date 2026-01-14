package email

import (
	"context"

	"github.com/resend/resend-go/v3"
)

type VerificationEmailParams struct {
	To                string
	VerificationToken string
}

func SendVerificationEmail(ctx context.Context, params VerificationEmailParams) error {
	client := GetClient()

	emailParams := resend.SendEmailRequest{
		To: []string{params.To},
		Template: &resend.EmailTemplate{
			Id: EmailVerificationTemplateID,
			Variables: map[string]any{
				"token": params.VerificationToken,
			},
		},
	}

	_, err := client.Emails.SendWithContext(ctx, &emailParams)
	return err
}
