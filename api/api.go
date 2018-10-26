package api

import (
	"github.com/aweris/stp/internal/server"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ApiHandler struct {
	server  *server.Server
	router  *mux.Router
	timeout time.Duration
}

// TODO : add services
func NewHandler(s *server.Server) *ApiHandler {
	api := &ApiHandler{server: s, router: mux.NewRouter(), timeout: time.Second * 5}

	// initialize routes
	api.registerHealthCheck()
	api.registerInventoryRoutes()

	return api
}

func (ah *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.router.ServeHTTP(w, r)
}
