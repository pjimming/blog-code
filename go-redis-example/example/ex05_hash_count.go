package example

import (
	"context"
	"fmt"
	"log"
	"strconv"
)

const ex05UserCountKey = "ex05_user_count"

// Ex05 hash数据结果的运用（参考掘金应用）
// go run main.go Ex05 init 初始化用户计数值
// go run main.go Ex05 get 1556564194374926  // 打印用户(1556564194374926)的所有计数值
// go run main.go Ex05 incr_like 1556564194374926 // 点赞数+1
// go run main.go Ex05 incr_collect 1556564194374926 // 点赞数+1
// go run main.go Ex05 decr_like 1556564194374926 // 点赞数-1
// go run main.go Ex05 decr_collect 1556564194374926 // 点赞数-1
func Ex05(ctx context.Context, args []string) {
	if len(args) <= 0 {
		panic("args can't be empty")
	}
	arg1 := args[0]
	switch arg1 {
	case "init":
		Ex05InitUserCount(ctx)
	case "get":
		userID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err)
		}
		Ex05GetUserCount(ctx, userID)
	//case "incr_like":
	//	userID, err := strconv.ParseInt(args[1], 10, 64)
	//	if err != nil {
	//		panic(err)
	//	}
	//	IncrByUserLike(ctx, userID)
	//case "incr_collect":
	//	userID, err := strconv.ParseInt(args[1], 10, 64)
	//	if err != nil {
	//		panic(err)
	//	}
	//	IncrByUserCollect(ctx, userID)
	//case "decr_like":
	//	userID, err := strconv.ParseInt(args[1], 10, 64)
	//	if err != nil {
	//		panic(err)
	//	}
	//	DecrByUserLike(ctx, userID)
	//case "decr_collect":
	//	userID, err := strconv.ParseInt(args[1], 10, 64)
	//	if err != nil {
	//		panic(err)
	//	}
	//	DecrByUserCollect(ctx, userID)
	default:
		panic("do not support now...")
	}
}

func Ex05InitUserCount(ctx context.Context) {
	pipe := RedisCli.Pipeline()
	userCounters := []map[string]interface{}{
		{"user_id": "1556564194374926", "got_digg_count": 10693, "got_view_count": 2238438, "followee_count": 176, "follower_count": 9895, "follow_collect_set_count": 0, "subscribe_tag_count": 95},
		{"user_id": "1111", "got_digg_count": 19, "got_view_count": 4},
		{"user_id": "2222", "got_digg_count": 1238, "follower_count": 379},
	}

	for _, counter := range userCounters {
		uid, err := strconv.ParseInt(counter["user_id"].(string), 10, 64)
		if err != nil {
			panic(err)
		}
		key := ex05GetUserCounterKey(uid)
		if err = pipe.Del(ctx, key).Err(); err != nil {
			panic(err)
		}
		if err = pipe.HMSet(ctx, key, counter).Err(); err != nil {
			panic(err)
		}

		log.Printf("设置uid[%d], key=%s", uid, key)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		// 再执行一次
		if _, err = pipe.Exec(ctx); err != nil {
			panic(err)
		}
	}
}

// ex05GetUserCounterKey 获取用户计数的key
func ex05GetUserCounterKey(userID int64) string {
	return fmt.Sprintf("%s_%d", ex05UserCountKey, userID)
}

func Ex05GetUserCount(ctx context.Context, userID int64) {
	pipe := RedisCli.Pipeline()
	pipe.HGetAll(ctx, ex05GetUserCounterKey(userID))
	results, err := RedisCli.HGetAll(ctx, ex05GetUserCounterKey(userID)).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("User[%d]:\n", userID)
	for k, v := range results {
		fmt.Printf("%s: %s\n", k, v)
	}
}
