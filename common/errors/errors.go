package errors

import (
	"github.com/micro/go-micro/v2/errors"
)

func newError(code int32, detail string) error {
	return errors.New("dmicro", detail, code)
}

func NewError(detail string) error {
	return errors.New("dmicro", detail, ErrCustom)
}
