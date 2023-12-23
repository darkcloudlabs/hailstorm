package proxy

import (
	"net/http"

	"github.com/darkcloudlabs/hailstorm/pkg/store"
	"github.com/go-chi/chi/v5"
)

// Proxy is an HTTP server that will proxy the request to the corresponding function.

type Proxy struct {
	router *chi.Mux
	store  store.Store
}

func NewProxy() *Proxy {
	return &Proxy{
		router: chi.NewRouter(),
		// hardcoded for now
		store: store.NewMemoryStore(),
	}
}

func (p *Proxy) Listen(addr string) error {
	return http.ListenAndServe(addr, p.router)
}

// func (p *Proxy) initRoutes() {
// 	p.router.Handle("/", p.handleRequest)
// }

func (p *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {

}
