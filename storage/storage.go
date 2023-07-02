package storage

import (
	"context"
	"errors"
	"io"
)

const TypeS3 = "s3"

// Storage defines interface for store data file
type Storage interface {
	UploadFile(ctx context.Context, objectKey string, reader io.Reader) (string, error)
	DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, objectKey string) error
	GetURL(ctx context.Context, objectKey string) (string, error)
	Exist(ctx context.Context, objectKey string) (bool, error)
}

// NewStorage creates a instance of FileStore with provided config
func NewStorage(ctx context.Context, storageType string, config interface{}) (Storage, error) {
	switch storageType {
	case TypeS3:
		cfg, ok := config.(S3Config)
		if !ok {
			return nil, errors.New("invalid config")
		}
		return newS3Storage(ctx, cfg)
	}

	return nil, errors.New("not found storage type")
}
