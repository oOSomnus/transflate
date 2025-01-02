package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"net/http"
	"path/filepath"
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
				"error": "Invalid username type",
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
				"error": "document invalid",
			},
		)
		return
	}
	//check whether it's pdf
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}
	// Open the uploaded file
	fileContent, err := utils.OpenFile(file)
	lang := c.DefaultPostForm("lang", "eng")

	transResponse, err := usecase.ProcessOCRAndTranslate(usernameStr, fileContent, lang)

	downLink, err := utils.CreateDownloadLinkWithMdString(transResponse)
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
