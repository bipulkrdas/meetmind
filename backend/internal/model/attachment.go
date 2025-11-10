package model

import (
    "time"
    "github.com/google/uuid"
)

type Attachment struct {
    ID           uuid.UUID  `json:"id" db:"id"`
    MessageID    *uuid.UUID `json:"message_id" db:"message_id"`
    FileName     string     `json:"file_name" db:"file_name"`
    FileType     string     `json:"file_type" db:"file_type"`
    FileSize     int64      `json:"file_size" db:"file_size"`
    StoragePath  string     `json:"-" db:"storage_path"`
    StorageURL   string     `json:"storage_url" db:"storage_url"`
    ThumbnailURL *string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}
