package service

import (
	"context"
	"io"

	"livekit-consulting/backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3TranscriptStorage interface {
	GetTranscriptFile(ctx context.Context, key string) (io.ReadCloser, error)
}

type s3TranscriptStorage struct {
	s3Client *s3.Client
	bucket   string
}

func NewS3TranscriptStorage(cfg *config.Config) (S3TranscriptStorage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.TranscriptAWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.TranscriptAWSAccessKeyID, cfg.TranscriptAWSSecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Use path-style addressing for S3 compatible storage
	})

	return &s3TranscriptStorage{
		s3Client: s3Client,
		bucket:   cfg.TranscriptAWSBucket,
	}, nil
}

func (s *s3TranscriptStorage) GetTranscriptFile(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}
