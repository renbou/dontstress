package main

import (
	"github.com/renbou/dontstress/lambda-api/auth"
	"strconv"

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

	app.Delete("/lab/:labid/task/:taskid", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))
		if ok := utils.Check(c, err); !ok {
			return err
		}

		task := models.Task{Num: taskId, LabId: labId}
		err = dao.TaskDao().Delete(&task)
		if ok := utils.Check(c, err); !ok {
			return err
		}
		c.Status(fiber.StatusOK)
		return nil
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
