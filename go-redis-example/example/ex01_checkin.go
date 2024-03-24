package example

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const continueCheckKey = "cc_uid_%d"

func Ex01(ctx context.Context, params []string) error {
	userID, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		err = errors.Errorf("参数错误：params = %+v, error = %v", params, err)
		return err
	}
	return addContinueDays(ctx, userID)
}

func addContinueDays(ctx context.Context, userID int64) error {
	key := getContinueCheckKey(userID)
	// 1. 签到天数+1
	if err := RedisCli.Incr(ctx, key).Err(); err != nil {
		err = errors.Errorf("user[%d]签到失败, %v", userID, err)
		return err
	}

	// 2. 设置签到时间为后天0点过期
	expAt := beginningOfDay().Add(48 * time.Hour)
	if err := RedisCli.ExpireAt(ctx, key, expAt).Err(); err != nil {
		return err
	}

	// 3. 打印用户签到天数
	day, err := getUserCheckInDays(ctx, userID)
	if err != nil {
		return err
	}
	log.Printf("User[%d]连续签到：%d天，过期时间：%s", userID, day, expAt.Format("2006-01-02 15:04:05"))
	return nil
}

func getUserCheckInDays(ctx context.Context, userID int64) (int64, error) {
	key := getContinueCheckKey(userID)
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
func beginningOfDay() time.Time {
	now := time.Now()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func getContinueCheckKey(userID int64) string {
	return fmt.Sprintf(continueCheckKey, userID)
}
