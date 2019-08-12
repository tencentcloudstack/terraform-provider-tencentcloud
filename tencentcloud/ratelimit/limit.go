package ratelimit

import (
	"fmt"
	"github.com/yangwenmai/ratelimit/simpleratelimit"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	limitConfig = make(map[string]int64)

	limitContainer = make(map[string]*simpleratelimit.RateLimiter)

	locker sync.Mutex
)

func ProCheck(namespace, action string) {

	key := fmt.Sprintf("%s.%s", namespace, action)

	var limit *simpleratelimit.RateLimiter

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

	limit = limitContainer[key]
	locker.Unlock()

	old := time.Now()

	for limit.Limit() {
		//Prevent wake at same time
		sleepMs := 10 + rand.Intn(30)
		time.Sleep(time.Duration(sleepMs) * time.Microsecond)

		if time.Now().Sub(old) > 5*time.Minute {
			log.Printf("[WARN] %s wait too long, we try to release it", key)
		}
	}

}

func Check(action string) {

	_, filePath, _, _ := runtime.Caller(1)

	items := strings.Split(filePath, `/`)
	items = strings.Split(items[len(items)-1], `\`)

	fileName := strings.TrimSuffix(items[len(items)-1], ".go")

	ProCheck(fileName, action)
}
