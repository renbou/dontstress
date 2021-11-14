package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	wsclient "github.com/renbou/dontstress/lambda-api/websockets/client"
	"net/http"
)

type Payload struct {
	Message      string `json:"message"`
	Id           string `json:"id"`
	ConnectionId string `json:"connectionId"`
	Lang         string `json:"lang"`
}

func handler(_ context.Context, req *events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	record := req.Records[0].Change.NewImage

	connectionIdAttr, ok := record["connectionId"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "No connectionId was passed",
		}, nil
	}
	connectionId := connectionIdAttr.String()

	client := wsclient.New()

	client.Post(connectionId, Payload{
		Message:      "Stream was invoked",
		Id:           record["id"].String(),
		ConnectionId: connectionId,
		Lang:         record["lang"].String(),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(handler)
}
