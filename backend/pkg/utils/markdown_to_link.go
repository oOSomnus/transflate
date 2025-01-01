package utils

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

/*
CreateDownloadLinkWithMdString generates a downloadable link for a PDF file created from a Markdown string.

Parameters:
  - mdString (string): The input Markdown content to be converted into a PDF.

Returns:
  - (string): A presigned URL for downloading the generated PDF file.
  - (error): An error if the process of creating the file, converting the Markdown, or generating the link fails.
*/
func CreateDownloadLinkWithMdString(mdString string) (string, error) {
	mdTmpFile, err := os.CreateTemp("", "respMd-*.md")
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Error creating temp file"))
	}
	//var htmlBuf bytes.Buffer
	//if err := goldmark.Convert([]byte(mdString), &htmlBuf); err != nil {
	//	log.Fatalln("Failed to parse response:", err)
	//}
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
	//pdfTmpFile, err := os.CreateTemp("", "respPdf-*.pdf")
	//if err != nil {
	//	log.Fatalln(errors.Wrap(err, "Error creating temp file"))
	//}
	//defer func(name string) {
	//	err := os.Remove(name)
	//	if err != nil {
	//		log.Fatalln(errors.Wrapf(err, "Failed to remove file: %s", name))
	//	}
	//}(pdfTmpFile.Name())
	//convert into pdf
	//cmd := exec.Command(
	//	"pandoc",
	//	mdTmpFIle.Name(),
	//	"-o", pdfTmpFile.Name(),
	//	"--pdf-engine=xelatex",
	//	"-V", "mainfont=Noto Sans CJK SC", // 替换为你本地安装的中文字体
	//)
	//err = cmd.Run()
	err = UploadFileToS3(bucketName, "mds/"+filepath.Base(mdTmpFile.Name()), mdTmpFile.Name(), 1)
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload file")
	}
	downLink, err := GeneratePresignedURL(bucketName, "mds/"+filepath.Base(mdTmpFile.Name()), time.Hour)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate presigned url")
	}
	return downLink, nil
}
