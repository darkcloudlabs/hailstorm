package store

import (
	"fmt"

	"github.com/darkcloudlabs/hailstorm/pkg/types"
	"github.com/google/uuid"
)

type AppStore interface {
	InsertApp(*types.App) error
	GetApp(uuid.UUID) (*types.App, error)
}

type MemoryAppStore struct {
	data map[uuid.UUID]*types.App
}

func NewMemoryAppStore() *MemoryAppStore {
	return &MemoryAppStore{
		data: make(map[uuid.UUID]*types.App),
	}
}

func (s *MemoryAppStore) InsertApp(app *types.App) error {
	s.data[app.ID] = app
	return nil
}

func (s *MemoryAppStore) GetApp(id uuid.UUID) (*types.App, error) {
	app, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("app with id (%s) not found", id)
	}
	return app, nil
}

type Store interface {
	AppStore
}

type MemoryStore struct {
	AppStore
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		AppStore: NewMemoryAppStore(),
	}
}
