package repository

import (
    "context"
    "time"

    "livekit-consulting/backend/internal/model"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type InviteRepository interface {
    Create(ctx context.Context, roomID, inviterID uuid.UUID, inviteeEmail, inviteeName, token string, expiresAt time.Time) error
    GetByToken(ctx context.Context, token string) (*model.Invite, error)
    MarkAsAccepted(ctx context.Context, id uuid.UUID) error
}

type inviteRepository struct {
    db *sqlx.DB
}

func NewInviteRepository(db *sqlx.DB) InviteRepository {
    return &inviteRepository{db: db}
}

func (r *inviteRepository) Create(ctx context.Context, roomID, inviterID uuid.UUID, inviteeEmail, inviteeName, token string, expiresAt time.Time) error {
    query := `
        INSERT INTO invites (room_id, inviter_id, invitee_email, invitee_name, token, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query, roomID, inviterID, inviteeEmail, inviteeName, token, expiresAt)
    return err
}

func (r *inviteRepository) GetByToken(ctx context.Context, token string) (*model.Invite, error) {
    var invite model.Invite
    query := `SELECT * FROM invites WHERE token = $1`
    err := r.db.GetContext(ctx, &invite, query, token)
    return &invite, err
}

func (r *inviteRepository) MarkAsAccepted(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE invites SET status = 'accepted', accepted_at = NOW() WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
