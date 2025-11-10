package repository

import (
	"context"
	"database/sql"
	"livekit-consulting/backend/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
	GetByName(ctx context.Context, name string) (*model.Room, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*model.Room, error)
	Update(ctx context.Context, room *model.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetRoomsByUser(ctx context.Context, userID uuid.UUID) ([]*model.Room, error)
	UpdateLastRead(ctx context.Context, roomID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int, error)
}

type roomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
	query := `
        INSERT INTO rooms (room_name, description, owner_id, livekit_room_name, room_sid)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at, is_active
    `
	return r.db.QueryRowxContext(ctx, query, room.RoomName, room.Description, room.OwnerID, room.LiveKitRoomName, room.RoomSID).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt, &room.IsActive)
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	var room model.Room
	query := `
        SELECT id, room_name, room_sid, description, owner_id, livekit_room_name, 
               metadata, created_at, updated_at, is_active, last_message_seq, last_message_at
        FROM rooms
        WHERE id = $1 AND is_active = true
    `
	err := r.db.GetContext(ctx, &room, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &room, err
}

func (r *roomRepository) GetByName(ctx context.Context, name string) (*model.Room, error) {
	var room model.Room
	query := `
        SELECT id, room_name, room_sid, description, owner_id, livekit_room_name, 
               metadata, created_at, updated_at, is_active, last_message_seq, last_message_at
        FROM rooms
        WHERE room_name = $1 AND is_active = true
    `
	err := r.db.GetContext(ctx, &room, query, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &room, err
}

func (r *roomRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*model.Room, error) {
	var rooms []*model.Room
	query := `
        SELECT id, room_name, room_sid, description, owner_id, livekit_room_name, 
               metadata, created_at, updated_at, is_active, last_message_seq, last_message_at
        FROM rooms
        WHERE owner_id = $1 AND is_active = true
        ORDER BY created_at DESC
    `
	err := r.db.SelectContext(ctx, &rooms, query, ownerID)
	return rooms, err
}

func (r *roomRepository) Update(ctx context.Context, room *model.Room) error {
	query := `
        UPDATE rooms
        SET room_name = $1, description = $2, livekit_room_name = $3, room_sid = $4, updated_at = NOW()
        WHERE id = $5
    `
	_, err := r.db.ExecContext(ctx, query, room.RoomName, room.Description, room.LiveKitRoomName, room.RoomSID, room.ID)
	return err
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE rooms SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *roomRepository) GetRoomsByUser(ctx context.Context, userID uuid.UUID) ([]*model.Room, error) {
	var rooms []*model.Room
	query := `
        SELECT r.id, r.room_name, r.room_sid, r.description, r.owner_id, r.livekit_room_name, 
               r.metadata, r.created_at, r.updated_at, r.is_active, r.last_message_seq, r.last_message_at
        FROM rooms r
        JOIN room_participants rp ON r.id = rp.room_id
        WHERE rp.user_id = $1 AND r.is_active = true
        ORDER BY r.created_at DESC
    `
	err := r.db.SelectContext(ctx, &rooms, query, userID)
	return rooms, err
}

func (r *roomRepository) UpdateLastRead(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `
        UPDATE room_participants 
        SET last_viewed_at = $1
        WHERE room_id = $2 AND user_id = $3
    `

	_, err := r.db.ExecContext(ctx, query, time.Now(), roomID, userID)
	return err
}

func (r *roomRepository) GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM messages m
        INNER JOIN room_participants rp ON m.room_id = rp.room_id
        WHERE m.room_id = $1 
              AND rp.user_id = $2 
              AND m.created_at > rp.last_viewed_at
              AND m.deleted_at IS NULL
    `

	var count int
	err := r.db.GetContext(ctx, &count, query, roomID, userID)
	return count, err
}