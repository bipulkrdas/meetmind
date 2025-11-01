package email

import (
    "context"
    mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type MailjetProvider struct {
    client    *mailjet.Client
    fromEmail string
    fromName  string
}

func NewMailjetProvider(apiKey, secretKey, fromEmail, fromName string) *MailjetProvider {
    client := mailjet.NewMailjetClient(apiKey, secretKey)
    return &MailjetProvider{
        client:    client,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (p *MailjetProvider) SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error {
    messagesInfo := []mailjet.InfoMessagesV31{
        {
            From: &mailjet.RecipientV31{
                Email: p.fromEmail,
                Name:  p.fromName,
            },
            To: &mailjet.RecipientsV31{
                mailjet.RecipientV31{
                    Email: to,
                },
            },
            Subject:  subject,
            TextPart: textContent,
            HTMLPart: htmlContent,
        },
    }
    
    messages := mailjet.MessagesV31{Info: messagesInfo}
    _, err := p.client.SendMailV31(&messages)
    return err
}

func (p *MailjetProvider) SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error {
    // Implement template-based sending if needed
    return nil
}
