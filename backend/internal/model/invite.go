package model

import (
    "time"
    "github.com/google/uuid"
)

type Invite struct {
    ID           uuid.UUID `json:"id" db:"id"`
    RoomID       uuid.UUID `json:"room_id" db:"room_id"`
    InviterID    uuid.UUID `json:"inviter_id" db:"inviter_id"`
    InviteeEmail string    `json:"invitee_email" db:"invitee_email"`
    InviteeName  string    `json:"invitee_name" db:"invitee_name"`
    Token        string    `json:"token" db:"token"`
    Status       string    `json:"status" db:"status"`
    ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    AcceptedAt   *time.Time `json:"accepted_at" db:"accepted_at"`
}
