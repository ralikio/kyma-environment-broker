package steps

import (
	"context"
	"fmt"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewCheckRuntimeResourceStep(os storage.Operations, k8sClient client.Client, kimConfig broker.KimConfig, runtimeResourceStepTimeout time.Duration) *checkRuntimeResource {
	return &checkRuntimeResource{
		k8sClient:                  k8sClient,
		operationManager:           process.NewOperationManager(os),
		kimConfig:                  kimConfig,
		runtimeResourceStepTimeout: runtimeResourceStepTimeout,
	}
}

type checkRuntimeResource struct {
	k8sClient                  client.Client
	kimConfig                  broker.KimConfig
	operationManager           *process.OperationManager
	runtimeResourceStepTimeout time.Duration
}

func (_ *checkRuntimeResource) Name() string {
	return "Check_RuntimeResource"
}

func (s *checkRuntimeResource) Run(operation internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if !s.kimConfig.IsDrivenByKim(broker.PlanNamesMapping[operation.ProvisioningParameters.PlanID]) {
		log.Infof("Only provisioner is controlling provisioning process, skipping")
		return operation, 0, nil
	}

	runtime, err := s.GetRuntimeResource(operation.RuntimeID, operation.KymaResourceNamespace)
	if err != nil {
		log.Errorf("unable to get Runtime resource %s/%s", operation.KymaResourceNamespace, operation.RuntimeID)
		return s.operationManager.RetryOperation(operation, "unable to get Runtime resource", err, time.Second, 10*time.Second, log)
	}

	// check status
	state := runtime.Status.State
	log.Infof("Runtime resource state: %s", state)
	if state != imv1.RuntimeStateReady {
		if time.Since(operation.UpdatedAt) > s.runtimeResourceStepTimeout {
			description := fmt.Sprintf("Waiting for Runtime resource (%s/%s) ready state timeout.", operation.KymaResourceNamespace, operation.RuntimeID)
			log.Error(description)
			log.Infof("Runtime resource status: %v, timeout: %v", runtime.Status, s.runtimeResourceStepTimeout)
			return s.operationManager.OperationFailed(operation, description, nil, log)
		}
		return operation, 10 * time.Second, nil
	}
	return operation, 0, nil
}

func (s *checkRuntimeResource) GetRuntimeResource(name string, namespace string) (*imv1.Runtime, error) {
	runtime := imv1.Runtime{}
	err := s.k8sClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &runtime)
	if err != nil {
		return nil, err
	}
	return &runtime, nil
}
