package dbmodel

import (
	"time"
)

// InstanceFilter holds the filters when querying Instances
type BindingFilter struct {
	PageSize   int
	Page       int
	RuntimeIDs []string
}

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
