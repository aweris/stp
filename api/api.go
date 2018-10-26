package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ApiHandler struct {
	router  *mux.Router
	timeout time.Duration
}

// TODO : add services
func NewHandler() *ApiHandler {
	api := &ApiHandler{router: mux.NewRouter(), timeout: time.Second * 5}

	// initialize routes
	api.registerHealthCheck()

	return api
}

func (ah *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.router.ServeHTTP(w, r)
}