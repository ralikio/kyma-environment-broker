package provisioning

import (
	"fmt"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal/broker"

	"github.com/sirupsen/logrus"

	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/provisioner"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
)

// CheckRuntimeStep checks if the SKR is provisioned
type CheckRuntimeStep struct {
	provisionerClient   provisioner.Client
	operationManager    *process.OperationManager
	provisioningTimeout time.Duration
	kimConfig           broker.KimConfig
}

func NewCheckRuntimeStep(os storage.Operations,
	provisionerClient provisioner.Client,
	provisioningTimeout time.Duration,
	kimConfig broker.KimConfig) *CheckRuntimeStep {
	return &CheckRuntimeStep{
		provisionerClient:   provisionerClient,
		operationManager:    process.NewOperationManager(os),
		provisioningTimeout: provisioningTimeout,
		kimConfig:           kimConfig,
	}
}

var _ process.Step = (*CheckRuntimeStep)(nil)

func (s *CheckRuntimeStep) Name() string {
	return "Check_Runtime"
}

func (s *CheckRuntimeStep) Run(operation internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if operation.RuntimeID == "" {
		log.Errorf("Runtime ID is empty")
		return s.operationManager.OperationFailed(operation, "Runtime ID is empty", nil, log)
	}
	return s.checkRuntimeStatus(operation, log.WithField("runtimeID", operation.RuntimeID))
}

func (s *CheckRuntimeStep) checkRuntimeStatus(operation internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if s.kimConfig.IsDrivenByKimOnly(broker.PlanNamesMapping[operation.ProvisioningParameters.PlanID]) {
		log.Infof("KIM is driving the process for plan %s, skipping", broker.PlanNamesMapping[operation.ProvisioningParameters.PlanID])
		return operation, 0, nil
	}

	if time.Since(operation.UpdatedAt) > s.provisioningTimeout {
		log.Infof("operation has reached the time limit: updated operation time: %s", operation.UpdatedAt)
		return s.operationManager.OperationFailed(operation, fmt.Sprintf("operation has reached the time limit: %s", s.provisioningTimeout), nil, log)
	}

	if operation.ProvisionerOperationID == "" {
		msg := "Operation does not contain Provisioner Operation ID"
		log.Error(msg)
		return s.operationManager.OperationFailed(operation, msg, nil, log)
	}

	status, err := s.provisionerClient.RuntimeOperationStatus(operation.ProvisioningParameters.ErsContext.GlobalAccountID, operation.ProvisionerOperationID)
	if err != nil {
		log.Errorf("call to provisioner RuntimeOperationStatus failed: %s", err.Error())
		return operation, 1 * time.Minute, nil
	}
	log.Infof("call to provisioner returned %s status", status.State.String())

	switch status.State {
	case gqlschema.OperationStateSucceeded:
		return operation, 0, nil
	case gqlschema.OperationStateInProgress:
		return operation, 20 * time.Second, nil
	case gqlschema.OperationStatePending:
		return operation, 20 * time.Second, nil
	case gqlschema.OperationStateFailed:
		lastErr := provisioner.OperationStatusLastError(status.LastError)
		return s.operationManager.OperationFailed(operation, "provisioner client returns failed status", lastErr, log)
	}

	lastErr := provisioner.OperationStatusLastError(status.LastError)
	return s.operationManager.OperationFailed(operation, fmt.Sprintf("unsupported provisioner client status: %s", status.State.String()), lastErr, log)
}
