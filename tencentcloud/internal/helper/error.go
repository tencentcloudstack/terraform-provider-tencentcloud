package helper

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

type frame struct {
	frames   [3]uintptr
	function *string
	file     *string
	line     *int
}

func newFrame() frame {
	var s frame
	runtime.Callers(2, s.frames[:])
	return s
}

func (f frame) location() (function, file string, line int) {
	// quick path
	if f.function != nil {
		return *f.function, *f.file, *f.line
	}

	frames := runtime.CallersFrames(f.frames[:])
	if _, ok := frames.Next(); !ok {
		f.function = String("")
		f.file = String("")
		f.line = Int(0)

		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		f.function = String("")
		f.file = String("")
		f.line = Int(0)

		return "", "", 0
	}

	f.function = &fr.Function
	f.file = &fr.File
	f.line = &fr.Line

	return fr.Function, fr.File, fr.Line
}

type Error struct {
	Id        string
	RequestId string
	Cause     error
	frame     frame
	msg       string
	args      []interface{}
}

func (e Error) Error() string {
	if strings.ToUpper(os.Getenv("TF_LOG")) != "" {
		return e.debugError()
	}

	return e.error()
}

func (e Error) debugError() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("id: %s, request-id: %s\n", e.Id, e.RequestId))

	if e.msg != "" {
		sb.WriteString(fmt.Sprintf(e.msg, e.args...))
		sb.WriteRune('\n')
	}

	fn, file, line := e.frame.location()

	sb.WriteString(fmt.Sprintf("\t%s\n\t\t%s:%d\n\t- %v", fn, file, line, e.Cause))

	return sb.String()
}

func (e Error) error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("id: %s, request-id: %s\n", e.Id, e.RequestId))

	if e.msg != "" {
		sb.WriteString(fmt.Sprintf(e.msg, e.args...))
		sb.WriteRune('\n')
	}

	sb.WriteString(e.Cause.Error())

	return sb.String()
}

// if cause is *sdkErrors.TencentCloudSDKError, will use sdk error request-id
func WrapErrorf(cause error, id, requestId, msg string, args ...interface{}) error {
	if cause == nil {
		return nil
	}

	if sdkErr, ok := cause.(*sdkErrors.TencentCloudSDKError); ok && requestId == "" {
		requestId = sdkErr.RequestId
	}

	frame := newFrame()

	return Error{
		Id:        id,
		RequestId: requestId,
		Cause:     cause,
		frame:     frame,
		msg:       msg,
		args:      args,
	}
}

func WrapError(cause error, id, requestId string) error {
	if cause == nil {
		return nil
	}

	if sdkErr, ok := cause.(*sdkErrors.TencentCloudSDKError); ok && requestId == "" {
		requestId = sdkErr.RequestId
	}

	frame := newFrame()

	return Error{
		Id:        id,
		RequestId: requestId,
		Cause:     cause,
		frame:     frame,
	}
}

func UnwarpSDKError(err error) *sdkErrors.TencentCloudSDKError {
	var result *sdkErrors.TencentCloudSDKError
	if errors.As(err, &result) {
		return result
	}
	return nil
}
