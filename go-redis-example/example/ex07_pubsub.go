package example

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const ex07Channel = "es_ch"

func Ex07(ctx context.Context) {
	pubSub := RedisCli.Subscribe(ctx, ex07Channel)
	defer func(mPubSub *redis.PubSub) {
		if err := mPubSub.Unsubscribe(ctx, ex07Channel); err != nil {
			log.Fatal(err)
		}
		_ = mPubSub.Close()
	}(pubSub)

	go func() {
		for i := 0; i < 5; i++ {
			RedisCli.Publish(ctx, ex07Channel, i)
		}
	}()

	for msg := range pubSub.Channel() {
		arcId, err := strconv.ParseInt(msg.Payload, 10, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("读取文章[%d]标题、正文，发送到ES更新索引\n", arcId)
	}
}
