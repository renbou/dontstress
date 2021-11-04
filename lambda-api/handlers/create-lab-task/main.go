package main

import (
	"encoding/json"
	"github.com/renbou/dontstress/internal/utils"
	_ "io/ioutil"
	_ "mime/multipart"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
)

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()
	app.Post("/lab/:labid/tasks", func(c *fiber.Ctx) error {
		var task models.Task
		err := json.Unmarshal(c.Body(), &task)
		if err = utils.Check(c, err); err != nil {
			return err
		}
		task.LabId = c.Params("labid")
		count, err := dao.TaskDao().GetCount(task.LabId)
		if err = utils.Check(c, err); err != nil {
			return err
		}
		task.Num = count
		err = dao.TaskDao().Create(&task)
		if err = utils.Check(c, err); err != nil {
			return err
		}
		return c.JSON(task)
	})

	adapter := fiberadapter.New(app)

	println(request.RouteKey)
	if resp, err := adapter.ProxyV2(request); err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	} else {
		return resp, nil
	}
}

func main() {
	lambda.Start(handler)
}
