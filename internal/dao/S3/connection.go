package S3

import (
	"bytes"
	"os"

	"github.com/renbou/dontstress/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	REGION = "eu-west-1"
)

var (
	BucketName = os.Getenv("BUCKET_NAME")
)

func getSession() *s3.S3 {
	return s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func UploadFile(content string) (string, error) {
	id := utils.GetId()
	_, err := getSession().PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader([]byte(content)),
		Bucket: aws.String(BucketName),
		Key:    aws.String(id),
	})
	return id, err
}
