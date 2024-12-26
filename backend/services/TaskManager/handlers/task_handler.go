package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/oOSomnus/transflate/services/TaskManager/usecase"
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

func TaskSubmit(c *gin.Context) {
	// Check login status
	username, exists := c.Get("username")
	usernameStr, ok := username.(string)
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
	}
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "document invalid",
		})
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
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close grpc connection: %v", err)
		}
	}(conn)

	client := pb.NewOCRServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.ProcessPDF(ctx, &pb.PDFRequest{PdfData: fileContent})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
	respLines := response.Lines
	//merge strings
	var builder strings.Builder
	for _, line := range respLines {
		builder.WriteString(line)
	}
	mergedString := builder.String()
	//tokenize
	encoder, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
	tokens := encoder.Encode(mergedString, nil, nil)
	numTokens := len(tokens)
	err = usecase.DecreaseBalance(usernameStr, numTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
	//TODO: Translate
}

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
