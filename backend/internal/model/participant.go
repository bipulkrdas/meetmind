package model

import (
	"time"

	"github.com/google/uuid"
)

type RoomParticipant struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	RoomID          uuid.UUID  `json:"room_id" db:"room_id"`
	ParticipantID   *uuid.UUID `json:"participant_id" db:"participant_id"`
	UserID          *uuid.UUID `json:"user_id" db:"user_id"`
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	LiveKitIdentity *string    `json:"livekit_identity" db:"livekit_identity"`
	Role            string     `json:"role" db:"role"`
	JoinedAt        *time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	LastViewedAt    *time.Time `json:"last_viewed_at" db:"last_viewed_at"`
	LastReadSeqNo   int        `json:"last_read_seq_no" db:"last_read_seq_no"`
	IsActive        bool       `json:"is_active" db:"is_active"`
}

type AddParticipantRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2"`
}

type ParticipantInviteResponse struct {
	ParticipantID uuid.UUID `json:"participant_id"`
	InviteToken   string    `json:"invite_token"`
	InviteURL     string    `json:"invite_url"`
}
