package dbmodel

import (
	"time"
)

type BindingDTO struct {
	ID        string
	RuntimeID string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt *time.Time

	Kubeconfig         string
	ExpireationSeconds int
	GenerationMethod   string

	Version int
}
