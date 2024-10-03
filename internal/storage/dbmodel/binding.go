package dbmodel

import (
	"time"
)

type BindingDTO struct {
	ID        string
	RuntimeID string

	CreatedAt time.Time

	Kubeconfig        string
	ExpirationSeconds int64
	GenerationMethod  string

	Version int
}
