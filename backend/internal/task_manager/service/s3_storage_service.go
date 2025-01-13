package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"os"
	"time"
)

type S3StorageService interface {
	UploadFileToS3(bucketName, objectKey, filePath string, expirationDays int) error
	GeneratePresignedURL(bucketName, objectKey string, expiration time.Duration) (string, error)
}

type S3StorageServiceImpl struct {
	client *s3.Client
}

func NewS3StorageService() (*S3StorageServiceImpl, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				viper.GetString("aws.access.key.id"), viper.GetString("aws.secret.access.key"), "",
			),
		), config.WithRegion(viper.GetString("aws.region")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load config %w", err)
	}
	client := s3.NewFromConfig(cfg)
	return &S3StorageServiceImpl{client: client}, nil
}

// UploadFileToS3 uploads a local file to an S3 bucket at the specified key and sets an expiration metadata value.
func (s *S3StorageServiceImpl) UploadFileToS3(bucketName, objectKey, filePath string, expirationDays int) error {
	//format expiration
	expiration := time.Now().AddDate(0, 0, expirationDays).Format(time.RFC1123)

	// open file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("can't open file %w", err)
	}
	defer file.Close()

	// upload
	_, err = s.client.PutObject(
		context.Background(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
			Body:   file,
			Metadata: map[string]string{
				"Expires": expiration,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to uplaod file %w", err)
	}

	fmt.Printf("file uploaded to S3 successfully: s3://%s/%s\n", bucketName, objectKey)
	return nil
}

// GeneratePresignedURL generates a presigned URL for accessing an object in an S3 bucket with a specified expiration time.
// It requires the bucket name, object key, and expiration duration as inputs and returns the presigned URL or an error.
func (s *S3StorageServiceImpl) GeneratePresignedURL(bucketName, objectKey string, expiration time.Duration) (
	string, error,
) {

	presignClient := s3.NewPresignClient(s.client)

	// generate pre-signed url
	presignedURL, err := presignClient.PresignGetObject(
		context.Background(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
		}, s3.WithPresignExpires(expiration),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return presignedURL.URL, nil
}
