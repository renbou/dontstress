package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	_ "io/ioutil"
	_ "mime/multipart"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		// TODO: implement polling
		//labId := c.Params("labid")
		//taskId := c.Params("taskid")
		//file := models.File{Id: }
		//dynamodb.FileImpl{}.Create()
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
