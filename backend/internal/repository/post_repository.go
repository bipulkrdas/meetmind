package repository

import (
    "context"
    "livekit-consulting/backend/internal/model"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type PostRepository interface {
    Create(ctx context.Context, post *model.Post) error
    GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*model.PostWithCreator, error)
    Delete(ctx context.Context, postID uuid.UUID) error
}

type postRepository struct {
    db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
    return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
    query := `
        INSERT INTO posts (room_id, creator_id, message)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
    return r.db.QueryRowxContext(ctx, query, post.RoomID, post.CreatorID, post.Message).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *postRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) ([]*model.PostWithCreator, error) {
    var posts []*model.PostWithCreator
    query := `
        SELECT p.*, u.name as creator_name
        FROM posts p
        JOIN users u ON p.creator_id = u.id
        WHERE p.room_id = $1 AND p.is_deleted = false
        ORDER BY p.created_at DESC
    `
    err := r.db.SelectContext(ctx, &posts, query, roomID)
    return posts, err
}

func (r *postRepository) Delete(ctx context.Context, postID uuid.UUID) error {
    query := `UPDATE posts SET is_deleted = true, updated_at = NOW() WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, postID)
    return err
}
