package redisBroker

import (
	"context"
	"encoding/json"
	"fmt"
	"kenrms/jubibot-response-lambda/messageData"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	println("Initializing Redis client")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "jubi-responder-cache-loq0u0.serverless.use1.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// Test connectivity with Ping.
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to ElastiCache Redis:", err)
		return
	}
	fmt.Println("Ping response:", pong)
}

func SaveConversationHistory(channelID string, history []messageData.MessageData) error {
	ctx := context.Background()
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, channelID, historyJSON, time.Minute*60*24).Err()

	return err
}

func GetConversationHistory(channelID string) ([]messageData.MessageData, error) {
	ctx := context.Background()
	historyJSON, err := redisClient.Get(ctx, channelID).Result()
	if err != nil {
		// if a key doesn't exist, ensure one is created
		if err == redis.Nil {
			emptyHistory := make([]messageData.MessageData, 0)
			err = SaveConversationHistory(channelID, emptyHistory)
			if err != nil {
				return nil, err
			}

			return emptyHistory, nil
		}

	}

	var history []messageData.MessageData
	err = json.Unmarshal([]byte(historyJSON), &history)
	if err != nil {
		return nil, err
	}

	return history, err
}

func ClearConversationHistory(channelID string) error {
	ctx := context.Background()
	err := redisClient.Del(ctx, channelID).Err()

	return err
}
