package api

import (
	"io"
	"net/http"
	"time"

	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateDeployParams struct {
	Region string `json:"region"`
}

func (s *Server) handleCreateDeploy(w http.ResponseWriter, r *http.Request) error {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}
	app, err := s.store.GetAppByID(appID)
	if err != nil {
		return writeJSON(w, http.StatusNotFound, ErrorResponse(err))
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return writeJSON(w, http.StatusNotFound, ErrorResponse(err))
	}

	deploy := types.Deploy{
		ID:        uuid.New(),
		AppID:     app.ID,
		Blob:      b,
		CreatedAT: time.Now(),
	}
	if err := s.store.CreateDeploy(&deploy); err != nil {
		return writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse(err))
	}

	return writeJSON(w, http.StatusOK, deploy)
}
