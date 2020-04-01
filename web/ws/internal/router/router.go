package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

// 与 dd 不同的是，这里用了 mux 库。
func Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/say", say)
	r.HandleFunc("/join", join)

	return r
}
