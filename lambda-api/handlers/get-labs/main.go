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

var app = fiber.New()
var adapter = fiberadapter.New(app)

func initApp() {
	app.Get("/labs", func(c *fiber.Ctx) error {
		labs, err := dao.LabDao().GetAll()
		if ok := utils.Check(c, err); !ok {
			return err
		}

		labDtos := []dto.LabDTO{}
		for _, lab := range labs {
			labDtos = append(labDtos, *lab.ToDTO())
		}
		return c.JSON(labDtos)
	})
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if resp, err := adapter.ProxyV2(request); err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	} else {
		return resp, nil
	}
}

func main() {
	initApp()
	lambda.Start(handler)
}
