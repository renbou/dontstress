package main

import (
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Get("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))
		id := c.Query("id")
		testrun, err := dao.TestrunDao().GetById(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if testrun.Labid != labId || testrun.Taskid != taskId {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Run id does not match with lab id or task id",
			})
		}

		return c.JSON(testrun)
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
