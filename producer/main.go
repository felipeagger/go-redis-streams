package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/felipeagger/go-redis-streams/packages/event"
	"github.com/felipeagger/go-redis-streams/packages/utils"
	"github.com/go-redis/redis/v7"
)

const (
	streamName = "events"
)

func main() {
	client, err := utils.NewRedisClient()
	if err != nil {
		panic(err)
	}

	generateEvent(client)
}

func generateEvent(client *redis.Client) {
	for i := 0; i < 10; i++ {

		eventType := []event.Type{event.ViewType, event.LikeType}[rand.Intn(2)]
		extra := []string{"test", "gopher", "streams"}[rand.Intn(3)]
		userID := uint64(rand.Intn(1000))

		strCMD := client.XAdd(&redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"type": string(eventType),
				"data": &event.ViewEvent{
					Base: &event.Base{
						Type:     eventType,
						DateTime: time.Now(),
					},
					UserID: userID,
					Extra:  extra,
				},
			},
		})

		newID, err := strCMD.Result()
		if err != nil {
			fmt.Printf("produce event error:%v\n", err)
		} else {
			fmt.Printf("produce event success Type:%v UserID:%v Extra:%v offset:%v\n",
				string(eventType), userID, extra, newID)
		}

	}
}
