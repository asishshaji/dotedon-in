package email_service

import "context"

type IEmailService interface {
	SendEmail(ctx context.Context, toEmail, subject, content string) error
}
