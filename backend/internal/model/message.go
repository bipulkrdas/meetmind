package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of a message.
type MessageType string

const (
	// MessageTypeUserMessage is a standard message from a user.
	MessageTypeUserMessage MessageType = "user_message"
	// MessageTypeMeetingTranscript indicates the message contains a meeting transcript.
	MessageTypeMeetingTranscript MessageType = "meeting_transcript"
	// MessageTypeParticipantJoined indicates a participant has joined the room.
	MessageTypeParticipantJoined MessageType = "participant_joined"
)

type Message struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	RoomID      uuid.UUID        `json:"room_id" db:"room_id"`
	UserID      *uuid.UUID       `json:"user_id" db:"user_id"`
	Username    string           `json:"username" db:"username"` // Joined from users table
	SeqNo       int              `json:"seq_no" db:"seq_no"`
	Content     string           `json:"content" db:"content"`
	MessageType MessageType      `json:"message_type" db:"message_type"`
	Metadata    *MessageMetadata `json:"metadata" db:"metadata"`
	ExtraData   *ExtraData       `json:"extra_data,omitempty" db:"extra_data"`
	Edited      bool             `json:"edited" db:"edited"`
	Attachments []Attachment     `json:"attachments,omitempty"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ExtraData holds flexible JSON data for special message types.
type ExtraData struct {
	Transcript *TranscriptData `json:"transcript,omitempty"`
}

// Scan implements the sql.Scanner interface for ExtraData.
func (e *ExtraData) Scan(value interface{}) error {
	if value == nil {
		*e = ExtraData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte for ExtraData")
	}
	if len(bytes) == 0 {
		*e = ExtraData{}
		return nil
	}
	return json.Unmarshal(bytes, e)
}

// Value implements the driver.Valuer interface for ExtraData.
func (e ExtraData) Value() (driver.Value, error) {
	// If Transcript is nil, we should store a null value in the DB.
	if e.Transcript == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

// TranscriptData holds information about a meeting transcript.
type TranscriptData struct {
	Bucket       string    `json:"bucket"`
	Region       string    `json:"region"`
	S3Keys       S3Keys    `json:"s3_keys"`
	HTTPSUrls    HTTPSUrls `json:"https_urls"`
	SessionStart time.Time `json:"session_start"`
	SessionEnd   time.Time `json:"session_end"`
}

// S3Keys holds the S3 object keys for the transcript files.
type S3Keys struct {
	JSON string `json:"json"`
	Text string `json:"text"`
}

// HTTPSUrls holds the public HTTPS URLs for the transcript files.
type HTTPSUrls struct {
	JSON string `json:"json_https_url"`
	Text string `json:"text_https_url"`
}

type MessageMetadata struct {
	Reactions []Reaction `json:"reactions,omitempty"`
	Mentions  []string   `json:"mentions,omitempty"`
}

func (m *MessageMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = MessageMetadata{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	if len(bytes) == 0 {
		*m = MessageMetadata{}
		return nil
	}
	return json.Unmarshal(bytes, m)
}

func (m MessageMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type Reaction struct {
	Emoji   string      `json:"emoji"`
	UserIDs []uuid.UUID `json:"user_ids"`
	Count   int         `json:"count"`
}

type CreateMessageRequest struct {
	Content       string      `json:"content" validate:"required,max=5000"`
	AttachmentIDs []uuid.UUID `json:"attachment_ids,omitempty"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required,max=5000"`
}

type UpdateLastReadRequest struct {
	LastReadSequenceNumber int `json:"last_read_sequence_number"`
}