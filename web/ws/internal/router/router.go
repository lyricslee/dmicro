package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/say", say)
	r.HandleFunc("/join", join)

	return r
}
