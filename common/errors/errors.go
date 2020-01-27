package errors

import (
	"encoding/json"
)

type Error struct {
	Errno  int32  `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func NewError(errno int32, errmsg string) error {
	return &Error{
		Errno:  errno,
		Errmsg: errmsg,
	}
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Errno = -1
		e.Errmsg = err
	}
	if e.Errno == 0 && e.Errmsg == "" {
		e.Errno = -1
		e.Errmsg = err
	}
	return e
}
