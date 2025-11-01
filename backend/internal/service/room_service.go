package service

import (
	"context"
	"errors"
	"livekit-consulting/backend/internal/middleware"
	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/repository"

	"github.com/google/uuid"
	"github.com/livekit/protocol/livekit"
)

type RoomService struct {
    roomRepo        repository.RoomRepository
    participantRepo repository.ParticipantRepository
    livekitService  *LiveKitService
}

func NewRoomService(
    roomRepo repository.RoomRepository,
    participantRepo repository.ParticipantRepository,
    livekitService *LiveKitService,
) *RoomService {
    return &RoomService{
        roomRepo:        roomRepo,
        participantRepo: participantRepo,
        livekitService:  livekitService,
    }
}

func (s *RoomService) CreateRoom(ctx context.Context, userID uuid.UUID, req *model.CreateRoomRequest) (*model.Room, error) {
    livekitRoomName := "room_" + uuid.New().String()
    
    lkRoom, err := s.livekitService.CreateRoom(ctx, livekitRoomName)
    if err != nil {
        return nil, err
    }
    
    room := &model.Room{
        RoomName:        req.RoomName,
        Description:     req.Description,
        OwnerID:         userID,
        LiveKitRoomName: &livekitRoomName,
        RoomSID:         &lkRoom.Sid,
    }
    
    err = s.roomRepo.Create(ctx, room)
    if err != nil {
        s.livekitService.DeleteRoom(ctx, livekitRoomName)
        return nil, err
    }
    
    // Get user from context to get email and name
    user, ok := middleware.UserFrom(ctx)
    if !ok {
        return nil, errors.New("user not found in context")
    }

    err = s.participantRepo.Create(ctx, &model.RoomParticipant{
        RoomID: room.ID,
        UserID: &userID,
        Email:  user.Email,
        Name:   user.Name,
        Role:   "owner",
    })
    
    return room, nil
}

func (s *RoomService) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]*model.RoomResponse, error) {
    ownedRooms, err := s.roomRepo.GetByOwnerID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    participantRooms, err := s.roomRepo.GetRoomsByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    roomMap := make(map[uuid.UUID]*model.RoomResponse)
    
    for _, room := range ownedRooms {
        count, _ := s.participantRepo.CountByRoomID(ctx, room.ID)
        roomMap[room.ID] = &model.RoomResponse{
            Room:             *room,
            ParticipantCount: count,
            IsOwner:          true,
        }
    }
    
    for _, room := range participantRooms {
        if _, exists := roomMap[room.ID]; !exists {
            count, _ := s.participantRepo.CountByRoomID(ctx, room.ID)
            roomMap[room.ID] = &model.RoomResponse{
                Room:             *room,
                ParticipantCount: count,
                IsOwner:          false,
            }
        }
    }
    
    rooms := make([]*model.RoomResponse, 0, len(roomMap))
    for _, room := range roomMap {
        rooms = append(rooms, room)
    }
    
    return rooms, nil
}

func (s *RoomService) GetRoomDetails(ctx context.Context, roomID, userID uuid.UUID) (*model.RoomResponse, error) {
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    
    hasAccess, err := s.participantRepo.UserHasAccess(ctx, roomID, userID)
    if err != nil || !hasAccess {
        return nil, errors.New("access denied")
    }
    
    count, err := s.participantRepo.CountByRoomID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    
    return &model.RoomResponse{
        Room:             *room,
        ParticipantCount: count,
        IsOwner:          room.OwnerID == userID,
    }, nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, roomID, userID uuid.UUID) error {
    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return err
    }
    
    if room.OwnerID != userID {
        return errors.New("only room owner can delete")
    }
    
    if room.LiveKitRoomName != nil {
        s.livekitService.DeleteRoom(ctx, *room.LiveKitRoomName)
    }
    
    return s.roomRepo.Delete(ctx, roomID)
}

func (s *RoomService) CreateLiveKitRoom(ctx context.Context, roomID, userID uuid.UUID) (*livekit.Room, error) {
    participant, err := s.participantRepo.GetByRoomAndUserID(ctx, roomID, userID)
    if err != nil {
        return nil, err
    }

    if participant == nil || participant.Role != "owner" {
        return nil, errors.New("permission denied: only room owner can create livekit room")
    }

    room, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return nil, err
    }
    if room == nil {
        return nil, errors.New("room not found")
    }

    livekitRoomName := "room_" + uuid.New().String()
    if room.LiveKitRoomName != nil && *room.LiveKitRoomName != "" {
        livekitRoomName = *room.LiveKitRoomName
    }

    lkRoom, err := s.livekitService.CreateRoom(ctx, livekitRoomName)
    if err != nil {
        return nil, err
    }

    room.LiveKitRoomName = &livekitRoomName
    room.RoomSID = &lkRoom.Sid
    err = s.roomRepo.Update(ctx, room)
    if err != nil {
        // Try to delete the created livekit room for consistency
        s.livekitService.DeleteRoom(ctx, livekitRoomName)
        return nil, err
    }

    return lkRoom, nil
}
