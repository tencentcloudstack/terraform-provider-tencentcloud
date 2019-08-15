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

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/likexian/gokit/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const FILED_SP = "#"

var contextNil context.Context = nil

var logFirstTime = ""
var logAtomaticId int64 = 0

// readRetryTimeout is read retry timeout
const readRetryTimeout = 3 * time.Minute

// writeRetryTimeout is write retry timeout
const writeRetryTimeout = 5 * time.Minute

// retryableErrorCode is retryable error code
var retryableErrorCode = []string{
	// client
	"ClientError.HttpStatusCodeError",
	// commom
	"FailedOperation",
	"InternalError",
	"TradeUnknownError",
	"RequestLimitExceeded",
	"ResourceInUse",
	"ResourceInsufficient",
	"ResourceUnavailable",
	// cbs
	"ResourceBusy",
}

func init() {
	logFirstTime = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

// getLogId get logid  for trace, return a new logid if ctx is nil
func getLogId(ctx context.Context) string {
	if ctx != nil {
		logId, ok := ctx.Value("logId").(string)
		if ok {
			return logId
		}
	}

	return fmt.Sprintf("%s-%d", logFirstTime, atomic.AddInt64(&logAtomaticId, 1))
}

// logElapsed log func elapsed time, using in defer
func logElapsed(mark ...string) func() {
	start_at := time.Now()
	return func() {
		log.Printf("[DEBUG] [ELAPSED] %s elapsed %d ms\n", strings.Join(mark, " "), int64(time.Since(start_at)/time.Millisecond))
	}
}

// retryError returns retry error
func retryError(err error) *resource.RetryError {
	if isErrorRetryable(err) {
		return resource.RetryableError(err)
	}

	return resource.NonRetryableError(err)
}

// isErrorRetryable returns whether error is retryable
func isErrorRetryable(err error) bool {
	e, ok := err.(*errors.TencentCloudSDKError)
	if !ok {
		log.Printf("[CRITAL] NonRetryable error: %v", err)
		return false
	}

	longCode := e.Code
	if assert.IsContains(retryableErrorCode, longCode) {
		log.Printf("[CRITAL] Retryable error: %s", e.Error())
		return true
	}

	if strings.Contains(longCode, ".") {
		shortCode := strings.Split(longCode, ".")[0]
		if assert.IsContains(retryableErrorCode, shortCode) {
			log.Printf("[CRITAL] Retryable error: %s", e.Error())
			return true
		}
	}

	log.Printf("[CRITAL] NonRetryable error: %s", e.Error())

	return false
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

	return ioutil.WriteFile(filePath, jsonStr, 0422)
}
