package service

import (
	"gitlab.unanet.io/devops/eve/internal/data"
	"gitlab.unanet.io/devops/go/pkg/errors"
)

func CheckForNotFoundError(err error) error {
	if dErr, ok := err.(data.NotFoundError); ok {
		return errors.NotFoundf(dErr.Error())
	}
	return errors.Wrap(err)
}
