package main

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/dao/S3"
	"github.com/renbou/dontstress/internal/models"
	"github.com/renbou/dontstress/internal/utils"
)

type payload struct {
	Lang string `json:"lang" validate:"required"`
	Data string `json:"data" validate:"required"`
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/lab/:labid/task/:taskid/test", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))

		if ok := utils.Check(c, err); !ok {
			return err
		}

		var payload payload
		err = json.Unmarshal(c.Body(), &payload)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		if ok := utils.Validate(c, payload); !ok {
			return err
		}

		id, err := S3.UploadFile(payload.Data)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		file := models.File{Id: id, Lang: payload.Lang}

		err = dao.FileDao().Create(&file)

		if ok := utils.Check(c, err); !ok {
			return err
		}

		testRun := models.Run{
			Id:     utils.GetId(),
			LabId:  labId,
			TaskId: taskId,
			FileId: id,
			Status: "QUEUE",
		}

		err = dao.TestrunDao().Create(&testRun)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		return c.JSON(testRun.Id)
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
