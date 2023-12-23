package api

import (
	"net/http"

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

	if err := deploy.Runtime.HandleHTTP(w, r, deploy.Function); err != nil {
		return err
	}

	return nil
}
