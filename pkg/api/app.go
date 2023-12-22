package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateAppRequest struct {
	Name       string
	ExposeHTTP int
	Region     string
}

type CreateAppResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (p CreateAppRequest) validate() error {
	if len(p.Name) < 3 || len(p.Name) > 20 {
		return fmt.Errorf("application name should be min 3 and max 20 characters long")
	}
	if len(p.Region) < 3 || len(p.Region) > 20 {
		return fmt.Errorf("region should be min 3 and max 20 characters long")
	}
	return nil
}

func (s *Server) handleCreateApp(w http.ResponseWriter, r *http.Request) error {
	var params CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(fmt.Errorf("failed to parse request body")))
	}
	if err := params.validate(); err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}

	app := types.App{
		ID:         uuid.New(),
		Version:    types.APPVersion,
		Name:       params.Name,
		ExposeHTTP: params.ExposeHTTP,
		CreatedAt:  time.Now(),
	}

	if err := s.store.InsertApp(&app); err != nil {
		return writeJSON(w, http.StatusUnprocessableEntity, ErrorResponse(err))
	}

	resp := CreateAppResponse{
		ID:   app.ID,
		Name: app.Name,
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleCreateDeploy(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(fmt.Errorf("invalid ID format")))
	}
	app, err := s.store.GetApp(id)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}

	fmt.Println(app)
	return nil
}
