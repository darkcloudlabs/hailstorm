package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/google/uuid"
)

type CreateAppParams struct {
	Name string `json:"name"`
}

func (p CreateAppParams) validate() error {
	if len(p.Name) < 3 || len(p.Name) > 20 {
		return fmt.Errorf("name of the application should be longer then 3 and less than 20")
	}
	return nil
}

func (s *Server) handleCreateApp(w http.ResponseWriter, r *http.Request) error {
	var params CreateAppParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(ErrDecodeRequestBody))
	}
	if err := params.validate(); err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}
	app := types.App{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAT: time.Now(),
	}
	if err := s.store.CreateApp(&app); err != nil {
		return writeJSON(w, http.StatusBadRequest, ErrorResponse(err))
	}
	return writeJSON(w, http.StatusOK, app)
}
