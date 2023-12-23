package types

import (
	"time"

	"github.com/google/uuid"
)

type App struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAT time.Time `json:"createAt"`
}
