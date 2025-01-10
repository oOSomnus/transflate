package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task handler] ")
}

/*
TaskSubmit handles the task submission process, including PDF file validation, OCR processing, token balance deduction, and text translation.

Parameters:
  - c (*gin.Context): The HTTP context that contains the request and response writer.

Responses:
  - 200 OK: Returns the translated text as JSON.
  - 400 Bad Request: If the uploaded file is invalid or not a PDF.
  - 401 Unauthorized: If the user is not logged in.
  - 500 Internal Server Error: For errors during file reading, gRPC communication, token deduction, or other server-side issues.
*/
func TaskSubmit(c *gin.Context) {
	// Check login status
	log.Println("Received new task.")
	username, exists := c.Get("username")
	usernameStr, ok := username.(string)
	log.Println("Validating information ...")
	if !ok {
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"error": "Invalid username",
			},
		)
		return
	}
	if !exists {
		c.JSON(
			http.StatusUnauthorized, gin.H{
				"error": "User not authorized to submit task",
			},
		)
		return
	}
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Document invalid",
			},
		)
		return
	}
	//check whether it's pdf
	if filepath.Ext(file.Filename) != ".pdf" {
		log.Println("File extension not pdf")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}
	// Open the uploaded file
	fileContent, err := utils.OpenFile(file)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse file"})
	}
	lang := c.DefaultPostForm("lang", "eng")

	transResponse, err := usecase.ProcessOCRAndTranslate(usernameStr, fileContent, lang)
	if err != nil {
		log.Printf("Error processing OCR: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process OCR"})
	}
	//log.Printf("transresponse: %s", transResponse)

	downLink, err := CreateDownloadLinkWithMdString(transResponse)
	//covert md into html
	if err != nil {
		log.Fatalf("Failed to create download link: %v", err)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"data": downLink,
		},
	)
}

/*
CreateDownloadLinkWithMdString generates a downloadable link for a PDF file created from a Markdown string.

Parameters:
  - mdString (string): The input Markdown content to be converted into a PDF.

Returns:
  - (string): A presigned URL for downloading the generated PDF file.
  - (error): An error if the process of creating the file, converting the Markdown, or generating the link fails.
*/
func CreateDownloadLinkWithMdString(mdString string) (string, error) {
	//utils.LoadEnv()
	bucketName := viper.GetString("s3.bucket.name")
	mdTmpFile, err := os.CreateTemp("", "respMd-*.md")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Error creating temp file"))
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(errors.Wrapf(err, "Failed to remove file: %s", name))
		}
	}(mdTmpFile.Name())
	_, err = mdTmpFile.Write([]byte(mdString))
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Error writing to temp file"))
	}
	err = service.UploadFileToS3(bucketName, "mds/"+filepath.Base(mdTmpFile.Name()), mdTmpFile.Name(), 1)
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload file")
	}
	downLink, err := service.GeneratePresignedURL(bucketName, "mds/"+filepath.Base(mdTmpFile.Name()), time.Hour)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate presigned url")
	}
	return downLink, nil
}
