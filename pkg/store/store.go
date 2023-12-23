package store

import (
	"fmt"

	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/google/uuid"
)

type Store interface {
	CreateApp(*types.App) error
	CreateDeploy(*types.Deploy) error
	GetDeployByID(uuid.UUID) (*types.Deploy, error)
	GetAppByID(uuid.UUID) (*types.App, error)
}

type DeployStore interface {
	CreateDeply(*types.Deploy) error
	GetDeployByID(uuid.UUID) (*types.Deploy, error)
}

type MemoryStore struct {
	apps    map[uuid.UUID]*types.App
	deploys map[uuid.UUID]*types.Deploy
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		apps:    make(map[uuid.UUID]*types.App),
		deploys: make(map[uuid.UUID]*types.Deploy),
	}
}

func (s *MemoryStore) CreateApp(app *types.App) error {
	s.apps[app.ID] = app
	return nil
}

func (s *MemoryStore) GetAppByID(id uuid.UUID) (*types.App, error) {
	app, ok := s.apps[id]
	if !ok {
		return nil, fmt.Errorf("could not find app with id (%s)", id)
	}
	return app, nil
}

func (s *MemoryStore) CreateDeploy(deploy *types.Deploy) error {
	s.deploys[deploy.ID] = deploy
	return nil
}

func (s *MemoryStore) GetDeployByID(id uuid.UUID) (*types.Deploy, error) {
	deploy, ok := s.deploys[id]
	if !ok {
		return nil, fmt.Errorf("could not find deployment with id (%s)", id)
	}
	return deploy, nil
}
