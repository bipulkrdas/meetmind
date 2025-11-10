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

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3FileStorage struct {
	uploader      *manager.Uploader
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	attachRepo    repository.AttachmentRepository
}

func NewS3FileStorage(cfg *config.Config, attachRepo repository.AttachmentRepository) (FileStorage, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.StorageEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.StorageEndpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.StorageRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.StorageAccessKey, cfg.StorageSecretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3FileStorage{
		uploader:      manager.NewUploader(s3Client),
		s3Client:      s3Client,
		presignClient: s3.NewPresignClient(s3Client),
		bucket:        cfg.StorageBucket,
		attachRepo:    attachRepo,
	}, nil
}

func (s *S3FileStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, roomID uuid.UUID) (*model.Attachment, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%s%s", roomID.String(), uuid.New().String(), ext)

	contentType := file.Header.Get("Content-Type")

	_, err = s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filename),
		Body:        src,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	presignResult, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(7*24*time.Hour))
	if err != nil {
		return nil, err
	}

	attachment := &model.Attachment{
		ID:          uuid.New(),
		FileName:    file.Filename,
		FileType:    contentType,
		FileSize:    file.Size,
		StoragePath: filename,
		StorageURL:  presignResult.URL,
	}

	if err := s.attachRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *S3FileStorage) DeleteFile(ctx context.Context, attachmentID uuid.UUID) error {
	attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return err
	}

	_, err = s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(attachment.StoragePath),
	})
	if err != nil {
		return err
	}

	return s.attachRepo.Delete(ctx, attachmentID)
}

func (s *S3FileStorage) GetFile(ctx context.Context, attachmentID uuid.UUID) (io.Reader, error) {
	attachment, err := s.attachRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(attachment.StoragePath),
	})
	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}
