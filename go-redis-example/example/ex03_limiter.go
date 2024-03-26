package example

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"go-redis-example/common"
)

type Ex03Params struct {
}

const (
	ex03LimitKeyPreFix = "common_freq_limit" // 限流key前缀
	ex03MaxQPS         = 10                  // 限流次数
)

var (
	accessQueryNum = int32(0)
)

// 返回key格式为：comment_freq_limit-1669524458,
// 用来记录这1秒内的请求数量
func ex03LimitKey(currentTimeStamp time.Time) string {
	return fmt.Sprintf("%s-%d", ex03LimitKeyPreFix, currentTimeStamp.Unix())
}

// Ex03 简单限流
func Ex03(ctx context.Context) {
	eventLogger := common.NewConcurrentEventLog(ctx, 1000)
	// new一个并发执行器
	cInst := common.NewConcurrentRoutine(500, eventLogger)
	// 并发执行自定义函数
	cInst.Run(ctx, Ex03Params{}, ex03Work)
	// 输出日志
	eventLogger.PrintLogs()
	log.Printf("放行总数：%d", accessQueryNum)

	time.Sleep(1 * time.Second)
	fmt.Printf("\n------\n下一秒请求\n------\n")
	// 清空日志信息
	eventLogger = common.NewConcurrentEventLog(ctx, 1000)
	accessQueryNum = 0
	// new一个并发执行器
	cInst = common.NewConcurrentRoutine(10, eventLogger)
	// 并发执行用户自定义函数work
	cInst.Run(ctx, Ex03Params{}, ex03Work)
	// 按日志时间正序打印日志
	eventLogger.PrintLogs()
	log.Printf("放行总数：%d", accessQueryNum)
}

func ex03Work(ctx context.Context, cInstParams common.CInstParams) {
	routine := cInstParams.Routine
	eventLogger := cInstParams.ConcurrentEventLogger
	key := ex03LimitKey(time.Now())
	currentQPS, err := RedisCli.Incr(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	if currentQPS > ex03MaxQPS {
		// 超过流量限制，请求受限
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       common.LogFormat(routine, "被限流[%d]", currentQPS),
		})
		// sleep模拟业务耗时
		time.Sleep(50 * time.Millisecond)
		if err = RedisCli.Decr(ctx, key).Err(); err != nil {
			panic(err)
		}
	} else {
		// 流量放行
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       common.LogFormat(routine, "流量放行[%d]", currentQPS),
		})
		atomic.AddInt32(&accessQueryNum, 1)
		time.Sleep(20 * time.Millisecond)
	}
}
