package types

import (
	"time"

	"github.com/google/uuid"
)

const APPVersion = "1"

type App struct {
	ID         uuid.UUID `json:"id"`
	Version    string    `json:"version"`
	Name       string    `json:"name"`
	ExposeHTTP int       `json:"exposeHTTP"`
	CreatedAt  time.Time `json:"createdAt"`
}
