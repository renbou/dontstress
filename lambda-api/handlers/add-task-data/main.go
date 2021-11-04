package main

import (
	"encoding/json"
	_ "io/ioutil"
	_ "mime/multipart"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/dao/S3"
	"github.com/renbou/dontstress/internal/models"
)

type payload struct {
	Filetype string `json:"filetype"`
	File     struct {
		Lang string `json:"lang"`
		Data string `json:"data"`
	} `json:"file"`
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	app := fiber.New()

	app.Post("/lab/:labid/task/:taskid", func(c *fiber.Ctx) error {
		labId := c.Params("labid")
		taskId, err := strconv.Atoi(c.Params("taskid"))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var payload payload
		err = json.Unmarshal(c.Body(), &payload)

		// TODO: handle errors in a better way
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		id, err := S3.UploadFile(payload.File.Data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		file := models.File{Id: id, Lang: payload.File.Lang}

		err = dao.FileDao().Create(&file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		task := models.Task{LabId: labId, Num: taskId}
		if payload.Filetype == "generator" {
			task.Generator = id
		} else {
			task.Validator = id
		}

		err = dao.TaskDao().Update(&task)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(file)
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
