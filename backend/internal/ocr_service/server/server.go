package server

import (
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/oOSomnus/transflate/internal/ocr_service/pool"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[OCR Service] ")
}

// extractPageNumber extracts the first numeric sequence found in a filename and returns it as an integer. Returns -1 if no number is found or conversion fails.
func extractPageNumber(filename string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(filename)
	if match == "" {
		return -1
	}
	num, err := strconv.Atoi(match)
	if err != nil {
		return -1
	}
	return num
}

/*
OCRServiceServer represents the server for the OCR (Optical Character Recognition) service.

Embedded Struct:
  - pb.UnimplementedOCRServiceServer: Embeds the unimplemented server methods to comply with gRPC server requirements.
*/
type OCRServiceServer struct {
	pb.UnimplementedOCRServiceServer
}

// ProcessPDF handles the processing of a PDF file into text using OCR.
// It converts the PDF into images and performs OCR on each page concurrently.
// Returns a StringListResponse containing the extracted text and the total page count.
func (s *OCRServiceServer) ProcessPDF(ctx context.Context, req *pb.PDFRequest) (*pb.StringListResponse, error) {
	// Create temp folder
	log.Println("Received PDF Process request")
	log.Println("Creating temp file ...")
	lang := req.Language
	tmpFile, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name()) // Clear up temp files

	if _, err := tmpFile.Write(req.PdfData); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %v", err)
	}
	log.Println("Creating temp folder ...")
	// Create output directory for images
	outputDir, err := os.MkdirTemp("", "pdf-pages-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir) // Clear up temp files

	// Use pdftoppm to convert PDF pages to PNG images
	log.Println("Converting pdf to png images ...")
	outputPattern := filepath.Join(outputDir, "page")
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %v", err)
	}
	cmd := exec.Command("pdftoppm", "-png", tmpFile.Name(), outputPattern)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run pdftoppm: %v", err)
	}
	log.Println("Images converted successfully.")
	// Read generated image files
	files, err := os.ReadDir(outputDir)
	pageNumber := len(files)
	if err != nil {
		return nil, fmt.Errorf("failed to read converted images: %v", err)
	}

	// Sort files by name (to ensure page order)
	sort.Slice(
		files, func(i, j int) bool {
			return extractPageNumber(files[i].Name()) < extractPageNumber(files[j].Name())
		},
	)

	// Worker pool for concurrent OCR
	ocrResults := make([]string, len(files))

	// Acquiring number of cpus
	numCPU := runtime.NumCPU()

	// Acquiring gosseract pool
	gossPool := pool.NewGosseractPool(numCPU+1, lang)
	defer gossPool.Close()
	var wg sync.WaitGroup
	workerPool := make(chan struct{}, numCPU+1)
	log.Println("Starting worker pool ...")
	for i, file := range files {
		wg.Add(1)
		workerPool <- struct{}{} // Acquire a worker slot

		go func(index int, fileName string) {
			defer wg.Done()
			defer func() { <-workerPool }() // Release the worker slot
			client := gossPool.Get()
			defer gossPool.Put(client)

			imagePath := filepath.Join(outputDir, fileName)
			err := client.SetImage(imagePath)
			if err != nil {
				log.Printf("failed to set image %v: %v", fileName, err)
				return
			}

			text, err := client.Text()
			if err != nil {
				log.Printf("OCR failed for %s: %v", imagePath, err)
				text = ""
			}

			ocrResults[index] = text
		}(i, file.Name())
	}
	log.Println("Waiting worker pool to finish.")
	wg.Wait() // Wait for all workers to complete
	log.Println("Worker pool finished.")
	return &pb.StringListResponse{Lines: ocrResults, PageNum: uint32(pageNumber)}, nil
}
