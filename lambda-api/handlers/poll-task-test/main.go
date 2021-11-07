package main

import (
	"strconv"

	"github.com/renbou/dontstress/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
)

var app = fiber.New()
var adapter = fiberadapter.New(app)

func initApp() {
	app.Get("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))
		id := c.Query("id")
		testRun, err := dao.TestrunDao().GetById(id)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		if testRun.LabId != labId || testRun.TaskId != taskId {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Run id does not match with lab id or task id",
			})
		}

		return c.JSON(testRun.ToDTO())
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
