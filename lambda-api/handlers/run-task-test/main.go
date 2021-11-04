package main

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/dao/S3"
	"github.com/renbou/dontstress/internal/models"
	"github.com/renbou/dontstress/internal/utils"
)

type payload struct {
	Lang string `json:"lang"`
	Data string `json:"data"`
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var payload payload
		err = json.Unmarshal(c.Body(), &payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		id, err := S3.UploadFile(payload.Data)

		file := models.File{Id: id, Lang: payload.Lang}

		err = dao.FileDao().Create(&file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		testrun := models.Run{
			Id:     utils.GetId(),
			Labid:  labId,
			Taskid: taskId,
			Fileid: id,
			Status: "QUEUE",
		}

		err = dao.TestrunDao().Create(&testrun)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
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
