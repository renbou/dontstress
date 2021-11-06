package main

import (
	_ "io/ioutil"
	_ "mime/multipart"

	"github.com/renbou/dontstress/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/lab/:labid/tasks", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		tasks, err := dao.TaskDao().GetAll(labId)
		if ok := utils.Check(c, err); !ok {
			return err
		}
		return c.JSON(tasks)
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
