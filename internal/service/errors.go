package service

import (
	"github.com/unanet/eve/internal/data"
	"github.com/unanet/go/pkg/errors"
)

func CheckForNotFoundError(err error) error {
	if dErr, ok := err.(data.NotFoundError); ok {
		return errors.NotFoundf(dErr.Error())
	}
	return errors.Wrap(err)
}
