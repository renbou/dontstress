package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/guregu/dynamo"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/serverless/lib"
	_ "io/ioutil"
	_ "mime/multipart"
	"os"
)

const (
	REGION = "eu-west-1"
)

var (
	TasksDynamoName = os.Getenv("TASKS_TABLE_NAME")
)

type Upload struct {
	LabId    string `json:"id" form:"id" validate:"required"`
	TaskName string `json:"name" form:"name" validate:"required"`
}

type TaskDTO struct {
	Id        string `dynamo:"id"`
	LabId     string `dynamo:"labId"`
	//Lang      string `dynamo:"lang"`
	//Validator string `dynamo:"validator"`
	//Generator string `dynamo:"generator"`
}

func putData(upload *Upload) error {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(REGION)})
	table := db.Table(TasksDynamoName)
	w := TaskDTO{Id: upload.TaskName, LabId: upload.LabId}
	return table.Put(w).Run()
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/create-task", func(c *fiber.Ctx) error {
		upload := new(Upload)

		if err := c.BodyParser(upload); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validate := validator.New()
		if err := validate.Struct(upload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request" + err.Error(),
			})
		}

		lab, err := lib.GetData(&lib.Upload{LabId: upload.LabId})
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "lab is not found. " + err.Error(),
			})
		}

		fmt.Println(lab)

		err = putData(upload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to put data in dynamoDB " + err.Error(),
			})
		}

		return c.JSON(upload)
	})

	adapter := fiberadapter.New(app)

	if resp, err := adapter.ProxyV2(request); err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	} else {
		return resp, nil
	}
}

func main() {
	lambda.Start(handler)
}
