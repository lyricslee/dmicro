package util

import (
	"encoding/json"
	"net/http"

	"github.com/micro/go-micro/v2/errors"
)

func MakeErrBody(err error) (code int, body string) {
	e := errors.Parse(err.Error())

	//m := make(map[string]interface{})
	//m["errno"] = ce.Errno
	//m["errmsg"] = ce.Errmsg
	//m["t"] = time.Now().UnixNano()

	//if ce.Errno == -1 {
	//	code = http.StatusInternalServerError
	//} else {
	//	code = 499
	//}

	code = 499

	b, _ := json.Marshal(e)
	body = string(b)
	return
}

func MakeBody(obj interface{}) (code int, body string) {
	if obj == nil {
		obj = make(map[string]interface{})
	}

	code = http.StatusOK
	b, _ := json.Marshal(obj)
	body = string(b)

	return
}
