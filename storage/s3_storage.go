package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"hungdv136/gokit/logger"
)

const ErrCodeNotFound = "NotFound"

// S3Config defines config for s3 storage
type S3Config struct {
	AccessKeyID          string        `json:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey      string        `json:"secret_access_key" yaml:"secret_access_key"`
	Bucket               string        `json:"bucket" yaml:"bucket"`
	Region               string        `json:"region" yaml:"region"`
	Directory            string        `json:"directory" yaml:"directory"`
	PresignURLExpiration time.Duration `json:"presign_url_expiration" yaml:"presign_url_expiration"`
	MaxKeys              int64         `json:"max_keys" yaml:"max_keys"`

	// Set nil to use default value
	Endpoint         *string `json:"endpoint" yaml:"endpoint"`
	S3ForcePathStyle *bool   `json:"s3_force_path_style" yaml:"s3_force_path_style"`
	DisableSSL       *bool   `json:"disable_ssl" yaml:"disable_ssl"`
}

// S3Storage defines methods to access S3
type S3Storage struct {
	config  S3Config
	session *session.Session
	client  s3iface.S3API
}

// newS3Storage creates an instance of S3Storage
func newS3Storage(ctx context.Context, config S3Config) (*S3Storage, error) {
	awsConfig := &aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         config.Endpoint,
		DisableSSL:       config.DisableSSL,
		S3ForcePathStyle: config.S3ForcePathStyle,
	}

	if config.AccessKeyID != "" && config.SecretAccessKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, "")
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	return &S3Storage{config: config, session: sess, client: s3.New(sess)}, nil
}

// UploadFile reads from reader and uploads to S3
func (s *S3Storage) UploadFile(ctx context.Context, objectKey string, reader io.Reader) (string, error) {
	uploader := s3manager.NewUploader(s.session)
	path := s.getFilePath(objectKey)
	_, err := uploader.UploadWithContext(ctx,
		&s3manager.UploadInput{
			Bucket: aws.String(s.config.Bucket),
			Key:    aws.String(path),
			Body:   reader,
		})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("unable to upload %q to %q, %w", objectKey, s.config.Bucket, err))
		return "", err
	}

	return path, nil
}

// DownloadFile downloads file from S3 returns the body
func (s *S3Storage) DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(s.getFilePath(objectKey)),
	})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("unable to download %q to %q, %w", objectKey, s.config.Bucket, err))
		return nil, err
	}

	return result.Body, nil
}

// DeleteFile deletes file from S3
func (s *S3Storage) DeleteFile(ctx context.Context, objectKey string) error {
	_, err := s.client.DeleteObjectWithContext(ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.config.Bucket),
			Key:    aws.String(s.getFilePath(objectKey)),
		})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("unable to delete %q to %q, %w", objectKey, s.config.Bucket, err))
		return err
	}

	return nil
}

// GetURL returns the url of a file
func (s *S3Storage) GetURL(ctx context.Context, objectKey string) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(s.getFilePath(objectKey)),
	})
	req.SetContext(ctx)

	url, err := req.Presign(s.config.PresignURLExpiration)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("unable to presign %q to %q, %w", objectKey, s.config.Bucket, err))
		return "", err
	}

	return url, nil
}

// Exist checks if file is existed
func (s *S3Storage) Exist(ctx context.Context, objectKey string) (bool, error) {
	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(s.getFilePath(objectKey)),
	})
	if err != nil {
		//nolint:errorlint
		aErr, ok := err.(awserr.Error)
		if !ok {
			logger.Error(ctx, err)
			return false, err
		}

		code := aErr.Code()
		if code == s3.ErrCodeNoSuchKey || code == s3.ErrCodeNoSuchBucket || code == ErrCodeNotFound {
			return false, nil
		}

		logger.Error(ctx, err)
		return false, err
	}

	return true, nil
}

func (s *S3Storage) getFilePath(objectKey string) string {
	return filepath.Join(s.config.Directory, objectKey)
}
