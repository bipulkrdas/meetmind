package repository

import (
    "context"
    "database/sql"
    "livekit-consulting/backend/internal/model"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, name)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at, is_active
    `
    return r.db.QueryRowxContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Name).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    var user model.User
    query := `
        SELECT id, username, email, password_hash, name, created_at, updated_at, last_login, is_active
        FROM users
        WHERE id = $1 AND is_active = true
    `
    err := r.db.GetContext(ctx, &user, query, id)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    query := `
        SELECT id, username, email, password_hash, name, created_at, updated_at, last_login, is_active
        FROM users
        WHERE email = $1 AND is_active = true
    `
    err := r.db.GetContext(ctx, &user, query, email)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    query := `
        SELECT id, username, email, password_hash, name, created_at, updated_at, last_login, is_active
        FROM users
        WHERE username = $1 AND is_active = true
    `
    err := r.db.GetContext(ctx, &user, query, username)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
    query := `
        UPDATE users
        SET username = $1, email = $2, password_hash = $3, name = $4, updated_at = NOW()
        WHERE id = $5
    `
    _, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.Name, user.ID)
    return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
    query := `
        UPDATE users
        SET last_login = NOW()
        WHERE id = $1
    `
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}
