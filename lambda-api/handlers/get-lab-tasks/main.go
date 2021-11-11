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
	app.Get("/lab/:labid/tasks", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		tasks, err := dao.TaskDao().GetAll(labId)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		taskDtos := []dto.TaskDTO{}
		for _, task := range tasks {
			taskDtos = append(taskDtos, *task.ToDTO())
		}
		return c.JSON(taskDtos)
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
