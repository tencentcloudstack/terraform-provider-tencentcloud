package common

import (
	"fmt"
)

type APIError struct {
	Code       string
	Message    string
	CodeNumber int
	RequestId  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[APIError] Code=%s, Message=%s, CodeNumber=%d, RequestId=%s", e.Code, e.Message, e.CodeNumber, e.RequestId)
}

func NewAPIError(code, message string, codeNumber int, requestId string) error {
	return &APIError{
		Code:       code,
		Message:    message,
		CodeNumber: codeNumber,
		RequestId:  requestId,
	}
}
