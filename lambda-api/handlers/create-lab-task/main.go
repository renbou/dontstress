package main

import (
	"encoding/json"
	"github.com/renbou/dontstress/lambda-api/auth"

	"github.com/renbou/dontstress/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
)

var app = fiber.New()
var adapter = fiberadapter.New(app)

func initApp() {
	app.Use(auth.New())

	app.Post("/lab/:labid/tasks", func(c *fiber.Ctx) error {
		var task models.Task
		err := json.Unmarshal(c.Body(), &task)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		if ok := utils.Validate(c, task); !ok {
			return err
		}

		task.LabId = c.Params("labid")
		count, err := dao.TaskDao().GetCount(task.LabId)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		task.Num = count
		err = dao.TaskDao().Create(&task)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		return c.JSON(task.ToDTO())
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
