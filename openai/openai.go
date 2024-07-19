package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kenrms/jubibot-response-lambda/constants"
	"kenrms/jubibot-response-lambda/messageData"
	"net/http"
	"os"
)

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	TopP        float64         `json:"top_p"`
	N           int             `json:"n"`
	Stop        []string        `json:"stop"`
	User        string          `json:"user"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
		Index        int           `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func GetReplyFromOpenAI(conversation []messageData.MessageData) (string, error) {
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	messages := []OpenAIMessage{
		{
			Role:    "system",
			Content: constants.OPEN_AI_SYSTEM_PROMPT,
		},
	}

	for _, message := range conversation {
		if message.AuthorUsername == constants.BOT_NAME {
			messages = append(messages, OpenAIMessage{
				Role:    "assistant",
				Content: message.MessageContent,
			})
		} else {
			var body string
			if message.ReferencedMessageID != "" {
				body = fmt.Sprintf("This message is in response to another message. Here is the response: \n\n --- \n\n "+
					"Author: %s\n"+
					"Message: %s\n"+
					"---\n\n"+
					"This is the message this this is responding to:\n\n"+
					"---\n\n"+
					"Reference Author: %s\n"+
					"Refereance Message: %s",
					message.AuthorUsername,
					message.MessageContent,
					message.ReferencedMessageAuthor,
					message.ReferencedMessageContent)
			} else {
				body = fmt.Sprintf("Author: %s\nMessage: %s", message.AuthorUsername, message.MessageContent)
			}

			messages = append(messages, OpenAIMessage{
				Role:    "user",
				Content: body,
			})
		}
	}

	// TODO get this configuration info from config API
	openAIRequest := OpenAIRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   2048,
		Temperature: 0.7,
		TopP:        1,
		N:           1,
		Stop:        nil,
		User:        "discord-bot",
	}

	requestBody, err := json.Marshal(openAIRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling OpenAI request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating OpenAI request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openAIAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending OpenAI request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading OpenAI response: %v", err)
	}

	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling OpenAI response: %v", err)
	}

	// if openAI responds with a BadRequest, reply with an error saying that the input was bad
	if resp.StatusCode == http.StatusBadRequest {
		return "Sorry, I couldn't understand your input. Please try again.", nil
	}

	if len(openAIResponse.Choices) > 0 {
		return openAIResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI")
}
