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

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type GCSFileStorage struct {
	client     *storage.Client
	bucket     string
	attachRepo repository.AttachmentRepository
	env        string // Store the environment
}

func NewGCSFileStorage(ctx context.Context, cfg *config.Config, attachRepo repository.AttachmentRepository) (FileStorage, error) {

	var opts []option.ClientOption
	if cfg.StorageEndpoint != "" {
		// For local testing with a GCS emulator
		opts = append(opts, option.WithEndpoint(cfg.StorageEndpoint))
		opts = append(opts, option.WithoutAuthentication())
	} else {
		// In production, use Application Default Credentials
		// Ensure GOOGLE_APPLICATION_CREDENTIALS is set or gcloud auth is configured
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GCSFileStorage{
		client:     client,
		bucket:     cfg.StorageBucket,
		attachRepo: attachRepo,
		env:        cfg.Env, // Store the environment
	}, nil
}

func (s *GCSFileStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, roomID uuid.UUID) (*model.Attachment, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%s%s", roomID.String(), uuid.New().String(), ext)

	wc := s.client.Bucket(s.bucket).Object(filename).NewWriter(ctx)
	wc.ContentType = file.Header.Get("Content-Type")

	if _, err = io.Copy(wc, src); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}

	var signedURL string
	if s.env == "development" {
		signedURL = "" // Empty string for development environment
	} else {
		// Attempt to generate signed URL, relying on ADC to infer the signer
		opts := &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(7 * 24 * time.Hour),
		}
		var err error
		signedURL, err = storage.SignedURL(s.bucket, filename, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate signed URL in production: %w", err)
		}
	}

	attachment := &model.Attachment{
		ID:          uuid.New(),
		FileName:    file.Filename,
		FileType:    file.Header.Get("Content-Type"),
		FileSize:    file.Size,
		StoragePath: filename,
		StorageURL:  signedURL,
	}

	if err := s.attachRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *GCSFileStorage) DeleteFile(ctx context.Context, attachmentID uuid.UUID) error {
	attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return err
	}

	err = s.client.Bucket(s.bucket).Object(attachment.StoragePath).Delete(ctx)
	if err != nil {
		return err
	}

	return s.attachRepo.Delete(ctx, attachmentID)
}

func (s *GCSFileStorage) GetFile(ctx context.Context, attachmentID uuid.UUID) (io.Reader, error) {
	attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	return s.client.Bucket(s.bucket).Object(attachment.StoragePath).NewReader(ctx)
}