package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
)

const (
	REGION = "eu-west-1"
)

var (
	AdminsDynamoName   = os.Getenv("ADMINS_TABLE_NAME")
	LabsDynamoName     = os.Getenv("LABS_TABLE_NAME")
	TasksDynamoName    = os.Getenv("TASKS_TABLE_NAME")
	FilesDynamoName    = os.Getenv("FILES_TABLE_NAME")
	TestrunsDynamoName = os.Getenv("RUNS_TABLE_NAME")
)

var db = dynamo.New(session.Must(session.NewSession()), &aws.Config{Region: aws.String(REGION)})

func getDB() *dynamo.DB {
	return db
}
