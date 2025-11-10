package service

import (
	"context"
	"errors"
	"time"

	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/repository"

	"github.com/google/uuid"
)

// AgentWebhookPayload defines the structure for the incoming transcript webhook
type AgentWebhookPayload struct {
	Event           string `json:"event"`
	RoomName        string `json:"room_name"`
	SessionStart    string `json:"session_start"`
	SessionEnd      string `json:"session_end"`
	TranscriptPaths struct {
		JSON        string `json:"json"`
		Text        string `json:"text"`
		JSONHTTPS   string `json:"json_https_url"`
		TextHTTPS   string `json:"text_https_url"`
	} `json:"transcript_paths"`
	S3Keys struct {
		JSON string `json:"json"`
		Text string `json:"text"`
	} `json:"s3_keys"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	ItemCount int    `json:"item_count"`
	Timestamp string `json:"timestamp"`
}

type MessageService struct {
	messageRepo     repository.MessageRepository
	participantRepo repository.ParticipantRepository
	attachmentRepo  repository.AttachmentRepository
	roomRepo        repository.RoomRepository
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	participantRepo repository.ParticipantRepository,
	attachmentRepo repository.AttachmentRepository,
	roomRepo repository.RoomRepository,
) *MessageService {
	return &MessageService{
		messageRepo:     messageRepo,
		participantRepo: participantRepo,
		attachmentRepo:  attachmentRepo,
		roomRepo:        roomRepo,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, req *model.CreateMessageRequest, roomID, userID uuid.UUID) (*model.Message, error) {
	isMember, err := s.participantRepo.UserHasAccess(ctx, roomID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a member of this room")
	}

	message := &model.Message{
		RoomID:      roomID,
		UserID:      &userID,
		Content:     req.Content,
		MessageType: model.MessageTypeUserMessage,
	}

	_, err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	for _, attachmentID := range req.AttachmentIDs {
		err := s.attachmentRepo.SetMessageID(ctx, attachmentID, message.ID)
		if err != nil {
			// Log error but don't fail the whole process
		}
	}

	fullMessage, err := s.messageRepo.GetMessageWithAttachments(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	return fullMessage, nil
}

func (s *MessageService) CreateTranscriptMessage(ctx context.Context, payload *AgentWebhookPayload) (*model.Message, error) {
	room, err := s.roomRepo.GetByName(ctx, payload.RoomName)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	sessionStart, _ := time.Parse(time.RFC3339, payload.SessionStart)
	sessionEnd, _ := time.Parse(time.RFC3339, payload.SessionEnd)

	message := &model.Message{
		RoomID:      room.ID,
		UserID:      nil, // System message
		Content:     "Meeting transcript is available.",
		MessageType: model.MessageTypeMeetingTranscript,
		ExtraData: &model.ExtraData{
			Transcript: &model.TranscriptData{
				Bucket: payload.Bucket,
				Region: payload.Region,
				S3Keys: model.S3Keys{
					JSON: payload.S3Keys.JSON,
					Text: payload.S3Keys.Text,
				},
				HTTPSUrls: model.HTTPSUrls{
					JSON: payload.TranscriptPaths.JSONHTTPS,
					Text: payload.TranscriptPaths.TextHTTPS,
				},
				SessionStart: sessionStart,
				SessionEnd:   sessionEnd,
			},
		},
	}

	_, err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	fullMessage, err := s.messageRepo.GetMessageWithAttachments(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	return fullMessage, nil
}

func (s *MessageService) GetMessages(ctx context.Context, roomID, userID uuid.UUID, limit int, before *uuid.UUID) ([]*model.Message, error) {
	isMember, err := s.participantRepo.UserHasAccess(ctx, roomID, userID)
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

	for i, msg := range messages {
		fullMsg, err := s.messageRepo.GetMessageWithAttachments(ctx, msg.ID)
		if err == nil {
			messages[i] = fullMsg
		}
	}

	return messages, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, messageID, userID uuid.UUID, content string) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	if message.UserID == nil || *message.UserID != userID {
		return errors.New("unauthorized to edit this message")
	}

	return s.messageRepo.Update(ctx, messageID, content)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID, userID uuid.UUID) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	if message.UserID == nil || *message.UserID != userID {
		return errors.New("unauthorized to delete this message")
	}

	return s.messageRepo.Delete(ctx, messageID)
}

func (s *MessageService) SearchMessages(ctx context.Context, roomID, userID uuid.UUID, query string, limit int) ([]*model.Message, error) {
	isMember, err := s.participantRepo.UserHasAccess(ctx, roomID, userID)
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

	if message.Metadata == nil {
		message.Metadata = &model.MessageMetadata{
			Reactions: []model.Reaction{},
		}
	}

	found := false
	for i, reaction := range message.Metadata.Reactions {
		if reaction.Emoji == emoji {
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
		message.Metadata.Reactions = append(message.Metadata.Reactions, model.Reaction{
			Emoji:   emoji,
			UserIDs: []uuid.UUID{userID},
			Count:   1,
		})
	}

	return s.messageRepo.UpdateMetadata(ctx, messageID, message.Metadata)
}

func (s *MessageService) UpdateLastRead(ctx context.Context, roomID, userID uuid.UUID, lastReadSeqNo int) error {
    return s.participantRepo.UpdateLastRead(ctx, roomID, userID, lastReadSeqNo)
}
