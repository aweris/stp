package api

import (
	"net/http"
)


func (ah *ApiHandler) registerHealthCheck() {
	sub := ah.router.PathPrefix("/ping").Subrouter()

	sub.HandleFunc("", ah.checkCheckHandler).Subrouter().Methods("GET")
}

func (ah *ApiHandler) checkCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}