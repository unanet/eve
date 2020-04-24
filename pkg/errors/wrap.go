package errors

import (
	"github.com/pkg/errors"
)

const eveErrorMessage = "eve error"

type cause interface {
	Cause() error
}

type eveError interface {
	IsEveError() bool
}

func Wrap(err error) error {
	if ee, ok := err.(eveError); ok && ee.IsEveError() {
		return err
	} else if _, ok := err.(cause); ok {
		return err
	} else {
		return errors.Wrap(err, eveErrorMessage)
	}
}
