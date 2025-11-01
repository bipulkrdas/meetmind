package model

import (
    "time"
    "github.com/google/uuid"
)

type Post struct {
    ID        uuid.UUID `json:"id" db:"id"`
    RoomID    uuid.UUID `json:"room_id" db:"room_id"`
    CreatorID uuid.UUID `json:"creator_id" db:"creator_id"`
    Message   string    `json:"message" db:"message"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
}

type CreatePostRequest struct {
    Message string `json:"message" validate:"required,min=1,max=5000"`
}

type PostWithCreator struct {
    Post
    CreatorName string `json:"creator_name"`
}
