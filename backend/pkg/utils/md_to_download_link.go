package utils

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"github.com/yuin/goldmark"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CreateDownloadLinkWithMdString(mdString string) (string, error) {
	tmpFile, err := os.CreateTemp("", "respPdf-*")
	var htmlBuf bytes.Buffer
	if err := goldmark.Convert([]byte(mdString), &htmlBuf); err != nil {
		log.Fatalln("Failed to parse response:", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalln(errors.Wrapf(err, "Failed to remove file: %s", name))
		}
	}(tmpFile.Name())
	htmlContent := htmlBuf.String()
	//convert into pdf
	responsePdf := gopdf.GoPdf{}
	responsePdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	responsePdf.AddPage()
	err = responsePdf.SetFont("Arial", "B", 12)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Failed to set the font"))
	}
	if err := responsePdf.Cell(nil, htmlContent); err != nil {
		log.Fatalln("Failed to parse response:", err)
	}
	if err := responsePdf.WritePdf("output.pdf"); err != nil {
		log.Fatalln("Failed to write pdf:", err)
	}
	err = UploadFileToS3(bucketName, filepath.Base(tmpFile.Name()), tmpFile.Name(), 1)
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload file")
	}
	downLink, err := GeneratePresignedURL(bucketName, filepath.Base(tmpFile.Name()), time.Hour)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate presigned url")
	}
	return downLink, nil
}
