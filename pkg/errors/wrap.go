package errors

import (
	"database/sql/driver"
	"fmt"

	"github.com/pkg/errors"
)

const eveErrorMessage = "eve error"

type cause interface {
	Cause() error
}

type eveError interface {
	IsEveError() bool
}

type txError struct {
	TxError       error
	OriginalError error
}

func (tx txError) Error() string {
	return tx.OriginalError.Error()
}

func (tx txError) Unwrap() error {
	return tx.OriginalError
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	if ee, ok := err.(eveError); ok && ee.IsEveError() {
		return err
	} else if _, ok := err.(cause); ok {
		return err
	} else {
		return errors.Wrap(err, eveErrorMessage)
	}
}

func Wrapf(format string, a ...interface{}) error {
	return Wrap(fmt.Errorf(format, a...))
}

func WrapTx(tx driver.Tx, err error) error {
	if tx == nil {
		return Wrap(err)
	}
	txErr := tx.Rollback()
	if txErr != nil {
		err = txError{
			TxError:       txErr,
			OriginalError: err,
		}
		return Wrap(txErr)
	}
	return Wrap(err)
}
