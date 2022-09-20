package tencentcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"gopkg.in/yaml.v2"
)

const FILED_SP = "#"

var contextNil context.Context = nil

type contextLogId string

const logIdKey = contextLogId("logId")

const (
	PROVIDER_READ_RETRY_TIMEOUT  = "TENCENTCLOUD_READ_RETRY_TIMEOUT"
	PROVIDER_WRITE_RETRY_TIMEOUT = "TENCENTCLOUD_WRITE_RETRY_TIMEOUT"
	PROVIDER_WAIT_READ_TIMEOUT   = "TENCENTCLOUD_WAIT_READ_TIMEOUT"

	SWEEPER_NEED_PROTECT            = "SWEEPER_NEED_PROTECT"
	TENCENTCLOUD_COMMON_TIME_LAYOUT = "2006-01-02 15:04:05"
)

var logFirstTime = ""
var logAtomicId int64 = 0

// readRetryTimeout is read retry timeout
//const readRetryTimeout = 3 * time.Minute
var readRetry = getEnvDefault(PROVIDER_READ_RETRY_TIMEOUT, 3)
var readRetryTimeout = time.Duration(readRetry) * time.Minute

// writeRetryTimeout is write retry timeout
//const writeRetryTimeout = 5 * time.Minute
var writeRetry = getEnvDefault(PROVIDER_WRITE_RETRY_TIMEOUT, 5)
var writeRetryTimeout = time.Duration(writeRetry) * time.Minute

// writeRetryTimeout is write retry timeout
//const writeRetryTimeout = 5 * time.Minute
var waitRead = getEnvDefault(PROVIDER_WAIT_READ_TIMEOUT, 1)
var waitReadTimeout = time.Duration(waitRead) * time.Second

//const writeRetryTimeout = 5 * time.Minute
var needProtect = getEnvDefault(SWEEPER_NEED_PROTECT, 0)

// InternalError common internalError, do not add in retryableErrorCode,
// because when some product return this error, retry won't fix anything.
const InternalError = "InternalError"

// retryableErrorCode is retryable error code
var retryableErrorCode = []string{
	// client
	"ClientError.NetworkError",
	"ClientError.HttpStatusCodeError",
	// common
	"FailedOperation",
	"TradeUnknownError",
	"RequestLimitExceeded",
	"ResourceInUse",
	"ResourceUnavailable",
	// cbs
	"ResourceBusy",
}

func init() {
	logFirstTime = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

func getEnvDefault(key string, defVal int) int {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	int, err := strconv.Atoi(val)
	if err != nil {
		panic("TENCENTCLOUD_READ_RETRY_TIMEOUT or TENCENTCLOUD_WRITE_RETRY_TIMEOUT must be int.")
	}
	return int
}

// string to time.Time
func stringTotime(t string) time.Time {
	template := TENCENTCLOUD_COMMON_TIME_LAYOUT
	stamp, _ := time.ParseInLocation(template, t, time.Local)
	return stamp
}

func parseTimeFromCommonLayout(t *string) time.Time {
	if t == nil {
		return time.Time{}
	}
	result, err := time.Parse(TENCENTCLOUD_COMMON_TIME_LAYOUT, *t)
	if err != nil {
		return time.Time{}
	}
	return result
}

// getLogId get logId for trace, return a new logId if ctx is nil
func getLogId(ctx context.Context) string {
	if ctx != nil {
		logId, ok := ctx.Value(logIdKey).(string)
		if ok {
			return logId
		}
	}

	return fmt.Sprintf("%s-%d", logFirstTime, atomic.AddInt64(&logAtomicId, 1))
}

// logElapsed log func elapsed time, using in defer
func logElapsed(mark ...string) func() {
	startAt := time.Now()
	return func() {
		log.Printf("[DEBUG] [ELAPSED] %s elapsed %d ms\n", strings.Join(mark, " "), int64(time.Since(startAt)/time.Millisecond))
	}
}

// for Provider produced inconsistent result after apply
func inconsistentCheck(d *schema.ResourceData, meta interface{}) func() {
	oldJson, _ := json.Marshal(d.State())
	return func() {
		newJson, _ := json.Marshal(d.State())
		if !reflect.DeepEqual(oldJson, newJson) {
			log.Printf("[Resource id %s data changes after reading old:%s, new:%s", d.Id(), oldJson, newJson)
		}
	}
}

// retryError returns retry error
func retryError(err error, additionRetryableError ...string) *resource.RetryError {
	switch realErr := errors.Cause(err).(type) {
	case *sdkErrors.TencentCloudSDKError:
		if isExpectError(realErr, retryableErrorCode) {
			log.Printf("[CRITAL] Retryable defined error: %v", err)
			return resource.RetryableError(err)
		}

		if len(additionRetryableError) > 0 {
			if isExpectError(realErr, additionRetryableError) {
				log.Printf("[CRITAL] Retryable addition error: %v", err)
				return resource.RetryableError(err)
			}
		}

	default:
	}

	log.Printf("[CRITAL] NonRetryable error: %v", err)
	return resource.NonRetryableError(err)
}

// isExpectError returns whether error is expected error
func isExpectError(err error, expectError []string) bool {
	e, ok := err.(*sdkErrors.TencentCloudSDKError)
	if !ok {
		return false
	}

	longCode := e.Code
	if IsContains(expectError, longCode) {
		return true
	}

	if strings.Contains(longCode, ".") {
		shortCode := strings.Split(longCode, ".")[0]
		if IsContains(expectError, shortCode) {
			return true
		}
	}

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

func CheckNil(object interface{}, fields map[string]string) (nilFields []string) {
	// if object is a pointer, get value which object points to
	object = reflect.Indirect(reflect.ValueOf(object)).Interface()

	for i := 0; i < reflect.TypeOf(object).NumField(); i++ {
		fieldName := reflect.TypeOf(object).Field(i).Name

		if realName, ok := fields[fieldName]; ok {
			if realName == "" {
				realName = fieldName
			}

			if reflect.ValueOf(object).Field(i).IsNil() {
				nilFields = append(nilFields, realName)
			}
		}
	}

	return
}

// BuildTagResourceName builds the Tencent Cloud specific name of a resource description.
// The format is `qcs:project_id:service_type:region:account:resource`.
// For more information, go to https://cloud.tencent.com/document/product/598/10606.
func BuildTagResourceName(serviceType, resourceType, region, id string) string {
	switch serviceType {
	case "cos":
		return fmt.Sprintf("qcs::%s:%s:uid/:%s/%s", serviceType, region, resourceType, id)

	default:
		return fmt.Sprintf("qcs::%s:%s:uin/:%s/%s", serviceType, region, resourceType, id)
	}
}

// IsContains returns whether value is within array
func IsContains(array interface{}, value interface{}) bool {
	vv := reflect.ValueOf(array)
	if vv.Kind() == reflect.Ptr || vv.Kind() == reflect.Interface {
		if vv.IsNil() {
			return false
		}
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Slice:
		for i := 0; i < vv.Len(); i++ {
			if reflect.DeepEqual(value, vv.Index(i).Interface()) {
				return true
			}
		}
		return false
	case reflect.Map:
		s := vv.MapKeys()
		for i := 0; i < len(s); i++ {
			if reflect.DeepEqual(value, s[i].Interface()) {
				return true
			}
		}
		return false
	case reflect.String:
		ss := reflect.ValueOf(value)
		switch ss.Kind() {
		case reflect.String:
			return strings.Contains(vv.String(), ss.String())
		}
		return false
	default:
		return reflect.DeepEqual(array, value)
	}
}

func FindIntListIndex(list []int, elem int) int {
	for i, v := range list {
		if v == elem {
			return i
		}
	}
	return -1
}

func GetListIncrement(o []int, n []int) (result []int, err error) {
	result = append(result, n...)
	if len(o) > len(n) {
		err = fmt.Errorf("new list elem count %d less than old: %d", len(n), len(o))
		return
	}
	for _, v := range o {
		index := FindIntListIndex(result, v)
		if index == -1 {
			err = fmt.Errorf("elem %d not exist", v)
			return
		}
		if index+1 >= len(result) {
			result = result[:index]
		} else {
			result = append(result[:index], result[index+1:]...)
		}
	}
	return
}

func GetListDiffs(o []int, n []int) (adds []int, lacks []int) {
	fillArr := func(arr []int, v int, count int) []int {
		for i := 0; i < count; i++ {
			arr = append(arr, v)
		}
		return arr
	}
	diffs := map[int]int{}
	for _, v := range o {
		diffs[v] -= 1
	}
	for _, v := range n {
		diffs[v] += 1
	}
	log.Printf("DIFFS: %v", diffs)
	for num, count := range diffs {
		if count < 0 {
			lacks = fillArr(lacks, num, -count)
		} else if count > 0 {
			adds = fillArr(adds, num, count)
		}
	}
	return
}

//GoRoutine Limit
type GoRoutineLimit struct {
	Count int
	Chan  chan struct{}
}

func NewGoRoutine(num int) *GoRoutineLimit {
	return &GoRoutineLimit{
		Count: num,
		Chan:  make(chan struct{}, num),
	}
}

func (g *GoRoutineLimit) Run(f func()) {
	g.Chan <- struct{}{}
	go func() {
		f()
		<-g.Chan
	}()
}

// YamlParser yaml syntax Parser
func YamlParser(config string) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(config), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func YamlToBase64(config string) string {
	m := []byte(config)
	encodedStr := base64.StdEncoding.EncodeToString(m)
	return encodedStr
}

func Base64ToYaml(config string) (string, error) {
	yamlConfig, err := base64.StdEncoding.DecodeString(config)
	if err != nil {
		return "", err
	}
	return string(yamlConfig), nil
}
