package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/renbou/dontstress/internal/dao"
	"net/http"
)

func handler(_ context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := request.RequestContext.ConnectionID

	fileId, ok := request.QueryStringParameters["fileId"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "No id is passed in query parameters",
		}, nil
	}

	err := dao.FileDao().Update(fileId, connectionId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(handler)
}
