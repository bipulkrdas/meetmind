package service

import (
	"context"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type LiveKitService struct {
	apiKey    string
	apiSecret string
	url       string
}

func NewLiveKitService(apiKey, apiSecret, url string) *LiveKitService {
	return &LiveKitService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		url:       url,
	}
}

// GenerateToken creates a LiveKit access token for a participant
func (s *LiveKitService) GenerateToken(identity, roomName string, canPublish, canSubscribe bool) (string, error) {
	at := auth.NewAccessToken(s.apiKey, s.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	at.SetVideoGrant(grant).SetIdentity(identity).SetValidFor(24 * time.Hour)

	return at.ToJWT()
}

// CreateRoom creates a new room in LiveKit
func (s *LiveKitService) CreateRoom(ctx context.Context, roomName string) (*livekit.Room, error) {
	roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)

	room, err := roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    300, // 5 minutes
		MaxParticipants: 50,
	})

	return room, err
}

// ListParticipants returns all participants in a room
func (s *LiveKitService) ListParticipants(ctx context.Context, roomName string) ([]*livekit.ParticipantInfo, error) {
	roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)

	res, err := roomClient.ListParticipants(ctx, &livekit.ListParticipantsRequest{
		Room: roomName,
	})

	if err != nil {
		return nil, err
	}

	return res.Participants, nil
}

// DeleteRoom removes a room from LiveKit
func (s *LiveKitService) DeleteRoom(ctx context.Context, roomName string) error {
	roomClient := lksdk.NewRoomServiceClient(s.url, s.apiKey, s.apiSecret)

	_, err := roomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
		Room: roomName,
	})

	return err
}
