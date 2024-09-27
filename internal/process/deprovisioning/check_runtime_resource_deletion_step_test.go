package deprovisioning

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/kyma-project/kyma-environment-broker/internal/logger"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCheckRuntimeResourceDeletionStep_ResourceNotExists(t *testing.T) {
	// given
	err := imv1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)
	op := fixture.FixDeprovisioningOperationAsOperation(fixOperationID, fixInstanceID)
	op.RuntimeResourceName = "runtime-name"
	op.KymaResourceNamespace = "kyma-ns"
	memoryStorage := storage.NewMemoryStorage()
	assert.NoError(t, memoryStorage.Operations().InsertOperation(op))
	kcpClient := fake.NewClientBuilder().Build()
	log := logger.NewLogDummy()

	// when
	step := NewCheckRuntimeResourceDeletionStep(memoryStorage.Operations(), kcpClient)
	_, backoff, err := step.Run(op, log)

	// then
	assert.NoError(t, err)
	assert.Zero(t, backoff)
}

func TestCheckRuntimeResourceDeletionStep_Run(t *testing.T) {
	// given
	err := imv1.AddToScheme(scheme.Scheme)
	assert.NoError(t, err)
	op := fixture.FixDeprovisioningOperationAsOperation(fixOperationID, fixInstanceID)
	op.RuntimeResourceName = "runtime-name"
	op.KymaResourceNamespace = "kyma-ns"
	memoryStorage := storage.NewMemoryStorage()
	assert.NoError(t, memoryStorage.Operations().InsertOperation(op))
	kcpClient := fake.NewClientBuilder().WithRuntimeObjects(fixRuntimeResource("kyma-ns", "runtime-name")).Build()
	log := logger.NewLogDummy()

	// when
	step := NewCheckRuntimeResourceDeletionStep(memoryStorage.Operations(), kcpClient)
	_, backoff, err := step.Run(op, log)

	// then
	assert.NoError(t, err)
	assert.NotZero(t, backoff)
}
