package api

import (
	"context"
	"net/http"

	"github.com/darkcloudlabs/hailstorm/pkg/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) error {
	deployID, err := uuid.Parse(chi.URLParam(r, ("id")))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}

	deploy, err := s.store.GetDeployByID(deployID)
	if err != nil {
		return writeJSON(w, http.StatusNotFound, ErrorResponse(err))
	}

	run, err := runtime.New(deploy.Blob)
	if err != nil {
		return err
	}

	if err := run.HandleHTTP(w, r, deploy.Function); err != nil {
		return err
	}

	run.Close(context.Background())

	return nil
}
