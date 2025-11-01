package email

import (
    "context"
)

// EmailProvider is the interface that all email providers must implement
type EmailProvider interface {
    SendEmail(ctx context.Context, to, subject, htmlContent, textContent string) error
    SendTemplateEmail(ctx context.Context, to string, templateID string, data map[string]interface{}) error
}

// EmailService handles email operations
type EmailService struct {
    provider EmailProvider
    fromEmail string
    fromName  string
}

func NewEmailService(provider EmailProvider, fromEmail, fromName string) *EmailService {
    return &EmailService{
        provider:  provider,
        fromEmail: fromEmail,
        fromName:  fromName,
    }
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
    subject := "Password Reset Request"
    htmlContent := `
        <html>
        <body>
            <h2>Password Reset Request</h2>
            <p>You have requested to reset your password. Click the link below to reset:</p>
            <a href="` + resetURL + `">Reset Password</a>
            <p>This link will expire in 1 hour.</p>
            <p>If you did not request this, please ignore this email.</p>
        </body>
        </html>
    `
    textContent := "Password reset link: " + resetURL
    
    return s.provider.SendEmail(ctx, to, subject, htmlContent, textContent)
}

func (s *EmailService) SendRoomInviteEmail(ctx context.Context, to, roomName, inviteURL string) error {
    subject := "You've been invited to join a meeting room"
    htmlContent := `
        <html>
        <body>
            <h2>Meeting Room Invitation</h2>
            <p>You have been invited to join the room: <strong>` + roomName + `</strong></p>
            <p>Click the link below to join:</p>
            <a href="` + inviteURL + `">Join Room</a>
        </body>
        </html>
    `
    textContent := "Join room " + roomName + ": " + inviteURL
    
    return s.provider.SendEmail(ctx, to, subject, htmlContent, textContent)
}
