package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kenrms/message-processing-service/messageData"
	"kenrms/message-processing-service/openai"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, apiGatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var messageData messageData.MessageData
	err := json.Unmarshal([]byte(apiGatewayEvent.Body), &messageData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Error parsing request body: %v", err),
		}, nil
	}

	// Process the message data as needed
	fmt.Printf("Received message data: %+v\n", messageData)

	// TODO send message to OpenAI for a reply
	reply, err := openai.GetReplyFromOpenAI(messageData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error processing message data: %v", err),
		}, nil
	}

	// Return a successful response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       reply,
	}, nil
}
