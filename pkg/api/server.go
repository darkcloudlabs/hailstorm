package api

import (
	"log/slog"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/cluster"
	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router  *chi.Mux
	store   store.Store
	cluster *cluster.Cluster
}

func NewServer(c *cluster.Cluster, store store.Store) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			cluster: c,
			store:   store,
		}
	}
}

func (s *Server) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		_ = msg
		s.start(c)
	}
}

func (s *Server) start(c *actor.Context) {
	s.initRouter()
	go func() {
		if err := http.ListenAndServe(":3000", s.router); err != nil {
			slog.Error("HTTP listen", "err", err)
		}
	}()
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
