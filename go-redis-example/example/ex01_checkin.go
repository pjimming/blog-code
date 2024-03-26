package example

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

const continueCheckKey = "cc_uid_%d"

func Ex01(ctx context.Context, params []string) {
	userID, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("参数错误：params = %+v, error = %v", params, err)
		panic(err)
	}
	ex01AddContinueDays(ctx, userID)
}

// 用户签到
func ex01AddContinueDays(ctx context.Context, userID int64) {
	key := ex01GetContinueCheckKey(userID)
	// 1. 签到天数+1
	if err := RedisCli.Incr(ctx, key).Err(); err != nil {
		err = fmt.Errorf("user[%d]签到失败, %v", userID, err)
		panic(err)
	}

	// 2. 设置签到时间为后天0点过期
	expAt := ex01BeginningOfDay().Add(48 * time.Hour)
	if err := RedisCli.ExpireAt(ctx, key, expAt).Err(); err != nil {
		panic(err)
	}

	// 3. 打印用户签到天数
	day, err := ex01GetUserCheckInDays(ctx, userID)
	if err != nil {
		panic(err)
	}
	log.Printf("User[%d]连续签到：%d天，过期时间：%s", userID, day, expAt.Format("2006-01-02 15:04:05"))
}

// 获取用户签到天数
func ex01GetUserCheckInDays(ctx context.Context, userID int64) (int64, error) {
	key := ex01GetContinueCheckKey(userID)
	days, err := RedisCli.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	daysInt, err := strconv.ParseInt(days, 10, 64)
	if err != nil {
		return 0, err
	}
	return daysInt, nil
}

// 获取今天0点时间
func ex01BeginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// 获取记录签到天数的key
func ex01GetContinueCheckKey(userID int64) string {
	return fmt.Sprintf(continueCheckKey, userID)
}
