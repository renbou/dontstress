package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/bebra", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"бебра": "понюхана",
		})
	})

	app.Post("/bebra", func(c *fiber.Ctx) error {
		var bebra map[string]string
		json.Unmarshal(c.Body(), &bebra)

		return c.JSON(fiber.Map{
			"получена бебра": bebra,
		})
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
