package types

import (
	"time"

	"github.com/darkcloudlabs/hailstorm/pkg/runtime"
	"github.com/google/uuid"
)

type Deploy struct {
	ID        uuid.UUID        `json:"id"`
	AppID     uuid.UUID        `json:"appId"`
	Region    string           `json:"region"`
	Blob      []byte           `json:"-"`
	CreatedAT time.Time        `json:"createdAt"`
	Function  int              `json:"function"`
	Runtime   *runtime.Runtime `json:"-"`
}
