package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/serverless/handlers/dao/dynamodb"
	"github.com/renbou/dontstress/serverless/handlers/models"
	_ "io/ioutil"
	_ "mime/multipart"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()
	app.Post("/lab/:labid/tasks", func(c *fiber.Ctx) error {
		var task models.Task
		err := json.Unmarshal(c.Body(), &task)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		task.LabId = c.Params("labid")
		task.Num = 1
		err = dynamodb.TaskImpl{}.Create(task)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(task)
	})

	adapter := fiberadapter.New(app)

	println(request.RouteKey)
	if resp, err := adapter.ProxyV2(request); err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	} else {
		return resp, nil
	}
}

func main() {
	lambda.Start(handler)
}
