package entity

import (
	"context"
	"errors"

	"github.com/otyang/go-pkg/datastore"
)

type IRepository interface {
	GetUsers(ctx context.Context) ([]*User, error)
}

var ErrNotFoundInDB = errors.New("no result found")

func IsErrNotFound(err error) error {
	if datastore.IsErrNotFound(err) {
		return ErrNotFoundInDB
	}
	return err
}
