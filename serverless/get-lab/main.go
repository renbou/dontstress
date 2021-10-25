package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/serverless/lib"
	_ "io/ioutil"
	_ "mime/multipart"
)


func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/get-lab", func(c *fiber.Ctx) error {
		upload := new(lib.Upload)

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

		lab, err := lib.GetData(upload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to put data in dynamoDB",
			})
		}
		return c.JSON(lab)
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
