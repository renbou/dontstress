package connect

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//var app = fiber.New()
//var adapter = fiberadapter.New(app)
//
//func initApp() {
//	app.Use(auth.New())
//
//	app.Post("/labs", func(c *fiber.Ctx) error {
//		var lab models.Lab
//		err := json.Unmarshal(c.Body(), &lab)
//
//		if ok := utils.Validate(c, lab); !ok {
//			return err
//		}
//
//		if ok := utils.Check(c, err); !ok {
//			return err
//		}
//
//		lab.Id = utils.GetId()
//		err = dao.LabDao().Create(&lab)
//		if ok := utils.Check(c, err); !ok {
//			return err
//		}
//		return c.JSON(lab.Id)
//	})
//}

func handler(req *events.APIGatewayWebsocketProxyRequest) error {
	fmt.Println(req.Body)
	return nil
	//if resp, err := adapter.ProxyV2(request); err != nil {
	//	return events.APIGatewayV2HTTPResponse{}, err
	//} else {
	//	return resp, nil
	//}
}

func main() {
	//initApp()
	lambda.Start(handler)
}
