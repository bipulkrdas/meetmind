package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"livekit-consulting/backend/internal/config"
	"livekit-consulting/backend/internal/model"
	"livekit-consulting/backend/internal/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, roomID uuid.UUID) (*model.Attachment, error)
	DeleteFile(ctx context.Context, attachmentID uuid.UUID) error
	GetFile(ctx context.Context, attachmentID uuid.UUID) (io.Reader, error)
}

type MinioFileStorage struct {
	minioClient *minio.Client
	bucket      string
	attachRepo  repository.AttachmentRepository
}

func NewMinioFileStorage(cfg *config.Config, attachRepo repository.AttachmentRepository) (FileStorage, error) {
	minioClient, err := minio.New(cfg.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
		Secure: false, // Set to true for production with HTTPS
	})
	if err != nil {
		return nil, err
	}

	return &MinioFileStorage{
		minioClient: minioClient,
		bucket:      cfg.StorageBucket,
		attachRepo:  attachRepo,
	}, nil
}

func (s *MinioFileStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, roomID uuid.UUID) (*model.Attachment, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%s%s", roomID.String(), uuid.New().String(), ext)

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

	url, err := s.minioClient.PresignedGetObject(ctx, s.bucket, filename, 7*24*time.Hour, nil)
	if err != nil {
		return nil, err
	}

	attachment := &model.Attachment{
		ID:          uuid.New(),
		FileName:    file.Filename,
		FileType:    file.Header.Get("Content-Type"),
		FileSize:    file.Size,
		StoragePath: filename,
		StorageURL:  url.String(),
	}

	if err := s.attachRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *MinioFileStorage) DeleteFile(ctx context.Context, attachmentID uuid.UUID) error {
	attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return err
	}

	err = s.minioClient.RemoveObject(ctx, s.bucket, attachment.StoragePath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return s.attachRepo.Delete(ctx, attachmentID)
}

func (s *MinioFileStorage) GetFile(ctx context.Context, attachmentID uuid.UUID) (io.Reader, error) {
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
