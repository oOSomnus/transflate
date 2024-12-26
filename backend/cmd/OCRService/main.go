package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"sync"

	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/otiai10/gosseract/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"google.golang.org/grpc"
)

type OCRServiceServer struct {
	pb.UnimplementedOCRServiceServer
}

/*
ProcessPDF processes a PDF file, extracts images from its pages, and performs OCR on each page to extract text.

Parameters:
  - ctx (context.Context): The context for managing request-scoped values, deadlines, and cancellation signals.
  - req (*pb.PDFRequest): The request containing the PDF data to be processed, provided as a byte array in `PdfData`.

Returns:
  - (*pb.StringListResponse): A response containing a list of strings where each string represents the OCR result for a corresponding page of the PDF.
  - (error): An error if any issues occur during processing, such as file creation, image extraction, or OCR execution.
*/

func (s *OCRServiceServer) ProcessPDF(ctx context.Context, req *pb.PDFRequest) (*pb.StringListResponse, error) {
	// Create temp folder
	tmpFile, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clear up temp files

	if _, err := tmpFile.Write(req.PdfData); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Extract images
	outputDir, err := os.MkdirTemp("", "pdf-pages-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir) // Clear up temp files

	if err := api.ExtractImagesFile(tmpFile.Name(), outputDir, nil, nil); err != nil {
		return nil, fmt.Errorf("failed to extract images: %v", err)
	}

	// Read extracted image files
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read extracted images: %v", err)
	}

	// Sort files by name (to ensure page order)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Worker pool for concurrent OCR
	ocrResults := make([]string, len(files))
	var wg sync.WaitGroup
	workerPool := make(chan struct{}, 3) // Limit to 3 concurrent workers

	for i, file := range files {
		wg.Add(1)
		workerPool <- struct{}{} // Acquire a worker slot

		go func(index int, fileName string) {
			defer wg.Done()
			defer func() { <-workerPool }() // Release the worker slot

			client := gosseract.NewClient()
			defer client.Close()

			imagePath := filepath.Join(outputDir, fileName)
			client.SetImage(imagePath)

			text, err := client.Text()
			if err != nil {
				log.Printf("OCR failed for %s: %v", imagePath, err)
				text = fmt.Sprintf("Error: %v", err)
			}

			ocrResults[index] = text
		}(i, file.Name())
	}

	wg.Wait() // Wait for all workers to complete

	return &pb.StringListResponse{Lines: ocrResults}, nil
}

func main() {
	// Start gRPC service
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOCRServiceServer(grpcServer, &OCRServiceServer{})

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
