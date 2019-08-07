package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"sync/atomic"
	"time"
)

const FILED_SP = "#"

var firstLogTime = ""
var logAtomaticId int64 = 0

// GetLogId get logid  for trace, return a new logid if ctx is nil
func GetLogId(ctx context.Context) string {
	if ctx != nil {
		logId, ok := ctx.Value("logId").(string)
		if ok {
			return logId
		}
	}

	return fmt.Sprintf("%s-%d", firstLogTime, atomic.AddInt64(&logAtomaticId, 1))
}

// logElapsed log func elapsed time, using in defer
func logElapsed(mark ...string) func() {
	start_at := time.Now()
	return func() {
		log.Printf("[DEBUG] [ELAPSED] %s elapsed %d ms\n", strings.Join(mark, " "), int64(time.Since(start_at)/time.Millisecond))
	}
}

// writeToFile write data to file
func writeToFile(filePath string, data interface{}) error {
	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Get current user fail,reason %s", err.Error())
		}
		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat old file error,reason %s", err.Error())
	}

	if !os.IsNotExist(err) {
		if fileInfo.IsDir() {
			return fmt.Errorf("old filepath is a dir,can not delete")
		}
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("delete old file error,reason %s", err.Error())
		}
	}

	jsonStr, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("json decode error,reason %s", err.Error())
	}

	return ioutil.WriteFile(filePath, []byte(jsonStr), 422)
}
