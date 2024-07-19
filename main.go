package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kenrms/jubibot-response-lambda/constants"
	"kenrms/jubibot-response-lambda/messageData"
	"kenrms/jubibot-response-lambda/openai"
	"kenrms/jubibot-response-lambda/redisBroker"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, apiGatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response string
	var msgData messageData.MessageData
	err := json.Unmarshal([]byte(apiGatewayEvent.Body), &msgData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("Error parsing request body: %v", err),
		}, nil
	}

	// Process the message data as needed
	fmt.Printf("Received message data: %+v\n", msgData)

	commandResponse, err := HandleCommands(msgData)
	// if there is a command response, return it
	if commandResponse != "" {
		response = commandResponse
	} else {

		// Get conversation history from redisBroker
		channelConversation, err := redisBroker.GetConversationHistory(msgData.ChannelID)
		// append new incoming message to conversation history
		channelConversation = append(channelConversation, msgData)

		reply, err := openai.GetReplyFromOpenAI(channelConversation)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Error processing message data: %v", err),
			}, nil
		}

		replyMessageData := messageData.MessageData{
			AuthorUsername: constants.BOT_NAME,
			MessageContent: reply,
			ChannelID:      msgData.ChannelID,
		}

		// if the message is not from the bot, set IsParentMessageBot to true

		// add the reply to the conversation history
		channelConversation = append(channelConversation, replyMessageData)

		// save the updated conversation history to redisBroker
		err = redisBroker.SaveConversationHistory(msgData.ChannelID, channelConversation)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("Error saving conversation history: %v", err),
			}, nil
		}

		response = reply
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       response,
	}, nil
}

// TODO implement command pattern
func HandleCommands(messageData messageData.MessageData) (string, error) {
	// if the user sent command "j!clear", clear the conversation history
	if messageData.MessageContent == "j!clear" {
		// ensure user is an admin first
		if IsAdmin(messageData) {
			err := redisBroker.ClearConversationHistory(messageData.ChannelID)
			if err != nil {
				return "", err
			}

			return "Conversation history cleared.", nil
		} else {
			return "", fmt.Errorf("You are not authorized to clear the conversation history.")
		}
	}

	return "", nil
}

func IsAdmin(messageData messageData.MessageData) bool {
	// Implement your logic to check if the user is an admin
	// For example, you could check if the user's role is "admin" or has specific permissions

	// for now, just return true if the message author is "vonnycakes"
	return messageData.AuthorUsername == "vonnycakes"
}
