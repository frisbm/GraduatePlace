package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
}

func NewS3(config aws.Config) *S3 {
	client := s3.NewFromConfig(config)
	return &S3{
		client: client,
	}
}

func (s *S3) CreateBucket(ctx context.Context, bucketName string) error {
	_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

func (s *S3) UploadFile(ctx context.Context, bucketName, filename string, file []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(file),
	})
	return err
}

func (s *S3) GetObject(ctx context.Context, bucketName, filename string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		return nil, err
	}
	return obj.Body, nil
}
