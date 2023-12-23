package api

import (
	"net/http"
)

type CreateDeployParams struct {
	Region string `json:"region"`
}

func (s *Server) handleCreateDeploy(w http.ResponseWriter, r *http.Request) error {
	return nil
}
