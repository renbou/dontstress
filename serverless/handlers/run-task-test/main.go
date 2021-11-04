package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/serverless/handlers/dao/S3"
	"github.com/renbou/dontstress/serverless/handlers/dao/dynamodb"
	"github.com/renbou/dontstress/serverless/handlers/models"
	_ "io/ioutil"
	_ "mime/multipart"
)

type file struct {
	lang string
	data string
}

type payload struct {
	filetype string
	file     file
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		//labId := c.Params("labid")
		//taskId := c.Params("taskid")
		var payload payload
		err := json.Unmarshal([]byte(c.Body()), &payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		id, err := S3.UploadFile(payload.file.data)

		file := models.File{Id: id, Lang: payload.file.lang}
		err = dynamodb.FileImpl{}.Create(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(nil)
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
