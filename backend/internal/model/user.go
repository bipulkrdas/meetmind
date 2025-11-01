package model

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID  `json:"id" db:"id"`
    Username     string     `json:"username" db:"username"`
    Email        string     `json:"email" db:"email"`
    PasswordHash string     `json:"-" db:"password_hash"`
    Name         string     `json:"name" db:"name"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
    LastLogin    *time.Time `json:"last_login" db:"last_login"`
    IsActive     bool       `json:"is_active" db:"is_active"`
}

type UserSignUpRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2"`
    Password string `json:"password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserSignInRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token        string    `json:"token"`
    User         User      `json:"user"`
    LiveKitToken string    `json:"livekit_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type PasswordResetRequest struct {
    Email string `json:"email" validate:"required,email"`
}

type PasswordResetConfirm struct {
    Token       string `json:"token" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=8"`
}
