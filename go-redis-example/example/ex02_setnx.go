package example

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-redis-example/common"
)

const (
	resourceKey = "syncKey"              // 分布式锁的key
	expTime     = 800 * time.Millisecond // 锁的过期时间，避免死锁
)

type Ex02Params struct {
}

// Ex02 只是体验SetNX的特性，不是高可用的分布式锁实现
// 该实现存在的问题:
// (1) 业务超时解锁，导致并发问题。业务执行时间超过锁超时时间
// (2) redis主备切换临界点问题。主备切换后，A持有的锁还未同步到新的主节点时，B可在新主节点获取锁，导致并发问题。
// (3) redis集群脑裂，导致出现多个主节点
func Ex02(ctx context.Context) {
	eventLogger := common.NewConcurrentEventLog(ctx, 32)
	// new一个并发执行器
	cInst := common.NewConcurrentRoutine(10, eventLogger)
	// 并发执行自定义work
	cInst.Run(ctx, Ex02Params{}, ex02Work)
	// 按时间顺序输出日志
	eventLogger.PrintLogs()
}

func ex02Work(ctx context.Context, cInstParams common.CInstParams) {
	routine := cInstParams.Routine
	eventLogger := cInstParams.ConcurrentEventLogger
	defer ex02ReleaseLock(ctx, routine, eventLogger)
	for {
		// 1. 尝试获取锁
		acquired, err := RedisCli.SetNX(ctx, resourceKey, routine, expTime).Result()
		if err != nil {
			err = fmt.Errorf("[%s] error routine[%d], %v", time.Now().Format(time.RFC3339Nano), routine, err)
			eventLogger.Append(common.EventLog{
				EventTime: time.Now(),
				Log:       err.Error(),
			})
			panic(err)
		}

		if acquired {
			// 2. 成功获取
			eventLogger.Append(common.EventLog{
				EventTime: time.Now(),
				Log:       fmt.Sprintf("[%s] routine[%d] 获取锁", time.Now().Format(time.RFC3339Nano), routine),
			})
			// 3. 模拟业务
			time.Sleep(10 * time.Millisecond)
			eventLogger.Append(common.EventLog{
				EventTime: time.Now(),
				Log:       fmt.Sprintf("[%s] routine[%d] 完成业务逻辑", time.Now().Format(time.RFC3339Nano), routine),
			})
			return
		}
		// 没有获取到锁，等待后重试
		time.Sleep(100 * time.Millisecond)
	}
}

func ex02ReleaseLock(ctx context.Context, routine int, eventLogger *common.ConcurrentEventLogger) {
	routineMark, _ := RedisCli.Get(ctx, resourceKey).Result()
	if strconv.FormatInt(int64(routine), 10) != routineMark {
		// 其它协程误删lock
		panic(fmt.Sprintf("del err lock[%s] can not del by [%d]", routineMark, routine))
	}
	result, err := RedisCli.Del(ctx, resourceKey).Result()
	if result == 1 {
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       fmt.Sprintf("[%s] routine[%d] 释放锁", time.Now().Format(time.RFC3339Nano), routine),
		})
	} else {
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       fmt.Sprintf("[%s] routine[%d] no lock to del", time.Now().Format(time.RFC3339Nano), routine),
		})
	}
	if err != nil {
		err = fmt.Errorf("[%s] error routine=%d, %v", time.Now().Format(time.RFC3339Nano), routine, err)
		eventLogger.Append(common.EventLog{
			EventTime: time.Now(),
			Log:       err.Error(),
		})
		panic(err)
	}
}
