package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	fiberadapter "github.com/renbou/aws-lambda-go-api-proxy/fiber"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
	"github.com/renbou/dontstress/internal/utils"
	"github.com/renbou/dontstress/lambda-api/auth"
)

var app = fiber.New()
var adapter = fiberadapter.New(app)

func initApp() {
	app.Use(auth.New())

	app.Post("/labs", func(c *fiber.Ctx) error {
		var lab models.Lab
		err := json.Unmarshal(c.Body(), &lab)

		if ok := utils.Validate(c, lab); !ok {
			return err
		}

		if ok := utils.Check(c, err); !ok {
			return err
		}

		lab.Id = utils.GetId()
		err = dao.LabDao().Create(&lab)
		if ok := utils.Check(c, err); !ok {
			return err
		}
		return c.JSON(lab.Id)
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
