package deprovisioning

import (
	"context"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CheckRuntimeResourceDeletionStep struct {
	operationManager *process.OperationManager
	kcpClient        client.Client
}

func NewCheckRuntimeResourceDeletionStep(operations storage.Operations, kcpClient client.Client) *CheckRuntimeResourceDeletionStep {
	return &CheckRuntimeResourceDeletionStep{
		operationManager: process.NewOperationManager(operations),
		kcpClient:        kcpClient,
	}
}

func (step *CheckRuntimeResourceDeletionStep) Name() string {
	return "Check_RuntimeResource_Deletion"
}

func (step *CheckRuntimeResourceDeletionStep) Run(operation internal.Operation, logger logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	namespace := operation.KymaResourceNamespace
	if namespace == "" {
		logger.Warnf("namespace for Kyma resource not specified, setting 'kcp-system'")
		namespace = "kcp-system"
	}
	resourceName := operation.RuntimeResourceName
	if resourceName == "" {
		logger.Infof("Runtime resource name is empty, using runtime-id")
		resourceName = operation.RuntimeID
	}
	if resourceName == "" {
		logger.Infof("Empty runtime ID, skipping")
		return operation, 0, nil
	}

	runtime := &imv1.Runtime{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: namespace,
		},
	}

	err := step.kcpClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      resourceName,
	}, runtime)

	if err == nil {
		logger.Infof("Runtime resource still exists")
		//TODO: extract the timeout as a configuration setting
		return step.operationManager.RetryOperationWithoutFail(operation, step.Name(), "Runtime resource still exists", 20*time.Second, 15*time.Minute, logger, nil)
	}

	if !errors.IsNotFound(err) {
		if meta.IsNoMatchError(err) {
			logger.Info("No CRD installed, skipping")
			return operation, 0, nil
		}

		logger.Warnf("unable to check Runtime resource existence: %s", err)
		return step.operationManager.RetryOperationWithoutFail(operation, step.Name(), "unable to check Runtime resource existence", backoffForK8SOperation, timeoutForK8sOperation, logger, err)
	}

	return step.operationManager.UpdateOperation(operation, func(op *internal.Operation) {
		op.RuntimeResourceName = ""
	}, logger)
}
