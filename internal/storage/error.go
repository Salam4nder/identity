package storage

import "errors"

func UserNotFoundErr() error {
	return errors.New("user not found")
}
