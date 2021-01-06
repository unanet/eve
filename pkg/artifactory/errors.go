package artifactory

import (
	"fmt"
)

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Errors []Status `json:"errors,omitempty"` // Individual errors
}

// Status is the individual error provided by the API
type Status struct {
	Status  int    `json:"status"`  // Validation error status code
	Message string `json:"message"` // Message describing the error. Errors with Code == "custom" will always have this set.
}

func (e *Status) Error() string {
	return fmt.Sprintf("%d error caused by %s", e.Status, e.Message)
}

func (r ErrorResponse) Error() string {
	return fmt.Sprintf("Artifactory Errors: %+v", r.Errors)
}

type ServiceUnavailableError struct {
	message string
}

func (e ServiceUnavailableError) Error() string {
	return e.message
}

func (e ServiceUnavailableError) IsEveError() bool {
	return true
}

func ServiceUnavailableErrorf(format string, a ...interface{}) ServiceUnavailableError {
	return ServiceUnavailableError{
		message: fmt.Sprintf(format, a...),
	}
}

type NotFoundError struct {
	message string
}

func (e NotFoundError) Error() string {
	return e.message
}

func (e NotFoundError) IsEveError() bool {
	return true
}

func NotFoundErrorf(format string, a ...interface{}) NotFoundError {
	return NotFoundError{
		message: fmt.Sprintf(format, a...),
	}
}

type InvalidRequestError struct {
	message string
}

func (e InvalidRequestError) Error() string {
	return e.message
}

func (e InvalidRequestError) IsEveError() bool {
	return true
}

func InvalidRequestErrorf(format string, a ...interface{}) InvalidRequestError {
	return InvalidRequestError{
		message: fmt.Sprintf(format, a...),
	}
}
