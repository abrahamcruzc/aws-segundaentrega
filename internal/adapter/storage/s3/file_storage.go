package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileStorage struct {
	client *s3.Client
	bucketName string
	endpoint string
}

func NewFileStorage(client *s3.Client, bucketName, endpoint string) *FileStorage {
	return &FileStorage{
		client: client,
		bucketName: bucketName,
		endpoint: endpoint,
	}
}

func (f *FileStorage) Upload(ctx context.Context, key string, file io.Reader, contentType string) (string, error) {
	_, err := f.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(f.bucketName),
		Key: aws.String(key),
		Body: file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error al subir archivo a s3: %w", err)
	}

	return f.GetURL(ctx, key), nil
}

func (f *FileStorage) GetURL(ctx context.Context, key string) string {
	return fmt.Sprintf("%s/%s/%s", f.endpoint, f.bucketName, key)
}

func (f *FileStorage) Delete(ctx context.Context, key string) error {
	_, err := f.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(f.bucketName),
		Key: aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error al eliminar archivo de S3: %w", err)
	}
	return nil
}

func (f *FileStorage) CreateBucket(ctx context.Context) error {
	_, err := f.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(f.bucketName),
	})
	if err != nil {
		return nil
	}
	return nil
}
