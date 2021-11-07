package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/dto"
	"github.com/renbou/dontstress/internal/utils"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/labs", func(c *fiber.Ctx) error {
		labs, err := dao.LabDao().GetAll()
		if ok := utils.Check(c, err); !ok {
			return err
		}

		var labdtos []dto.LabDTO
		for _, lab := range labs {
			labdtos = append(labdtos, *lab.ToDTO())
		}
		return c.JSON(labdtos)
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
