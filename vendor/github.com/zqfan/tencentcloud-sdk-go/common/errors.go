package common

import (
	"fmt"
)

type APIError struct {
	Code       string
	Message    string
	CodeNumber int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[APIError] Code=%s, Message=%s, CodeNumber=%d", e.Code, e.Message, e.CodeNumber)
}

func NewAPIError(code, message string, codeNumber int) error {
	return &APIError{
		Code:       code,
		Message:    message,
		CodeNumber: codeNumber,
	}
}
