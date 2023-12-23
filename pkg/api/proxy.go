package api

import (
	"context"
	"net/http"
	"os"

	"github.com/darkcloudlabs/hailstorm/pkg/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(chi.URLParam(r, ("id")))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}
	// Get the app by ID

	b, err := os.ReadFile("testdata/app.wasm")
	if err != nil {
		return err
	}

	run, err := runtime.New(b)
	if err != nil {
		return err
	}

	if err := run.HandleHTTP(w, r); err != nil {
		return err
	}

	run.Close(context.Background())

	_ = id
	return nil
}
