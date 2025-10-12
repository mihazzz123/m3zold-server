package email

import "context"

// EmailSender интерфейс для отправки email
type EmailSender interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
	SendWelcomeEmail(ctx context.Context, email, name string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
}

// EmailTemplate тип для шаблонов email
type EmailTemplate string

const (
	TemplateVerification  EmailTemplate = "verification"
	TemplateWelcome       EmailTemplate = "welcome"
	TemplatePasswordReset EmailTemplate = "password_reset"
)

// EmailData данные для отправки email
type EmailData struct {
	To          string
	Subject     string
	Template    EmailTemplate
	Data        map[string]interface{}
	Attachments []Attachment
}

// Attachment вложение email
type Attachment struct {
	Name        string
	Content     []byte
	ContentType string
}
