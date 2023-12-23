package types

import (
	"time"

	"github.com/google/uuid"
)

type Deploy struct {
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"createdAt"`
}
