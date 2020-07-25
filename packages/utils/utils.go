package utils

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
)

//NewRedisClient create a new instace of client redis
func NewRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	return client, err

}
