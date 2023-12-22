package types

import (
	"time"

	"github.com/google/uuid"
)

type Deploy struct {
	ID        uuid.UUID `json:"id"`
	AppID     uuid.UUID `json:"appID"`
	CreatedAT time.Time `json:"createdAt"`
}
