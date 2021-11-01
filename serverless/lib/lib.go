package lib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
)

const (
	REGION = "eu-west-1"
)

type Upload struct {
	LabId string `json:"id" form:"id" validate:"required"`
}

type Lab struct {
	Id string `json:"id" dynamo:"id"`
}

var (
	LabsDynamoName = os.Getenv("LABS_TABLE_NAME")
)

func GetData(upload *Upload) (Lab, error) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(REGION)})
	table := db.Table(LabsDynamoName)
	var result Lab
	var kek []Lab
	table.Scan().All(&kek)
	err := table.Get("id", upload.LabId).
		One(&result)
	fmt.Println(LabsDynamoName)
	fmt.Println(kek)
	return result, err
}