package repository

import (
	"context"
	"database/sql"
	"livekit-consulting/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ParticipantRepository interface {
	Create(ctx context.Context, participant *model.RoomParticipant) error
	GetByRoomAndEmail(ctx context.Context, roomID uuid.UUID, email string) (*model.RoomParticipant, error)
	GetByRoomAndUserID(ctx context.Context, roomID, userID uuid.UUID) (*model.RoomParticipant, error)
	CountByRoomID(ctx context.Context, roomID uuid.UUID) (int, error)
	UserHasAccess(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*model.RoomParticipant, error)
	Delete(ctx context.Context, participantID uuid.UUID) error
}

type participantRepository struct {
	db *sqlx.DB
}

func NewParticipantRepository(db *sqlx.DB) ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) Create(ctx context.Context, participant *model.RoomParticipant) error {
	query := `
		INSERT INTO room_participants (room_id, user_id, email, name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, is_active
	`
	return r.db.QueryRowxContext(ctx, query, participant.RoomID, participant.UserID, participant.Email, participant.Name, participant.Role).Scan(&participant.ID, &participant.CreatedAt, &participant.IsActive)
}

func (r *participantRepository) GetByRoomAndEmail(ctx context.Context, roomID uuid.UUID, email string) (*model.RoomParticipant, error) {
	var participant model.RoomParticipant
	query := `SELECT * FROM room_participants WHERE room_id = $1 AND email = $2`
	err := r.db.GetContext(ctx, &participant, query, roomID, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &participant, nil
}

func (r *participantRepository) GetByRoomAndUserID(ctx context.Context, roomID, userID uuid.UUID) (*model.RoomParticipant, error) {
	var participant model.RoomParticipant
	query := `SELECT * FROM room_participants WHERE room_id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &participant, query, roomID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &participant, nil
}

func (r *participantRepository) CountByRoomID(ctx context.Context, roomID uuid.UUID) (int, error) {
    var count int
    query := `SELECT count(*) FROM room_participants WHERE room_id = $1 AND is_active = true`
    err := r.db.GetContext(ctx, &count, query, roomID)
    return count, err
}

func (r *participantRepository) UserHasAccess(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
    var exists bool
    query := `SELECT exists(SELECT 1 FROM room_participants WHERE room_id = $1 AND user_id = $2 AND is_active = true)`
    err := r.db.GetContext(ctx, &exists, query, roomID, userID)
    return exists, err
}

func (r *participantRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*model.RoomParticipant, error) {
    var participants []*model.RoomParticipant
    query := `SELECT * FROM room_participants WHERE room_id = $1 AND is_active = true`
    err := r.db.SelectContext(ctx, &participants, query, roomID)
    return participants, err
}

func (r *participantRepository) Delete(ctx context.Context, participantID uuid.UUID) error {
    query := `UPDATE room_participants SET is_active = false WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, participantID)
    return err
}
