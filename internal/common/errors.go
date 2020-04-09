package common

import (
	"fmt"
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

func (re *RestError) Unwrap() error {
	return re.OriginalError
}

func ErrTrap(message, methodName string, err error) error {
	//if errors.Cause(err).Error() == SQLNOROWS {
	//	wrappedError := &RestError{
	//		Code:          404,
	//		Message:       message,
	//		OriginalError: err,
	//	}
	//	return errors.Wrap(wrappedError, fmt.Sprintf("Service.%s", methodName))
	//}
	return fmt.Errorf("service.%s: %w", methodName, err)
}
