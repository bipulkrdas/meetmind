package repository

import (
	"context"
	"database/sql"
	"time"

	"livekit-consulting/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) (*model.Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Message, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID, limit int, before *uuid.UUID) ([]*model.Message, error)
	Update(ctx context.Context, id uuid.UUID, content string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, roomID uuid.UUID, searchTerm string, limit int) ([]*model.Message, error)
	GetMessageWithAttachments(ctx context.Context, id uuid.UUID) (*model.Message, error)
	UpdateMetadata(ctx context.Context, id uuid.UUID, metadata *model.MessageMetadata) error
}

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *model.Message) (*model.Room, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var room model.Room
	err = tx.GetContext(ctx, &room, "UPDATE rooms SET last_message_seq = last_message_seq + 1, last_message_at = NOW() WHERE id = $1 RETURNING *", message.RoomID)
	if err != nil {
		return nil, err
	}

	message.SeqNo = room.LastMessageSeq

	query := `
        INSERT INTO messages (id, room_id, user_id, seq_no, content, message_type, metadata, extra_data, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at
    `
	message.ID = uuid.New()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	if message.MessageType == "" {
		message.MessageType = model.MessageTypeUserMessage
	}

	err = tx.QueryRowxContext(
		ctx,
		query,
		message.ID,
		message.RoomID,
		message.UserID,
		message.SeqNo,
		message.Content,
		message.MessageType,
		message.Metadata,
		message.ExtraData,
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &room, tx.Commit()
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	query := `
        SELECT m.id, m.room_id, m.user_id, u.name as username, m.seq_no, m.content, m.message_type, m.extra_data, m.metadata,
               m.edited, m.created_at, m.updated_at, m.deleted_at
        FROM messages m
        LEFT JOIN users u ON m.user_id = u.id
        WHERE m.id = $1 AND m.deleted_at IS NULL
    `

	var message model.Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *messageRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID, limit int, before *uuid.UUID) ([]*model.Message, error) {
	var messages []model.Message
	var err error

	if before != nil {
		query := `
			SELECT m.id, m.room_id, m.user_id, u.name as username, m.seq_no, m.content, m.message_type, m.extra_data, m.metadata, 
				   m.edited, m.created_at, m.updated_at
			FROM messages m
			LEFT JOIN users u ON m.user_id = u.id
			WHERE m.room_id = $1 AND m.deleted_at IS NULL 
				  AND m.created_at < (SELECT created_at FROM messages WHERE id = $2)
			ORDER BY m.created_at DESC
			LIMIT $3
		`
		err = r.db.SelectContext(ctx, &messages, query, roomID, *before, limit)
	} else {
		query := `
			SELECT m.id, m.room_id, m.user_id, u.name as username, m.seq_no, m.content, m.message_type, m.extra_data, m.metadata, 
				   m.edited, m.created_at, m.updated_at
			FROM messages m
			LEFT JOIN users u ON m.user_id = u.id
			WHERE m.room_id = $1 AND m.deleted_at IS NULL
			ORDER BY m.created_at DESC
			LIMIT $2
		`
		err = r.db.SelectContext(ctx, &messages, query, roomID, limit)
	}

	if err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	messagePtrs := make([]*model.Message, len(messages))
	for i := range messages {
		messagePtrs[i] = &messages[i]
	}

	return messagePtrs, nil
}

func (r *messageRepository) Update(ctx context.Context, id uuid.UUID, content string) error {
	query := `
        UPDATE messages
        SET content = $1, edited = true, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, content, time.Now(), id)
	return err
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE messages
        SET deleted_at = $1
        WHERE id = $2
    `

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *messageRepository) Search(ctx context.Context, roomID uuid.UUID, searchTerm string, limit int) ([]*model.Message, error) {
	query := `
        SELECT m.id, m.room_id, m.user_id, u.name as username, m.seq_no, m.content, m.message_type, m.extra_data, m.metadata,
               m.edited, m.created_at, m.updated_at
        FROM messages m
        LEFT JOIN users u ON m.user_id = u.id
        WHERE m.room_id = $1
              AND m.deleted_at IS NULL
              AND to_tsvector('english', m.content) @@ plainto_tsquery('english', $2)
        ORDER BY m.created_at DESC
        LIMIT $3
    `

	var messages []model.Message
	err := r.db.SelectContext(ctx, &messages, query, roomID, searchTerm, limit)
	if err != nil {
		return nil, err
	}

	messagePtrs := make([]*model.Message, len(messages))
	for i := range messages {
		messagePtrs[i] = &messages[i]
	}

	return messagePtrs, nil
}

func (r *messageRepository) GetMessageWithAttachments(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	message, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	attachmentQuery := `
        SELECT id, message_id, file_name, file_type, file_size,
               storage_url, thumbnail_url, created_at
        FROM attachments
        WHERE message_id = $1
        ORDER BY created_at
    `

	var attachments []model.Attachment
	err = r.db.SelectContext(ctx, &attachments, attachmentQuery, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	message.Attachments = attachments
	return message, nil
}

func (r *messageRepository) UpdateMetadata(ctx context.Context, id uuid.UUID, metadata *model.MessageMetadata) error {
	query := `
        UPDATE messages
        SET metadata = $1, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, metadata, time.Now(), id)
	return err
}