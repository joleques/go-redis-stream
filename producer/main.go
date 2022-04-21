package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/joleques/go-redis-stream/Utils"
	"math/rand"
	"os"
)

var (
	streamName string = os.Getenv("STREAM")
	client     *redis.Client
)

func init() {
	var err error
	client, err = Utils.NewRedisClient()
	if err != nil {
		panic(err)
	}
}

func main() {
	generateEvent()
}

func generateEvent() {
	var userID uint64 = 0
	for i := 0; i < 10; i++ {

		userID++

		message := []string{"Golang producer one!", "Golang producer two", "Golang producer three"}[rand.Intn(3)]

		newID, err := produceMsg(map[string]interface{}{
			"type": "string",
			"data": message,
		})

		checkError(err, newID, userID)
	}
}

func produceMsg(event map[string]interface{}) (string, error) {

	return client.XAdd(&redis.XAddArgs{
		Stream: streamName,
		Values: event,
	}).Result()
}

func checkError(err error, newID string, userID uint64) {
	if err != nil {
		fmt.Printf("produce event error:%v\n", err)
	} else {
		fmt.Printf("produce event success UserID:%v offset:%v\n", userID, newID)
	}
}
