package main

import (
	"github.com/renbou/dontstress/internal/utils"
	"github.com/renbou/dontstress/lambda-api/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Use(auth.New())

	app.Delete("/lab/:labid", func(c *fiber.Ctx) error {
		lab := models.Lab{Id: c.Params("labid")}
		err := dao.LabDao().Delete(&lab)
		if ok := utils.Check(c, err); !ok {
			return err
		}
		c.Status(fiber.StatusOK)
		return nil
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
