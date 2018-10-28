package api

import (
	"github.com/aweris/stp/initialize"
	"net/http"
)

func (ah *ApiHandler) registerDemoHandler() {
	sub := ah.router.PathPrefix("/demo").Subrouter()

	sub.HandleFunc("", ah.loadTestData).Subrouter().Methods("POST")
}

func (ah *ApiHandler) loadTestData(w http.ResponseWriter, r *http.Request) {

	//Loading Test data
	initialize.LoadTestData(ah.server)

	w.WriteHeader(http.StatusNoContent)
}
