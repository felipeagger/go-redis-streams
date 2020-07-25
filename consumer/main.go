package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/felipeagger/go-redis-streams/consumer/handler"
	"github.com/felipeagger/go-redis-streams/packages/event"
	"github.com/felipeagger/go-redis-streams/packages/utils"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

var (
	waitGrp       sync.WaitGroup
	client        *redis.Client
	start         string = ">"
	streamName    string = os.Getenv("STREAM")
	consumerGroup string = os.Getenv("GROUP")
	consumerName  string = uuid.NewV4().String()
)

func init() {
	var err error
	client, err = utils.NewRedisClient()
	if err != nil {
		panic(err)
	}

	createConsumerGroup()
}

func main() {
	fmt.Printf("Initializing Consumer:%v\nConsumerGroup: %v \nStream: %v\n",
		consumerName, consumerGroup, streamName)

	go consumeEvents()
	go consumePendingEvents()

	//Gracefully disconection
	chanOS := make(chan os.Signal)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)
	<-chanOS

	waitGrp.Wait()
	client.Close()
}

func createConsumerGroup() {

	if _, err := client.XGroupCreateMkStream(streamName, consumerGroup, "0").Result(); err != nil {

		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)
			panic(err)
		}

	}
}

// start consume events
func consumeEvents() {

	for {
		func() {
			fmt.Println("new round ", time.Now().Format(time.RFC3339))

			streams, err := client.XReadGroup(&redis.XReadGroupArgs{
				Streams:  []string{streamName, start},
				Group:    consumerGroup,
				Consumer: consumerName,
				Count:    10,
				Block:    0,
			}).Result()

			if err != nil {
				log.Printf("err on consume events: %+v\n", err)
				return
			}

			for _, stream := range streams[0].Messages {
				waitGrp.Add(1)
				go processStream(stream, false, handler.HandlerFactory())
			}
			waitGrp.Wait()
		}()
	}

}

func consumePendingEvents() {

	ticker := time.Tick(time.Second * 30)
	for {
		select {
		case <-ticker:

			func() {

				var streamsRetry []string
				pendingStreams, err := client.XPendingExt(&redis.XPendingExtArgs{
					Stream: streamName,
					Group:  consumerGroup,
					Start:  "0",
					End:    "+",
					Count:  10,
					//Consumer string
				}).Result()

				if err != nil {
					panic(err)
				}

				for _, stream := range pendingStreams {
					streamsRetry = append(streamsRetry, stream.ID)
				}

				if len(streamsRetry) > 0 {

					streams, err := client.XClaim(&redis.XClaimArgs{
						Stream:   streamName,
						Group:    consumerGroup,
						Consumer: consumerName,
						Messages: streamsRetry,
						MinIdle:  30 * time.Second,
					}).Result()

					if err != nil {
						log.Printf("err on process pending: %+v\n", err)
						return
					}

					for _, stream := range streams {
						waitGrp.Add(1)
						go processStream(stream, true, handler.HandlerFactory())
					}
					waitGrp.Wait()
				}

				fmt.Println("process pending streams at ", time.Now().Format(time.RFC3339))

			}()

		}

	}

}

func processStream(stream redis.XMessage, retry bool, handlerFactory func(t event.Type) handler.Handler) {
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
	err = h.Handle(newEvent, retry)
	if err != nil {
		fmt.Printf("error on process event:%v\n", newEvent)
		fmt.Println(err)
		return
	}

	//client.XDel(streamName, stream.ID)
	client.XAck(streamName, consumerGroup, stream.ID)

	//time.Sleep(2 * time.Second)
}
