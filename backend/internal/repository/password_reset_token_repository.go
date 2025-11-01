package repository

import (
	"context"
	"time"
	"livekit-consulting/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetByToken(ctx context.Context, token string) (*model.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}

type passwordResetTokenRepository struct {
	db *sqlx.DB
}

func NewPasswordResetTokenRepository(db *sqlx.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *passwordResetTokenRepository) GetByToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	query := `
		SELECT id, user_id, token, expires_at, used
		FROM password_reset_tokens
		WHERE token = $1
	`
	err := r.db.GetContext(ctx, &resetToken, query, token)
	return &resetToken, err
}

func (r *passwordResetTokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE password_reset_tokens
		SET used = true
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
