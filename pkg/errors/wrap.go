package errors

import (
	"github.com/pkg/errors"
)

func WrapUnexpected(err error) error {
	return errors.Wrap(err, "unexpected error occurred")
}
