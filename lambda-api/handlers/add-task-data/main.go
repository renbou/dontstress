package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/dao/S3"
	"github.com/renbou/dontstress/internal/models"
	"github.com/renbou/dontstress/internal/utils"
	"github.com/renbou/dontstress/lambda-api/auth"
	"strconv"
)

type payload struct {
	Filetype string `json:"type" validate:"required"`
	File     struct {
		Lang string `json:"lang" validate:"required"`
		Data string `json:"data" validate:"required"`
	} `json:"file" validate:"required"`
}

var app = fiber.New()
var adapter = fiberadapter.New(app)

func initApp() {
	app.Use(auth.New())

	app.Post("/lab/:labid/task/:taskid", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))

		if ok := utils.Check(c, err); !ok {
			return err
		}

		var payload payload
		err = json.Unmarshal(c.Body(), &payload)

		if ok := utils.Validate(c, payload); !ok {
			return err
		}

		if ok := utils.Check(c, err); !ok {
			return err
		}

		id, err := S3.UploadFile(payload.File.Data)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		file := models.File{Id: id, Lang: payload.File.Lang}

		err = dao.FileDao().Create(&file)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		task := models.Task{LabId: labId, Num: taskId}
		if payload.Filetype == "generator" {
			task.Generator = id
		} else {
			task.Validator = id
		}

		err = dao.TaskDao().Update(&task)
		if ok := utils.Check(c, err); !ok {
			return err
		}

		return c.JSON(fiber.Map{
			"type": payload.Filetype,
			"file": file.ToDTO(payload.File.Data),
		})
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
