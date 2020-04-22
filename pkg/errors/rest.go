package errors

import (
	"fmt"
	"net/http"
)

// RestError represents a Rest HTTP Error that can be returned from a controller
type RestError struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	OriginalError error  `json:"-"`
}

func (re *RestError) Error() string {
	return re.Message
}

type UnexpectStatusCodeError struct {
	UnexpectedCode int
	OriginalError  error
}

func (e *UnexpectStatusCodeError) Error() string {
	return fmt.Sprintf("The following Exit Code: %d, was unexpected", e.UnexpectedCode)
}

func (re *RestError) Unwrap() error {
	return re.OriginalError
}

func NotFoundf(format string, a ...interface{}) *RestError {
	return NotFound(fmt.Sprintf(format, a...))
}

func NotFound(message string) *RestError {
	return &RestError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func UnexpectedStatusCode(status int, err error) *UnexpectStatusCodeError {
	return &UnexpectStatusCodeError{
		UnexpectedCode: status,
		OriginalError:  err,
	}
}
