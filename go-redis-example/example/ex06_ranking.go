package example

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const ex06RankKey = "ex06_rank_zset"

type ex06ItemScore struct {
	ItemName string
	Score    float64
}

// Ex06 排行榜
// go run main.go Ex06 init // 初始化积分
// go run main.go Ex06 rev_order // 输出完整榜单
// go run main.go Ex06 order_page 1 // 逆序分页输出，page=1
// go run main.go Ex06 get_rank user2 // 获取user2的排名
// go run main.go Ex06 get_score user2 // 获取user2的分数
// go run main.go Ex06 add_user_score user2 10 // 为user2设置为10分
// zadd ex06_rank_zset 15 andy
// zincrby ex06_rank_zset -9 andy // andy 扣9分，排名掉到最后一名
func Ex06(ctx context.Context, args []string) {
	arg1 := args[0]
	switch arg1 {
	case "init":
		ex06Init(ctx)
	case "rev_order":
		ex06GetOrderListAll(ctx)
	case "order_page":
		pageSize := int64(2)
		if len(args[1]) > 0 {
			offset, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				panic(err)
			}
			ex06GetOrderListByPage(ctx, offset, pageSize)
		}
	case "get_rank":
		ex06GetUserRankByName(ctx, args[1])
	case "get_score":
		ex06GetUserScoreByName(ctx, args[1])
	case "add_user_score":
		if len(args) < 3 {
			fmt.Printf("参数错误，可能是缺少需要增加的分值。eg：go run main.go  Ex06 add_user_score user2 10\n")
			return
		}
		score, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			panic(err)
		}
		ex06AddUserScore(ctx, args[1], score)
	default:
		panic("unsupported type")
	}
}

func ex06Init(ctx context.Context) {
	initList := []redis.Z{
		{Member: "user1", Score: 10},
		{Member: "user2", Score: 232},
		{Member: "user3", Score: 129},
		{Member: "user4", Score: 232},
	}
	// 清空榜单
	if err := RedisCli.Del(ctx, ex06RankKey).Err(); err != nil {
		panic(err)
	}

	nums, err := RedisCli.ZAdd(ctx, ex06RankKey, initList...).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("初始化榜单Item数量：%d\n", nums)
}

// 获取全部榜单
// 榜单逆序输出
// ZRANGE ex06_rank_zset +inf -inf BYSCORE  rev WITHSCORES
// 正序输出
// ZRANGE ex06_rank_zset 0 -1 WITHSCORES
func ex06GetOrderListAll(ctx context.Context) {
	resList, err := RedisCli.ZRevRangeWithScores(ctx, ex06RankKey, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("\n榜单：")
	for i, z := range resList {
		fmt.Printf("第%d名，name=%s, score=%.2f\n", i+1, z.Member, z.Score)
	}
}

// 分页获取榜单
func ex06GetOrderListByPage(ctx context.Context, page, pageSize int64) {
	// zrange ex06_rank_zset 300 0 byscore rev limit 1 2 withscores // 取300分到0分之间的排名
	// zrange ex06_rank_zset -inf +inf byscore withscores 正序输出
	// zrange ex06_rank_zset +inf -inf byscore rev WITHSCORES 逆序输出所有排名
	// zrange ex06_rank_zset +inf -inf byscore rev limit 0 2 withscores 逆序分页输出排名
	offset := int((page - 1) * pageSize)
	zRangeArgs := redis.ZRangeArgs{
		Key:     ex06RankKey,
		ByScore: true,
		Rev:     true,
		Start:   "-inf",
		Stop:    "+inf",
		Offset:  int64(offset),
		Count:   pageSize,
	}
	resList, err := RedisCli.ZRangeArgsWithScores(ctx, zRangeArgs).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("榜单(page=%d, pageSize=%d)\n", page, pageSize)
	for i, z := range resList {
		rank := i + 1 + offset
		fmt.Printf("第%d名 %s\t%.2f\n", rank, z.Member, z.Score)
	}
}

// 获取用户排名
func ex06GetUserRankByName(ctx context.Context, name string) {
	rank, err := RedisCli.ZRevRank(ctx, ex06RankKey, name).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, rank=%d\n", name, rank+1)
}

// 获取用户分数信息
func ex06GetUserScoreByName(ctx context.Context, name string) {
	score, err := RedisCli.ZScore(ctx, ex06RankKey, name).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, score=%.2f\n", name, score)
}

// ex06AddUserScore 增加用户分数
func ex06AddUserScore(ctx context.Context, name string, score float64) {
	num, err := RedisCli.ZIncrBy(ctx, ex06RankKey, score, name).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("name=%s, add_score=%.2f, score=%.2f\n", name, score, num)
}
