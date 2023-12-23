package api

import (
	"net/http"

	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	store  store.Store
}

func NewServer(store store.Store) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Listen(addr string) error {
	s.initRouter()
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) initRouter() {
	s.router = chi.NewRouter()
	s.router.Get("/status", handleStatus)
	s.router.Post("/app", makeAPIHandler(s.handleCreateApp))
	s.router.Post("/app/{id}/deploy", makeAPIHandler(s.handleCreateDeploy))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
