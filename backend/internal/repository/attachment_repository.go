package repository

import (
	"context"

	"livekit-consulting/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *model.Attachment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetMessageID(ctx context.Context, attachmentID, messageID uuid.UUID) error
}

type attachmentRepository struct {
	db *sqlx.DB
}

func NewAttachmentRepository(db *sqlx.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *model.Attachment) error {
	query := `
		INSERT INTO attachments (id, file_name, file_type, file_size, storage_path, storage_url, thumbnail_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	return r.db.QueryRowxContext(ctx, query,
		attachment.ID,
		attachment.FileName,
		attachment.FileType,
		attachment.FileSize,
		attachment.StoragePath,
		attachment.StorageURL,
		attachment.ThumbnailURL,
	).Scan(&attachment.CreatedAt)
}

func (r *attachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	var attachment model.Attachment
	query := `SELECT * FROM attachments WHERE id = $1`
	err := r.db.GetContext(ctx, &attachment, query, id)
	return &attachment, err
}

func (r *attachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM attachments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *attachmentRepository) SetMessageID(ctx context.Context, attachmentID, messageID uuid.UUID) error {
	query := `UPDATE attachments SET message_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, messageID, attachmentID)
	return err
}
