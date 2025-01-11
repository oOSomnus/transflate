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

// UploadFileToS3 uploads a local file to an S3 bucket at the specified key and sets an expiration metadata value.
func UploadFileToS3(bucketName, objectKey, filePath string, expirationDays int) error {
	// load default config
	cfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				viper.GetString("aws.access.key.id"), viper.GetString("aws.secret.access.key"), "",
			),
		), config.WithRegion(viper.GetString("aws.region")),
	)
	if err != nil {
		return fmt.Errorf("unable to load config %w", err)
	}

	//format expiration
	expiration := time.Now().AddDate(0, 0, expirationDays).Format(time.RFC1123)
	// create s3 client
	client := s3.NewFromConfig(cfg)

	// open file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("can't open file %w", err)
	}
	defer file.Close()

	// upload
	_, err = client.PutObject(
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

// DownloadFileFromS3 downloads a file from an S3 bucket to a local file path.
func DownloadFileFromS3(bucketName, objectKey, destination string) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("can't load s3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// get obj
	resp, err := client.GetObject(
		context.Background(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to download file %w", err)
	}
	defer resp.Body.Close()

	// save to local
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create local file %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("file successfully downloaded: %s\n", destination)
	return nil
}

// GeneratePresignedURL generates a presigned URL for accessing an object in an S3 bucket with a specified expiration time.
// It requires the bucket name, object key, and expiration duration as inputs and returns the presigned URL or an error.
func GeneratePresignedURL(bucketName, objectKey string, expiration time.Duration) (string, error) {
	//utils.LoadEnv()
	cfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				viper.GetString("aws.access.key.id"), viper.GetString("aws.secret.access.key"), "",
			),
		), config.WithRegion(viper.GetString("aws.region")),
	)
	if err != nil {
		return "", fmt.Errorf("unable to load s3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	presignClient := s3.NewPresignClient(client)

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
