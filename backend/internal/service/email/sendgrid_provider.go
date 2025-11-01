package email

import (
    "context"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct {
    apiKey    string
    fromEmail string
    fromName  string
}

func NewSendGridProvider(apiKey, fromEmail, fromName string) *SendGridProvider {
    return &SendGridProvider{
        apiKey:    apiKey,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (p *SendGridProvider) SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error {
    from := mail.NewEmail(p.fromName, p.fromEmail)
    toEmail := mail.NewEmail("", to)
    message := mail.NewSingleEmail(from, subject, toEmail, textContent, htmlContent)
    
    client := sendgrid.NewSendClient(p.apiKey)
    _, err := client.Send(message)
    return err
}

func (p *SendGridProvider) SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error {
    // Implement template-based sending if needed
    return nil
}
