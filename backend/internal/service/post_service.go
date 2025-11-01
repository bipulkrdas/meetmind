package service

import (
    "context"
    "livekit-consulting/backend/internal/model"
    "livekit-consulting/backend/internal/repository"

    "github.com/google/uuid"
)

type PostService struct {
    postRepo repository.PostRepository
    roomRepo repository.RoomRepository
}

func NewPostService(postRepo repository.PostRepository, roomRepo repository.RoomRepository) *PostService {
    return &PostService{postRepo: postRepo, roomRepo: roomRepo}
}

func (s *PostService) CreatePost(ctx context.Context, userID, roomID uuid.UUID, req *model.CreatePostRequest) (*model.Post, error) {
    post := &model.Post{
        RoomID:    roomID,
        CreatorID: userID,
        Message:   req.Message,
    }
    err := s.postRepo.Create(ctx, post)
    return post, err
}

func (s *PostService) GetPosts(ctx context.Context, roomID uuid.UUID) ([]*model.PostWithCreator, error) {
    return s.postRepo.GetByRoomID(ctx, roomID)
}

func (s *PostService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
    // In a real application, you would check if the user has permission to delete the post.
    return s.postRepo.Delete(ctx, postID)
}
