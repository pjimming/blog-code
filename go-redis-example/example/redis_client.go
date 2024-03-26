package example

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisCli *redis.Client

func init() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Your Redis Address
		Password: "123456",         // Your Redis Password
	})
	if err := RedisCli.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
