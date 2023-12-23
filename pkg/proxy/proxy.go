package proxy

import (
	"context"
	"net/http"

	"github.com/darkcloudlabs/hailstorm/pkg/runtime"
	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Proxy is an HTTP server that will proxy the request to the corresponding function.

type Proxy struct {
	router *chi.Mux
	store  store.Store
}

func New(store store.Store) *Proxy {
	return &Proxy{
		router: chi.NewRouter(),
		store:  store,
	}
}

func (p *Proxy) Listen(addr string) error {
	p.initRoutes()
	return http.ListenAndServe(addr, p.router)
}

func (p *Proxy) initRoutes() {
	p.router.Handle("/{id}", http.HandlerFunc(p.handleRequest))
}

func (p *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	deployID, err := uuid.Parse(chi.URLParam(r, ("id")))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	deploy, err := p.store.GetDeployByID(deployID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	run, err := runtime.New(deploy.Blob)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	defer run.Close(context.Background())

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	if err := run.HandleHTTP(w, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}
