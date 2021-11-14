package wsclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ApiGatewayManagementApi = apigatewaymanagementapi.ApiGatewayManagementApi
type WsClient struct {
	api *ApiGatewayManagementApi
}

func New() *WsClient {
	sess := session.Must(session.NewSession())
	client := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("WS_ENDPOINT")))

	return &WsClient{
		api: client,
	}
}

func (client WsClient) Post(connectionId string, object interface{}) {
	data, err := json.Marshal(object)

	output, err := client.api.PostToConnection(
		&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionId,
			Data:         data,
		},
	)

	if err != nil {
		fmt.Println(output)
		fmt.Println(err.Error())
	}
}
