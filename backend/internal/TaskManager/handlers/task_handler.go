package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/TaskManager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"github.com/pkoukk/tiktoken-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid username type",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authorized to submit task",
		})
		return
	}
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "document invalid",
		})
		return
	}
	//check whether it's pdf
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}
	// Open the uploaded file
	fileContent, err := openFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read the uploaded file",
		})
		return
	}
	ocrConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close grpc connection: %v", err)
		}
	}(ocrConn)

	lang := c.DefaultPostForm("lang", "eng")

	ocrClient := pb.NewOCRServiceClient(ocrConn)
	ocrCtx, ocrCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer ocrCancel()
	log.Println("Start processing file ...")
	ocrResponse, err := ocrClient.ProcessPDF(ocrCtx, &pb.PDFRequest{PdfData: fileContent, Language: lang})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	if ocrResponse == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	respLines := ocrResponse.Lines
	log.Println("File processed successfully.")
	//merge strings
	var builder strings.Builder
	for _, line := range respLines {
		builder.WriteString(line)
	}
	mergedString := builder.String()
	c.JSON(http.StatusOK, gin.H{
		"data": mergedString,
	})
	mergedString = utils.RemoveNonUnicodeCharacters(mergedString)
	mergedString = utils.ReplaceMultipleSpaces(mergedString)
	//tokenize
	encoder, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	tokens := encoder.Encode(mergedString, nil, nil)
	numTokens := len(tokens)
	err = usecase.DecreaseBalance(usernameStr, numTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	//translate
	transConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	defer func() {
		err := transConn.Close()
		if err != nil {
			log.Printf("Failed to close grpc connection: %v", err)
		}
	}()
	translateClient := pbt.NewTranslateServiceClient(transConn)
	transCtx, transCancel := context.WithTimeout(context.Background(), 3600*time.Second)
	defer transCancel()
	transResponse, err := translateClient.ProcessTranslation(transCtx, &pbt.TranslateRequest{Text: mergedString})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": transResponse.Lines,
	})
	return
}

/*
openFile reads the contents of a given multipart file header into a byte slice.

Parameters:
  - file (*multipart.FileHeader): The multipart file header to open and read.

Returns:
  - ([]byte): A byte slice containing the file's contents.
  - (error): An error if the file cannot be opened, read, or closed properly.
*/
func openFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Printf("Failed to close src: %v", err)
		}
	}(src)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
