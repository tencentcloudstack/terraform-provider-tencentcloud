package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const DefaultLogMsg = "\n %v \n %v \n %s\n"

var firstLogTime = ""
var logAtomaticId int64 = 0

//debug log
func Dlog(ctx context.Context, title string, head, content interface{}) {
	debugEnv := os.Getenv("DEBUG")
	debugEnv = "tencentcloud"
	if strings.Contains(debugEnv, "tencentcloud") {
		logId, _ := ctx.Value("logId").(string)
		fmt.Printf(logId+DefaultLogMsg, title, content, debug.Stack())
		log.Printf(logId+DefaultLogMsg, title, content, debug.Stack())
	}
}

func GetLogId() string {
	if firstLogTime == "" {
		firstLogTime = fmt.Sprintf("%x", time.Now().Unix())
		log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
		logOut := &lumberjack.Logger{
			Filename:   "./app.log",
			MaxSize:    2 * 1024, //2G
			MaxBackups: 10,
			Compress:   true,
			LocalTime:  true,
		}
		log.SetOutput(logOut)
	}
	return fmt.Sprintf("%s-%d", firstLogTime, atomic.AddInt64(&logAtomaticId, 1))
}

func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

func writeToFile(filePath string, data interface{}) error {
	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	os.Remove(filePath)

	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}
	ioutil.WriteFile(filePath, []byte(out), 422)
	return nil
}
