package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/felipeagger/go-redis-streams/consumer/handler"
	"github.com/felipeagger/go-redis-streams/packages/event"
	"github.com/felipeagger/go-redis-streams/packages/utils"
	"github.com/go-redis/redis/v7"
)

const (
	streamName    = "events"
	consumerGroup = "consumersOne"
	consumerName  = "consumerOne"
)

var (
	waitGrp sync.WaitGroup
	mutex   sync.Mutex
	start   string = ">" //"0" //"0-0" //"-"
)

func main() {
	client, err := utils.NewRedisClient()
	if err != nil {
		panic(err)
	}

	createConsumerGroup(client)

	fmt.Printf("Initializing consumerGroup: %v on Stream: %v ...\n", consumerGroup, streamName)
	go consumeEvents(client)

	//Gracefully end
	chanOS := make(chan os.Signal)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)
	<-chanOS

	waitGrp.Wait()
	client.Close()
}

func createConsumerGroup(client *redis.Client) {

	if _, err := client.XGroupCreateMkStream(streamName, consumerGroup, "0").Result(); err != nil {

		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)
			panic(err)
		}

	}
}

// start consume events
func consumeEvents(client *redis.Client) {

	for {
		func() {

			/*
				redisRange, err := client.XRange(streamName, start, "+").Result()
				if err != nil {
					panic(err)
				}
			*/

			streams, err := client.XReadGroup(&redis.XReadGroupArgs{
				Streams:  []string{streamName, start},
				Group:    consumerGroup,
				Consumer: consumerName,
				//NoAck:    true,
				Count: 10,
				Block: 0, // Wait for new messages without a timeout.
			}).Result()
			if err != nil {
				log.Printf("err: %+v\n", err)
				return
			}

			fmt.Println("new round")

			for _, stream := range streams[0].Messages {

				waitGrp.Add(1)
				go processStream(stream, client, handler.HandlerFactory())

			}

		}()
	}

}

func processStream(stream redis.XMessage, client *redis.Client, handlerFactory func(t event.Type) handler.Handler) {
	defer waitGrp.Done()

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

		client.XClaimJustID(&redis.XClaimArgs{
			Stream:   streamName,
			Group:    consumerGroup,
			Consumer: "consumerTwo",
			Messages: []string{stream.ID},
			//MinIdle:  5 * time.Second,
		})
		return
	}

	//Delete stream from redis
	//client.XDel(streamName, stream.ID)
	client.XAck(streamName, consumerGroup, stream.ID)

	/*
		mutex.Lock()
		start = stream.ID
		//fmt.Printf("new start: %v \n", start)
		mutex.Unlock()
	*/
}
