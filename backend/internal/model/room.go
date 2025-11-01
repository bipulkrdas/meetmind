package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Metadata map[string]interface{}

func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	if len(bytes) == 0 {
		*m = nil
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

type Room struct {
	ID              uuid.UUID `json:"id" db:"id"`
	RoomName        string    `json:"room_name" db:"room_name"`
	RoomSID         *string   `json:"room_sid" db:"room_sid"`
	Description     *string   `json:"description" db:"description"`
	OwnerID         uuid.UUID `json:"owner_id" db:"owner_id"`
	LiveKitRoomName *string   `json:"livekit_room_name" db:"livekit_room_name"`
	Metadata        Metadata  `json:"metadata" db:"metadata"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	IsActive        bool      `json:"is_active" db:"is_active"`
}

type CreateRoomRequest struct {
    RoomName    string  `json:"room_name" validate:"required,min=3,max=100"`
    Description *string `json:"description"`
}

type RoomResponse struct {
    Room             Room                `json:"room"`
    ParticipantCount int                 `json:"participant_count"`
    IsOwner          bool                `json:"is_owner"`
}
