package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/oOSomnus/transflate/pkg/utils"
	"os"
	"time"
)

const (
	BucketName = "transflate-bucket"
)

/*
UploadFileToS3 uploads a local file to a specified S3 bucket.

Parameters:
  - bucketName (string): The name of the target S3 bucket.
  - objectKey (string): The key (path) to use for the file in the S3 bucket.
  - filePath (string): The local path of the file to upload.

Returns:
  - (error): An error if the configuration fails, the file cannot be opened, or the upload operation fails.
*/
func UploadFileToS3(bucketName, objectKey, filePath string, expirationDays int) error {
	// load default config
	utils.LoadEnv()
	cfg, err := config.LoadDefaultConfig(context.Background())
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

/*
DownloadFileFromS3 downloads a file from an S3 bucket and saves it to a specified local destination.

Parameters:
  - bucketName (string): The name of the S3 bucket containing the file to download.
  - objectKey (string): The key of the object in the S3 bucket to download.
  - destination (string): The local file path where the downloaded file should be saved.

Returns:
  - (error): An error if the S3 configuration fails, the file cannot be downloaded, the local file cannot be created, or there is an issue saving the file.
*/

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

/*
generatePresignedURL generates a pre-signed URL for accessing an object in an S3 bucket.

Parameters:
  - bucketName (string): The name of the S3 bucket containing the object.
  - objectKey (string): The key of the object in the S3 bucket for which the pre-signed URL is generated.
  - expiration (time.Duration): The duration for which the pre-signed URL is valid.

Returns:
  - (string): The generated pre-signed URL for accessing the specified object.
  - (error): An error if the configuration fails, the S3 client cannot be initialized, or the pre-signed URL generation fails.
*/

func GeneratePresignedURL(bucketName, objectKey string, expiration time.Duration) (string, error) {
	utils.LoadEnv()
	cfg, err := config.LoadDefaultConfig(context.Background())
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
