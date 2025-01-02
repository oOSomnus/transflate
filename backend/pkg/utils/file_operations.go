package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
)

/*
OpenFile reads the contents of a given multipart file header into a byte slice.

Parameters:
  - file (*multipart.FileHeader): The multipart file header to open and read.

Returns:
  - ([]byte): A byte slice containing the file's contents.
  - (error): An error if the file cannot be opened, read, or closed properly.
*/
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
