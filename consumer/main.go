package main

import (
	"fmt"
	"sync"

	"github.com/felipeagger/go-redis-streams/consumer/handler"
	"github.com/felipeagger/go-redis-streams/packages/event"
	"github.com/felipeagger/go-redis-streams/packages/utils"
	"github.com/go-redis/redis/v7"
)

const (
	streamName = "events"
)

var (
	mutex sync.Mutex
	start string = "-"
)

func main() {
	client, err := utils.NewRedisClient()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Initializing consumer on Stream: %v ...\n", streamName)

	go consumeEvents(client)

	quit := make(chan bool)
	<-quit
}

// start consume events
func consumeEvents(client *redis.Client) {

	for {
		func() {

			redisRange, err := client.XRange(streamName, start, "+").Result()
			if err != nil {
				panic(err)
			}

			for _, stream := range redisRange {

				go processStream(stream, client, handler.HandlerFactory())

			}

		}()
	}

}

func processStream(stream redis.XMessage, client *redis.Client, handlerFactory func(t event.Type) handler.Handler) {

	typeEvent := stream.Values["type"].(string)
	newEvent, _ := event.New(event.Type(typeEvent))

	err := newEvent.UnmarshalBinary([]byte(stream.Values["data"].(string)))
	if err != nil {
		fmt.Printf("error on unmarshal stream:%v\n", stream.ID)
		return
	}

	newEvent.SetID(stream.ID)

	h := handlerFactory(newEvent.GetType())
	err = h.Handle(newEvent)
	if err != nil {
		fmt.Printf("error on process event:%v\n", newEvent)
		fmt.Println(err)
		return
	}

	//Delete stream from redis
	client.XDel(streamName, stream.ID)

	mutex.Lock()
	start = stream.ID
	mutex.Unlock()
}
