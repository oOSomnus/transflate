package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
)

// OpenFile reads the contents of a multipart file and returns it as a byte slice. It closes the file after reading.
func OpenFile(file *multipart.FileHeader) ([]byte, error) {
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
