package storage

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	storagePorts "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/storage"
)

type StorageClient struct {
	client *storage.Client
}

func NewClient(ctx context.Context) storagePorts.StoragePort {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	return &StorageClient{
		client: client,
	}
}

func (c *StorageClient) Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader) error {
	bucket := c.client.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	_, err := io.Copy(writer, reader)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to upload object: %v", err))
	}
	err = writer.Close()
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to close writer: %v", err))
	}
	return nil
}

func (c *StorageClient) Download(ctx context.Context, path value.StorageInfo) (io.ReadCloser, error) {
	bucketName := path.BucketName()
	objectName := path.ObjectName()
	bucket := c.client.Bucket(bucketName)
	object := bucket.Object(objectName)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to download object: %v", err))
	}
	return reader, nil
}

func (c *StorageClient) Delete(ctx context.Context, path value.StorageInfo) error {
	bucketName := path.BucketName()
	objectName := path.ObjectName()
	bucket := c.client.Bucket(bucketName)
	object := bucket.Object(objectName)
	err := object.Delete(ctx)
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete object: %v", err))
	}
	return nil
}
