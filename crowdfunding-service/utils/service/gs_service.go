package service

import (
	"fmt"
	"os"
	"bytes"
	"net/http"
	"mime/multipart"
	"io/ioutil"
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"log"
	"io"
	"../../setting"
	"google.golang.org/api/option"
)

var gsBucket *storage.BucketHandle

type GSService struct {
}

func (gsService GSService) Init() (error) {
	if gsBucket == nil {
		opt := option.WithCredentialsFile(setting.CurrentConfig().GSCredentialsFile)
		ctx := context.Background()
		client, err := storage.NewClient(ctx, opt)
		if err != nil {
			log.Print("Failed to create client: %v", err)
			return err
		}
		bucketName := setting.CurrentConfig().GSBucketName
		gsBucket = client.Bucket(bucketName)
	}
	return nil
}

func (gsService GSService) UploadOsFile(sourceFile *os.File, directory string, fileName string) error {
	err := gsService.Init()
	if err != nil {
		log.Print(err)
		return err
	}
	fileInfo, _ := sourceFile.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer
	sourceFile.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	directory = directory + "/" + fileName
	ctx := context.Background()
	w := gsBucket.Object(directory).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.CacheControl = "public, max-age=86400"
	w.ContentType = fileType
	if _, err := io.Copy(w, fileBytes); err != nil {
		log.Print(err)
		return err
	}
	if err := w.Close(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (gsService GSService) UploadFormFile(sourceFile multipart.File, directory string, fileName string, sourceFileHeader *multipart.FileHeader) error {
	err := gsService.Init()
	if err != nil {
		log.Print(err)
		return err
	}
	buffer, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		fmt.Println("Read file error: ", err)
		return err
	}
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	directory = directory + "/" + fileName

	ctx := context.Background()
	w := gsBucket.Object(directory).NewWriter(ctx)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.CacheControl = "public, max-age=86400"
	w.ContentType = fileType
	if _, err := io.Copy(w, fileBytes); err != nil {
		log.Print(err)
		return err
	}
	if err := w.Close(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
