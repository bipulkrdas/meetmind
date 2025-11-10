

## RoomParticipant
### Roomparticipant Model struct
type RoomParticipant struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	RoomID          uuid.UUID  `json:"room_id" db:"room_id"`
	ParticipantID   *uuid.UUID `json:"participant_id" db:"participant_id"`
	UserID          *uuid.UUID `json:"user_id" db:"user_id"`
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	LiveKitIdentity *string    `json:"livekit_identity" db:"livekit_identity"`
	Role            string     `json:"role" db:"role"`
	JoinedAt        *time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	LastViewedAt    time.Time  `json:"last_viewed_at" db:"last_viewed_at"`
	LastReadSeqNo   int        `json:"last_read_seq_no" db:"last_read_seq_no"`
	IsActive        bool       `json:"is_active" db:"is_active"`
}

### RoomParticipant database table schema
CREATE TABLE room_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    participant_id UUID, -- Can be NULL for external participants
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    livekit_identity VARCHAR(255), -- LiveKit participant identity
    role VARCHAR(50) DEFAULT 'participant', -- owner, moderator, participant
    joined_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at TIMESTAMP
    last_read_seq_no INTEGER DEFAULT 0 
    is_active BOOLEAN DEFAULT true,
    UNIQUE(room_id, email)
);

### Attachment Model

go
// internal/models/attachment.go
package models

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

### Attachments Table database schema

```sql
CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    storage_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_attachments_message_id ON attachments(message_id);
```


### Message Model struct
type Message struct {
    ID          uuid.UUID       `json:"id" db:"id"`
    RoomID      uuid.UUID       `json:"room_id" db:"room_id"`
    UserID      *uuid.UUID      `json:"user_id" db:"user_id"`
    Username    string          `json:"username" db:"username"` // Joined from users table
    SeqNon int `json:"seq_no" db:"seq_no"`
    Content     string          `json:"content" db:"content"`
    Metadata    *MessageMetadata `json:"metadata,omitempty" db:"metadata"`
    Edited      bool            `json:"edited" db:"edited"`
    Attachments []Attachment    `json:"attachments,omitempty"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
}

type MessageMetadata struct {
    Reactions []Reaction `json:"reactions,omitempty"`
    Mentions  []string   `json:"mentions,omitempty"`
}

type Reaction struct {
    Emoji   string      `json:"emoji"`
    UserIDs []uuid.UUID `json:"user_ids"`
    Count   int         `json:"count"`
}

type CreateMessageRequest struct {
    Content       string      `json:"content" binding:"required,max=5000"`
    AttachmentIDs []uuid.UUID `json:"attachment_ids,omitempty"`
}

type UpdateMessageRequest struct {
    Content string `json:"content" binding:"required,max=5000"`
}

### messages DB table schema
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    seq_no INTEGER DEFAULT 0,
    content TEXT NOT NULL,
    metadata JSONB, -- For reactions, mentions, etc.
    edited BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_read_at TIMESTAMP WITH TIME ZONE
    deleted_at TIMESTAMP WITH TIME ZONE
);

## API Routes at @backend/cmd/api/main.go

authAPI.HandleFunc("/rooms/{roomId}/messages", postHandler.CreateMessage).Methods("POST")
authAPI.HandleFunc("/rooms/{roomId}/messages", postHandler.GetPosts).Methods("GET")
authAPI.HandleFunc("/rooms/{roomId}/update_last_read_for_user", postHandler.updateLastRead).Methods("POST")
[this will receive a json {"last_read_sequence_number": <number>} and update the "room_participant" table for the following columns:     last_viewed_at TIMESTAMP
    last_read_seq_no INTEGER DEFAULT 0, where last_viewed_at will be current timestamp and last_read_seq_no will be the value from JSON request.]


## Repository layer - Message and Attachments
## 8. Repository Layer

### 8.1 Message Repository

```go
// internal/repository/message_repo.go
package repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "your-project/internal/models"
)

type MessageRepository struct {
    db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
    return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {


    USE the following logic for creating message(in the repository layer)


   First update:
       // when a new message is created, we need to update the "room" table for the following columns
    //   last_message_seq INTEGER DEFAULT 0,>> this will be the 
    // last_message_at TIMESTAMP, 
   
   // Here make a transaction
    // 1. First update room table with last_msg_seq = last_msg_seq + 1
    // 2. Update  last_message_at with the current timestamp.

    get the last_message_seq from database and assign it to a variable:
    lastMessageSeq = "last_essage_seq"

    In the transaction , next part of the transaction, create the message record in the "message" table using the following query.
    
    query := `" 
        INSERT INTO messages (id, room_id, user_id, seq_no, content, metadata, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `

    message.ID = uuid.New()
    message.CreatedAt = time.Now()
    message.UpdatedAt = time.Now()
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        message.ID,
        message.RoomID,
        message.UserID,
        lastMessageSeq
        message.Content,
        message.Metadata,
        message.CreatedAt,
        message.UpdatedAt,
    ).Scan(&message.ID, &message.CreatedAt)
    
    return err
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
    query := `
        SELECT m.id, m.room_id, m.user_id, u.username, m.content, m.metadata, 
               m.edited, m.created_at, m.updated_at, m.deleted_at
        FROM messages m
        LEFT JOIN users u ON m.user_id = u.id
        WHERE m.id = $1 AND m.deleted_at IS NULL
    `
    
    var message models.Message
    err := r.db.GetContext(ctx, &message, query, id)
    if err != nil {
        return nil, err
    }
    
    return &message, nil
}

func (r *MessageRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID, limit int, before *uuid.UUID) ([]models.Message, error) {
    var query string
    var args []interface{}
    
    if before != nil {
        query = `
            SELECT m.id, m.room_id, m.user_id, u.username, m.content, m.metadata, 
                   m.edited, m.created_at, m.updated_at
            FROM messages m
            LEFT JOIN users u ON m.user_id = u.id
            WHERE m.room_id = $1 AND m.deleted_at IS NULL 
                  AND m.created_at < (SELECT created_at FROM messages WHERE id = $2)
            ORDER BY m.created_at DESC
            LIMIT $3
        `
        args = []interface{}{roomID, before, limit}
    } else {
        query = `
            SELECT m.id, m.room_id, m.user_id, u.username, m.content, m.metadata, 
                   m.edited, m.created_at, m.updated_at
            FROM messages m
            LEFT JOIN users u ON m.user_id = u.id
            WHERE m.room_id = $1 AND m.deleted_at IS NULL
            ORDER BY m.created_at DESC
            LIMIT $2
        `
        args = []interface{}{roomID, limit}
    }
    
    var messages []models.Message
    err := r.db.SelectContext(ctx, &messages, query, args...)
    if err != nil {
        return nil, err
    }
    
    // Reverse to get chronological order
    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }
    
    return messages, nil
}

func (r *MessageRepository) Update(ctx context.Context, id uuid.UUID, content string) error {
    query := `
        UPDATE messages 
        SET content = $1, edited = true, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `
    
    _, err := r.db.ExecContext(ctx, query, content, time.Now(), id)
    return err
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE messages 
        SET deleted_at = $1
        WHERE id = $2
    `
    
    _, err := r.db.ExecContext(ctx, query, time.Now(), id)
    return err
}

func (r *MessageRepository) Search(ctx context.Context, roomID uuid.UUID, searchTerm string, limit int) ([]models.Message, error) {
    query := `
        SELECT m.id, m.room_id, m.user_id, u.username, m.content, m.metadata, 
               m.edited, m.created_at, m.updated_at
        FROM messages m
        LEFT JOIN users u ON m.user_id = u.id
        WHERE m.room_id = $1 
              AND m.deleted_at IS NULL
              AND to_tsvector('english', m.content) @@ plainto_tsquery('english', $2)
        ORDER BY m.created_at DESC
        LIMIT $3
    `
    
    var messages []models.Message
    err := r.db.SelectContext(ctx, &messages, query, roomID, searchTerm, limit)
    return messages, err
}

func (r *MessageRepository) GetMessageWithAttachments(ctx context.Context, id uuid.UUID) (*models.Message, error) {
    message, err := r.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Get attachments
    attachmentQuery := `
        SELECT id, message_id, file_name, file_type, file_size, 
               storage_url, thumbnail_url, created_at
        FROM attachments
        WHERE message_id = $1
        ORDER BY created_at
    `
    
    var attachments []models.Attachment
    err = r.db.SelectContext(ctx, &attachments, attachmentQuery, id)
    if err != nil && err != sql.ErrNoRows {
        return nil, err
    }
    
    message.Attachments = attachments
    return message, nil
}

func (r *MessageRepository) UpdateMetadata(ctx context.Context, id uuid.UUID, metadata *models.MessageMetadata) error {
    query := `
        UPDATE messages 
        SET metadata = $1, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `
    
    _, err := r.db.ExecContext(ctx, query, metadata, time.Now(), id)
    return err
}
`

### Room Repository (add these functions to the already existing functins)
func (r *RoomRepository) UpdateLastRead(ctx context.Context, roomID, userID uuid.UUID) error {
    query := `
        UPDATE room_members 
        SET last_read_at = $1
        WHERE room_id = $2 AND user_id = $3
    `
    
    _, err := r.db.ExecContext(ctx, query, time.Now(), roomID, userID)
    return err
}

func (r *RoomRepository) GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
    query := `
        SELECT COUNT(*)
        FROM messages m
        INNER JOIN room_members rm ON m.room_id = rm.room_id
        WHERE m.room_id = $1 
              AND rm.user_id = $2 
              AND m.created_at > rm.last_read_at
              AND m.deleted_at IS NULL
    `
    
    var count int
    err := r.db.GetContext(ctx, &count, query, roomID, userID)
    return count, err
}
```


## 9. Service Layer

## 14. File Upload Service
### Here is a sample implementation using minio. Make the File service interface based and implament AWS S3, GCS (Google cloud storage) and Minio, all three of them. I  production, we can use any one of these service for file storage. We want to be flexible and provide option to our users on which storage to choose.
```go
// internal/service/file_service.go
package service

import (
    "context"
    "fmt"
    "io"
    "mime/multipart"
    "path/filepath"
    "time"

    "github.com/google/uuid"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "your-project/internal/config"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type FileService struct {
    minioClient *minio.Client
    bucket      string
    attachRepo  *repository.AttachmentRepository
}

func NewFileService(cfg config.StorageConfig, attachRepo *repository.AttachmentRepository) (*FileService, error) {
    minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
        Secure: false, // Set to true for production with HTTPS
    })
    if err != nil {
        return nil, err
    }
    
    return &FileService{
        minioClient: minioClient,
        bucket:      cfg.Bucket,
        attachRepo:  attachRepo,
    }, nil
}

func (s *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader, roomID uuid.UUID) (*models.Attachment, error) {
    // Open file
    src, err := file.Open()
    if err != nil {
        return nil, err
    }
    defer src.Close()
    
    // Generate unique filename
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("%s/%s%s", roomID.String(), uuid.New().String(), ext)
    
    // Upload to MinIO
    _, err = s.minioClient.PutObject(
        ctx,
        s.bucket,
        filename,
        src,
        file.Size,
        minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
    )
    if err != nil {
        return nil, err
    }
    
    // Generate presigned URL (valid for 7 days)
    url, err := s.minioClient.PresignedGetObject(ctx, s.bucket, filename, 7*24*time.Hour, nil)
    if err != nil {
        return nil, err
    }
    
    // Create attachment record
    attachment := &models.Attachment{
        ID:          uuid.New(),
        FileName:    file.Filename,
        FileType:    file.Header.Get("Content-Type"),
        FileSize:    file.Size,
        StoragePath: filename,
        StorageURL:  url.String(),
        CreatedAt:   time.Now(),
    }
    
    if err := s.attachRepo.Create(ctx, attachment); err != nil {
        return nil, err
    }
    
    return attachment, nil
}

func (s *FileService) DeleteFile(ctx context.Context, attachmentID uuid.UUID) error {
    attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
    if err != nil {
        return err
    }
    
    // Delete from MinIO
    err = s.minioClient.RemoveObject(ctx, s.bucket, attachment.StoragePath, minio.RemoveObjectOptions{})
    if err != nil {
        return err
    }
    
    // Delete from database
    return s.attachRepo.Delete(ctx, attachmentID)
}

func (s *FileService) GetFile(ctx context.Context, attachmentID uuid.UUID) (io.Reader, error) {
    attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
    if err != nil {
        return nil, err
    }
    
    obj, err := s.minioClient.GetObject(ctx, s.bucket, attachment.StoragePath, minio.GetObjectOptions{})
    if err != nil {
        return nil, err
    }
    
    return obj, nil
}
```


### 9.2 Message Service

```go
// internal/service/message_service.go
package service

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "your-project/internal/models"
    "your-project/internal/repository"
    "your-project/internal/websocket"
)

type MessageService struct {
    messageRepo *repository.MessageRepository
    roomRepo    *repository.RoomRepository
  
}

func NewMessageService(
    messageRepo *repository.MessageRepository,
    roomRepo *repository.RoomRepository,
) *MessageService {
    return &MessageService{
        messageRepo: messageRepo,
        roomRepo:    roomRepo,
    }
}

func (s *MessageService) CreateMessage(ctx context.Context, req *models.CreateMessageRequest, roomID, userID uuid.UUID) (*models.Message, error) {
    // Check if user is a member of the room
    isMember, err := s.roomRepo.IsMember(ctx, roomID, userID)
    if err != nil {
        return nil, err
    }
    if !isMember {
        return nil, errors.New("user is not a member of this room")
    }
    
    // Create message
    message := &models.Message{
        RoomID:  roomID,
        UserID:  &userID,
        Content: req.Content,
    }
    
    if err := s.messageRepo.Create(ctx, message); err != nil {
        return nil, err
    }
    
    // Get full message with user info
    fullMessage, err := s.messageRepo.GetByID(ctx, message.ID)
    if err != nil {
        return nil, err
    }

    
    return fullMessage, nil
}

func (s *MessageService) GetMessages(ctx context.Context, roomID, userID uuid.UUID, limit int, before *uuid.UUID) ([]models.Message, error) {
    // Check if user is a member
    isMember, err := s.roomRepo.IsMember(ctx, roomID, userID)
    if err != nil {
        return nil, err
    }
    if !isMember {
        return nil, errors.New("user is not a member of this room")
    }
    
    messages, err := s.messageRepo.GetByRoomID(ctx, roomID, limit, before)
    if err != nil {
        return nil, err
    }
    
    
    return messages, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, messageID, userID uuid.UUID, content string) error {
    // Get message to verify ownership
    message, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return err
    }
    
    if message.UserID == nil || *message.UserID != userID {
        return errors.New("unauthorized to edit this message")
    }
    
    if err := s.messageRepo.Update(ctx, messageID, content); err != nil {
        return err
    }
    
    // Broadcast update
    updatedMessage, _ := s.messageRepo.GetByID(ctx, messageID)
    s.hub.BroadcastToRoom(message.RoomID, map[string]interface{}{
        "type":       "message_updated",
        "message_id": messageID,
        "message":    updatedMessage,
    }, nil)
    
    return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID, userID uuid.UUID) error {
    message, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return err
    }
    
    if message.UserID == nil || *message.UserID != userID {
        return errors.New("unauthorized to delete this message")
    }
    
    if err := s.messageRepo.Delete(ctx, messageID); err != nil {
        return err
    }
    
    // Broadcast deletion
    s.hub.BroadcastToRoom(message.RoomID, map[string]interface{}{
        "type":       "message_deleted",
        "message_id": messageID,
    }, nil)
    
    return nil
}

func (s *MessageService) SearchMessages(ctx context.Context, roomID, userID uuid.UUID, query string, limit int) ([]models.Message, error) {
    isMember, err := s.roomRepo.IsMember(ctx, roomID, userID)
    if err != nil {
        return nil, err
    }
    if !isMember {
        return nil, errors.New("user is not a member of this room")
    }
    
    return s.messageRepo.Search(ctx, roomID, query, limit)
}

func (s *MessageService) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
    message, err := s.messageRepo.GetByID(ctx, messageID)
    if err != nil {
        return err
    }
    
    // Initialize metadata if nil
    if message.Metadata == nil {
        message.Metadata = &models.MessageMetadata{
            Reactions: []models.Reaction{},
        }
    }
    
    // Find or create reaction
    found := false
    for i, reaction := range message.Metadata.Reactions {
        if reaction.Emoji == emoji {
            // Check if user already reacted
            for _, id := range reaction.UserIDs {
                if id == userID {
                    return errors.New("user already reacted with this emoji")
                }
            }
            message.Metadata.Reactions[i].UserIDs = append(reaction.UserIDs, userID)
            message.Metadata.Reactions[i].Count++
            found = true
            break
        }
    }
    
    if !found {
        message.Metadata.Reactions = append(message.Metadata.Reactions, models.Reaction{
            Emoji:   emoji,
            UserIDs: []uuid.UUID{userID},
            Count:   1,
        })
    }
    
    // Update metadata
    if err := s.messageRepo.UpdateMetadata(ctx, messageID, message.Metadata); err != nil {
        return err
    }
    
    // Broadcast update
    updatedMessage, _ := s.messageRepo.GetByID(ctx, messageID)
    
    return nil
}
```
## File Upload Handler


Route:  authAPI.HandleFunc("/rooms/{roomId}/attachments"

response will have the fileId.

 write the handler to satisfy the the following request and requirements from UI client (webapp):

 export async function uploadFiles(
  roomId: string,
  file: File,
  onProgress?: (progress: number) => void
): Promise<string> {
  return new Promise((resolve, reject) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('roomId', roomId);

    const xhr = new XMLHttpRequest();

    // Track upload progress
    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable && onProgress) {
        const progress = Math.round((e.loaded / e.total) * 100);
        onProgress(progress);
      }
    });

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const response = JSON.parse(xhr.responseText);
          resolve(response.fileId);
        } catch (error) {
          reject(new Error('Invalid response from server'));
        }
      } else {
        reject(new Error(`Upload failed with status ${xhr.status}`));
      }
    });

    xhr.addEventListener('error', () => {
      reject(new Error('Network error during upload'));
    });

    xhr.addEventListener('abort', () => {
      reject(new Error('Upload aborted'));
    });

    xhr.open('POST', `${API_BASE}/rooms/{roomId}/attachments`);
    xhr.setRequestHeader('Authorization', `Bearer ${getAuthToken()}`);
    xhr.send(formData);
  });
}


## Message Handler

```go
// internal/handlers/message_handler.go
package handler

type MessageHandler struct {
    messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
    return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
    roomID, err := uuid.Parse(c.Param("roomId"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid room ID")
        return
    }
    
    userID := c.GetString("user_id") // From JWT middleware
    uid, _ := uuid.Parse(userID)
    
    var req models.CreateMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    message, err := h.messageService.CreateMessage(c.Request.Context(), &req, roomID, uid)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusCreated, message)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
    roomID, err := uuid.Parse(c.Param("roomId"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid room ID")
        return
    }
    
    userID := c.GetString("user_id")
    uid, _ := uuid.Parse(userID)
    
    limit := 20
    if l := c.Query("limit"); l != "" {
        // Parse limit
    }
    
    var before *uuid.UUID
    if b := c.Query("before"); b != "" {
        beforeID, err := uuid.Parse(b)
        if err == nil {
            before = &beforeID
        }
    }
    
    messages, err := h.messageService.GetMessages(c.Request.Context(), roomID, uid, limit, before)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, messages)
}

func (h *MessageHandler) UpdateMessage(c *gin.Context) {
    messageID, err := uuid.Parse(c.Param("messageId"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid message ID")
        return
    }
    
    userID := c.GetString("user_id")
    uid, _ := uuid.Parse(userID)
    
    var req models.UpdateMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    if err := h.messageService.UpdateMessage(c.Request.Context(), messageID, uid, req.Content); err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "message updated successfully"})
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
    messageID, err := uuid.Parse(c.Param("messageId"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid message ID")
        return
    }
    
    userID := c.GetString("user_id")
    uid, _ := uuid.Parse(userID)
    
    if err := h.messageService.DeleteMessage(c.Request.Context(), messageID, uid); err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "message deleted successfully"})
}