package main

import (
	_ "io/ioutil"
	_ "mime/multipart"
	"strconv"

	"github.com/renbou/dontstress/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

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
		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
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
