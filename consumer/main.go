package main

import (
	"fmt"

	"github.com/felipeagger/go-redis-streams/consumer/handler"
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

	fmt.Printf("Initializing consumer on Stream: %v ...\n", streamName)

	// start fetch events
	events := eventFetcher(client)

	// start consume events
	consumeEvents(events, handler.HandlerFactory())

	quit := make(chan bool)
	<-quit
}

func consumeEvents(events chan event.Event, handlerFactory func(t event.Type) handler.Handler) {
	for {
		select {
		case event := <-events:
			h := handlerFactory(event.GetType())
			err := h.Handle(event)
			if err != nil {
				fmt.Printf("handle event error eventType:%v err:%v\n", event.GetType(), err)
			}
		}
	}
}

// start fetch new event starting from st.LatestEventID
func eventFetcher(client *redis.Client) chan event.Event {
	c := make(chan event.Event, 1000)
	start := "-"

	go func() {
		for {
			func() {

				redisRange, err := client.XRange(streamName, start, "+").Result()
				if err != nil {
					panic(err)
				}

				for _, stream := range redisRange {
					start = stream.ID

					tp := stream.Values["type"].(string)
					newEvent, _ := event.New(event.Type(tp))

					err = newEvent.UnmarshalBinary([]byte(stream.Values["data"].(string)))
					if err != nil {
						fmt.Printf("fail to unmarshal event:%v\n", stream.ID)
						return
					}

					//Delete stream from redis
					client.XDel(streamName, stream.ID)

					newEvent.SetID(stream.ID)
					c <- newEvent
				}

			}()
		}
	}()

	return c
}
