package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"fmt"
	"os"
	"bytes"
	"net/http"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"../../setting"
	"time"
	"mime/multipart"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"io/ioutil"
)

type S3Service struct {
}

func (s3Service S3Service) Init() (*aws.Config, error) {
	aws_access_key_id := setting.CurrentConfig().AwsKeyId
	aws_secret_access_key := setting.CurrentConfig().AwsSecretKey
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
		return nil, err
	}
	region := aws.String(endpoints.UsEast1RegionID)
	cfg := aws.NewConfig().WithCredentials(creds).WithRegion(*region)
	return cfg, nil
}

func (s3Service S3Service) UploadOsFile(sourceFile *os.File, directory string, fileName string) error {
	cfn, err := s3Service.Init();
	if err != nil {
		return err
	}
	svc := s3.New(session.New(), cfn)

	fileInfo, _ := sourceFile.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer

	sourceFile.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	directory = directory + "/" + fileName
	cacheControl := "max-age=604800"
	expired := time.Date(2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
	acl := "public-read"
	params := &s3.PutObjectInput{
		Bucket:        aws.String(setting.CurrentConfig().S3BucketName),
		Key:           aws.String(directory),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		CacheControl:  aws.String(cacheControl),
		Expires:       aws.Time(expired),
		ACL:           aws.String(acl),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
		return err
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))
	return nil
}

func (s3Service S3Service) UploadFormFile(sourceFile multipart.File, directory string, fileName string, sourceFileHeader *multipart.FileHeader) error {
	cfn, err := s3Service.Init();
	if err != nil {
		return err
	}
	svc := s3.New(session.New(), cfn)

	buffer, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		fmt.Println("Read file error: ", err)
		return err
	}

	fileBytes := bytes.NewReader(buffer)
	directory = directory + "/" + fileName
	cacheControl := "max-age=604800"
	expired := time.Date(2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
	acl := "public-read"
	params := &s3.PutObjectInput{
		Bucket:       aws.String(setting.CurrentConfig().S3BucketName),
		Key:          aws.String(directory),
		Body:         fileBytes,
		ContentType:  aws.String(sourceFileHeader.Header.Get("Content-Type")),
		CacheControl: aws.String(cacheControl),
		Expires:      aws.Time(expired),
		ACL:          aws.String(acl),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Printf("bad response: %s", err)
		return err
	}
	fmt.Printf("response %s", awsutil.StringValue(resp))
	return nil
}
