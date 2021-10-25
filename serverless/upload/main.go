package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	_ "io/ioutil"
	"mime/multipart"
	_ "mime/multipart"
	"os"
)

const (
	REGION = "eu-west-1"
)

var (
	BucketName = os.Getenv("BUCKET_NAME")
	TasksDynamoName = os.Getenv("TASKS_TABLE_NAME")
	s3session  *s3.S3
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

type Upload struct {
	LabId    string `json:"id" form:"id" validate:"required"`
	TaskName string `json:"name" form:"name" validate:"required"`
	Language string `json:"lang" form:"lang" validate:"required"`
	Type     string `json:"type" form:"type" validate:"required,oneof=generator validator"`
}

type TaskDTO struct {
	id        string
	labId     string
	lang      string
	validator string
	generator string
}

//curl -X POST http://127.0.0.1:3000/upload -F "file=@test.c;filename=test.c" -F "lang=GCC" -F "id=abobalab" -F "name=pidortask" -F "type=validator"

func uploadFile(file *multipart.FileHeader) (string, error) {
	content, err := file.Open()
	if err != nil {
		return uploadFile(file)
	}
	id := uuid.New().String()
	_, err = s3session.PutObject(&s3.PutObjectInput{
		Body:   content,
		Bucket: aws.String(BucketName),
		Key:    aws.String(id),
	})
	return id, nil
}

func updateData(upload *Upload, fileId string) error {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("us-west-1")})
	table := db.Table(TasksDynamoName)
	task := TaskDTO{id: upload.TaskName, labId: upload.LabId, lang: upload.Language}
	if upload.Type == "validator" {
		task.validator = fileId
	} else {
		task.generator = fileId
	}
	return table.Update(upload.TaskName, task).Run()
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
		upload := new(Upload)

		if err := c.BodyParser(upload); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validate := validator.New()
		if err := validate.Struct(upload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		if file, err := c.FormFile("file"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File not found",
			})
		} else {
			id, err := uploadFile(file)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to upload to S3 Bucket",
				})
			}
			err = updateData(upload, id)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update dynamoDB",
				})
			}
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
