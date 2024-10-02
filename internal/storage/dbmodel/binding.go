package dbmodel

import (
	"time"
)

const (
// InstanceSucceeded        InstanceState = "succeeded"
// InstanceFailed           InstanceState = "failed"
// InstanceError            InstanceState = "error"
// InstanceProvisioning     InstanceState = "provisioning"
// InstanceDeprovisioning   InstanceState = "deprovisioning"
// InstanceUpgrading        InstanceState = "upgrading"
// InstanceUpdating         InstanceState = "updating"
// InstanceDeprovisioned    InstanceState = "deprovisioned"
// InstanceNotDeprovisioned InstanceState = "notDeprovisioned"
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
