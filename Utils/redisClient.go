package Utils

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"os"
)

func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	return client, err

}
