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

// UploadFileToS3 uploads a file to an S3 bucket with a specified key and expiration metadata.
// Parameters:
// - bucketName: The name of the target S3 bucket.
// - objectKey: The key for the uploaded object in the S3 bucket.
// - filePath: The local file path of the file to be uploaded.
// - expirationDays: The number of days until the file expires, stored as metadata.
// Returns an error if the upload process fails.
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

// DownloadFileFromS3 downloads a file from an S3 bucket to the specified local destination.
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

// GeneratePresignedURL generates a presigned URL for accessing an S3 object with a specified expiration duration.
// Parameters:
// - bucketName (string): The name of the S3 bucket.
// - objectKey (string): The key of the object in the S3 bucket.
// - expiration (time.Duration): The duration for which the presigned URL will remain valid.
// Returns:
// - (string): The generated presigned URL.
// - (error): An error if the URL generation fails.
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
