package models

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/liquiloans/sftp/config"
	"github.com/liquiloans/sftp/connection"
)

var S3Con *session.Session

func init() {
	connection.S3Connection()
	S3Con = connection.GetS3()
}

func UploadFileToS3(filename string, fileType string, bucket string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)
	currentDate := time.Now().Format("2006-01-02")
	///ll-sftp-storage/Prudent/2022-08-06/GetInvestorDashboard
	fileDir := fmt.Sprintf("/%s/%s/%s/%s", bucket, currentDate, fileType, filename)

	_, err = s3.New(S3Con).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.SFTP_BUCKET),
		Key:    aws.String(fileDir),
		//		ACL:                aws.String("private"),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
		//		ContentDisposition: aws.String("attachment"),
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("File ulpaded to S3 is ", filename)
}
