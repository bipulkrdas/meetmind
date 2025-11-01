package service

import (
	"context"
	"errors"
	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/repository"
	"livekit-consulting/backend/internal/service/email"
	"time"

	"github.com/google/uuid"
)

type ParticipantService struct {
	participantRepo repository.ParticipantRepository
	roomRepo        repository.RoomRepository
	inviteRepo      repository.InviteRepository
	emailService    *email.EmailService
	livekitService  *LiveKitService
	frontendURL     string
}

func NewParticipantService(
	participantRepo repository.ParticipantRepository,
	roomRepo repository.RoomRepository,
	inviteRepo repository.InviteRepository,
	emailService *email.EmailService,
	livekitService *LiveKitService,
	frontendURL string,
) *ParticipantService {
	return &ParticipantService{
		participantRepo: participantRepo,
		roomRepo:        roomRepo,
		inviteRepo:      inviteRepo,
		emailService:    emailService,
		livekitService:  livekitService,
		frontendURL:     frontendURL,
	}
}

func (s *ParticipantService) AddParticipant(
	ctx context.Context,
	roomID uuid.UUID,
	inviterID uuid.UUID,
	req *model.AddParticipantRequest,
) (*model.ParticipantInviteResponse, error) {
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	existing, err := s.participantRepo.GetByRoomAndEmail(ctx, roomID, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("participant already added")
	}

	participantID := uuid.New()
	participant := &model.RoomParticipant{
		ParticipantID: &participantID,
		RoomID:        roomID,
		Email:         req.Email,
		Name:          req.Name,
		Role:          "participant",
	}

	err = s.participantRepo.Create(ctx, participant)
	if err != nil {
		return nil, err
	}

	inviteToken := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	err = s.inviteRepo.Create(ctx, roomID, inviterID, req.Email, req.Name, inviteToken, expiresAt)
	if err != nil {
		return nil, err
	}

	inviteURL := s.frontendURL + "/join/" + roomID.String() + "/prep?token=" + inviteToken
	err = s.emailService.SendRoomInviteEmail(ctx, req.Email, room.RoomName, inviteURL)
	if err != nil {
		// Log error but don't fail
	}

	return &model.ParticipantInviteResponse{
		ParticipantID: participantID,
		InviteToken:   inviteToken,
		InviteURL:     inviteURL,
	}, nil
}

func (s *ParticipantService) GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]*model.RoomParticipant, error) {
	return s.participantRepo.GetByRoomID(ctx, roomID)
}

func (s *ParticipantService) RemoveParticipant(ctx context.Context, roomID, participantID uuid.UUID) error {
	return s.participantRepo.Delete(ctx, participantID)
}

func (s *ParticipantService) GenerateParticipantToken(
	ctx context.Context,
	roomID uuid.UUID,
	inviteToken string,
) (string, error) {
	invite, err := s.inviteRepo.GetByToken(ctx, inviteToken)
	if err != nil {
		return "", errors.New("invalid invite")
	}

	if invite.RoomID != roomID || time.Now().After(invite.ExpiresAt) {
		return "", errors.New("invalid or expired invite")
	}

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return "", err
	}

	identity := invite.InviteeEmail
	token, err := s.livekitService.GenerateToken(
		identity,
		*room.LiveKitRoomName,
		true,
		true,
	)

	if err != nil {
		return "", err
	}

	s.inviteRepo.MarkAsAccepted(ctx, invite.ID)

	return token, nil
}

func (s *ParticipantService) GenerateInternalParticipantToken(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (string, error) {
	participant, err := s.participantRepo.GetByRoomAndUserID(ctx, roomID, userID)
	if err != nil {
		return "", err
	}
	if participant == nil {
		return "", errors.New("access denied: user is not a participant of this room")
	}

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return "", err
	}
	if room == nil {
		return "", errors.New("room not found")
	}
	if room.LiveKitRoomName == nil || *room.LiveKitRoomName == "" {
		return "", errors.New("livekit room not created for this room yet")
	}

	identity := participant.Email
	token, err := s.livekitService.GenerateToken(
		identity,
		*room.LiveKitRoomName,
		true, // canPublish
		true, // canSubscribe
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
