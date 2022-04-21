package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/joleques/go-redis-stream/Utils"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
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
	client, err = Utils.NewRedisClient()
	if err != nil {
		panic(err)
	}
	createConsumerGroup()
}

func createConsumerGroup() {

	if _, err := client.XGroupCreateMkStream(streamName, consumerGroup, "0").Result(); err != nil {

		if !strings.Contains(fmt.Sprint(err), "BUSYGROUP") {
			fmt.Printf("Error on create Consumer Group: %v ...\n", consumerGroup)
			panic(err)
		}

	}
}
func consumeEvents() {

	for {
		func() {
			fmt.Println("Iniciando o consumo ", time.Now().Format(time.RFC3339))

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
				go processStream(stream)
			}
			waitGrp.Wait()
		}()
	}
}

func processStream(stream redis.XMessage) {
	defer waitGrp.Done()
	value := stream.Values["data"].(string)
	fmt.Println("Mensagem processada:", value)
	//client.XDel(streamName, stream.ID)
	client.XAck(streamName, consumerGroup, stream.ID)
	//time.Sleep(2 * time.Second)
}

func main() {
	fmt.Println("Initializing Consumer:", consumerName)
	fmt.Println("ConsumerGroup:", consumerGroup)
	fmt.Println("Stream:", streamName)

	go consumeEvents()

	//Gracefully disconection
	chanOS := make(chan os.Signal)
	signal.Notify(chanOS, syscall.SIGINT, syscall.SIGTERM)
	<-chanOS
	waitGrp.Wait()
	client.Close()
}
