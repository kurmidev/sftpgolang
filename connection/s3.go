package connection

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var s3Con *session.Session

var AWS_S3_REGION string = GoDotEnvVariable("AWS_S3_REGION")
var AWS_SECRET string = GoDotEnvVariable("AWS_SECRET")
var AWS_TOKEN string = GoDotEnvVariable("AWS_TOKEN")

func S3Connection() {

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_S3_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_SECRET, AWS_TOKEN, ""),
	})
	if err != nil {
		log.Fatal(err)
	}
	s3Con = session
}

func GetS3() *session.Session {
	return s3Con
}
