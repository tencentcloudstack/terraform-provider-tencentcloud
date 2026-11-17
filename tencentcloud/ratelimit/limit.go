package ratelimit

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/yangwenmai/ratelimit/simpleratelimit"
)

var (
	limitConfig = make(map[string]int64)

	limitContainer = make(map[string]*simpleratelimit.RateLimiter)

	locker sync.Mutex
)

func actionLimiter(namespace, action string) *simpleratelimit.RateLimiter {
	key := fmt.Sprintf("%s.%s", namespace, action)
	locker.Lock()

	limitNumber := limitConfig[key]
	if limitNumber == 0 {
		limitNumber = DefaultLimit

		if limitConfig[namespace] != 0 {
			limitNumber = limitConfig[namespace]
		}
		limitConfig[key] = limitNumber
	}

	if limitContainer[key] == nil {
		limitContainer[key] = simpleratelimit.New(int(limitNumber), time.Second)
	}

	limit := limitContainer[key]
	locker.Unlock()
	return limit
}

func ProCheck(namespace, action string) {
	key := fmt.Sprintf("%s.%s", namespace, action)
	limit := actionLimiter(namespace, action)
	old := time.Now()

	for limit.Limit() {
		//Prevent wake at same time
		sleepMs := 10 + rand.Intn(30)
		time.Sleep(time.Duration(sleepMs) * time.Microsecond)

		if time.Since(old) > 5*time.Minute {
			log.Printf("[WARN] %s wait too long, we try to release it", key)
			break
		}
	}

}

// ProCheckContext 保留限流键，取消或超时后不放行请求。
func ProCheckContext(ctx context.Context, namespace, action string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	limit := actionLimiter(namespace, action)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !limit.Limit() {
			return nil
		}
		timer := time.NewTimer(time.Duration(10+rand.Intn(30)) * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func CheckContext(ctx context.Context, action string) error {
	_, filePath, _, _ := runtime.Caller(1)
	items := strings.Split(filePath, `/`)
	items = strings.Split(items[len(items)-1], `\`)
	fileName := strings.TrimSuffix(items[len(items)-1], ".go")
	return ProCheckContext(ctx, fileName, action)
}

func Check(action string) {

	_, filePath, _, _ := runtime.Caller(1)

	items := strings.Split(filePath, `/`)
	items = strings.Split(items[len(items)-1], `\`)

	fileName := strings.TrimSuffix(items[len(items)-1], ".go")

	ProCheck(fileName, action)
}
