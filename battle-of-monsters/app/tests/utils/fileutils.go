package utilstests

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func GenerateMultipartFormFile(fileName string) (*multipart.Writer, bytes.Buffer) {
	var csvFile *os.File
	var csvContent bytes.Buffer
	var fileWriter io.Writer

	fileDir, _ := os.Getwd()
	filePath := path.Join(fileDir, fileName)
	csvFile, _ = os.Open(filePath)

	writer := multipart.NewWriter(&csvContent)
	fileWriter, _ = writer.CreateFormFile("file", filepath.Base(csvFile.Name()))
	io.Copy(fileWriter, csvFile)
	csvFile.Close()
	writer.Close()

	return writer, csvContent
}
