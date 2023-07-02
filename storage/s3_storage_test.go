package storage

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/hungdv136/gokit/util"
	"github.com/stretchr/testify/require"
)

func TestS3Storage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockS3Storage := newMockS3()
	objectKey := time.Now().String()
	data := util.RandomString(64, util.AlphaNumericCharacters)
	content := []byte(data)

	// upload file
	key, err := mockS3Storage.UploadFile(ctx, objectKey, bytes.NewReader(content))
	require.NoError(t, err)
	require.NotEmpty(t, key)

	// download file
	download, err := mockS3Storage.DownloadFile(ctx, objectKey)
	require.NoError(t, err)
	downloadContent := make([]byte, len(content))
	n, _ := download.Read(downloadContent)
	require.Len(t, content, n)
	require.Equal(t, content, downloadContent)

	existed, err := mockS3Storage.Exist(ctx, objectKey)
	require.NoError(t, err)
	require.True(t, existed)

	existed, err = mockS3Storage.Exist(ctx, uuid.NewString())
	require.NoError(t, err)
	require.False(t, existed)

	presignedURL, err := mockS3Storage.GetURL(ctx, objectKey)
	require.NoError(t, err)
	require.NotEmpty(t, presignedURL)

	// delete file
	err = mockS3Storage.DeleteFile(ctx, objectKey)
	require.NoError(t, err)
	err = download.Close()
	require.NoError(t, err)
}

func newMockS3() Storage {
	s, err := NewStorage(context.Background(), TypeS3, S3Config{
		Bucket:               "test",
		Directory:            "raw_data",
		Region:               "ap-southeast-1",
		AccessKeyID:          "test",
		SecretAccessKey:      "test",
		S3ForcePathStyle:     aws.Bool(true),
		DisableSSL:           aws.Bool(true),
		Endpoint:             aws.String("localhost:4566"),
		PresignURLExpiration: time.Minute,
	})
	if err != nil {
		panic(err)
	}

	return s
}
