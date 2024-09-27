package deprovisioning

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-project/kyma-environment-broker/internal/kubeconfig"

	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"
	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var siCRD = []byte(`
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: serviceinstances.services.cloud.sap.com
spec:
  group: services.cloud.sap.com
  names:
    kind: ServiceInstance
    listKind: ServiceInstanceList
    plural: serviceinstances
    singular: serviceinstance
  scope: Namespaced
`)

var sbCRD = []byte(`
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: servicebindings.services.cloud.sap.com
spec:
  group: services.cloud.sap.com
  names:
    kind: ServiceBinding
    listKind: ServiceBindingList
    plural: servicebindings
    singular: servicebinding
  scope: Namespaced
`)

func TestRemoveServiceInstanceStep(t *testing.T) {
	t.Run("should remove all service instances and bindings from btp operator as part of trial suspension", func(t *testing.T) {
		// given
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		si := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "ServiceInstance",
			"metadata": map[string]interface{}{
				"name":      "test-instance",
				"namespace": "kyma-system",
			},
		}}
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)
		decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
		obj, gvk, err := decoder.Decode(siCRD, nil, nil)
		fmt.Println(gvk)
		require.NoError(t, err)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(obj, ns).Build()}
		clientProvider := kubeconfig.NewFakeK8sClientProvider(k8sCli)
		err = k8sCli.Create(context.TODO(), si)
		require.NoError(t, err)

		op := fixture.FixSuspensionOperationAsOperation(fixOperationID, fixInstanceID)
		op.State = "in progress"
		step := NewBTPOperatorCleanupStep(ms.Operations(), clientProvider)

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op, entry)

		// then
		assert.NoError(t, err)

		// given
		emptySI := &unstructured.Unstructured{}
		emptySI.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "services.cloud.sap.com",
			Version: "v1",
			Kind:    "ServiceInstance",
		})

		// then
		assert.True(t, k8sCli.cleanupInstances)
		assert.True(t, k8sCli.cleanupBindings)
	})

	t.Run("should skip btp-cleanup if not trial", func(t *testing.T) {
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		si := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "ServiceInstance",
			"metadata": map[string]interface{}{
				"name":      "test-instance",
				"namespace": "kyma-system",
			},
		}}
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)
		decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
		obj, gvk, err := decoder.Decode(siCRD, nil, nil)
		fmt.Println(gvk)
		require.NoError(t, err)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(obj, ns).Build()}
		err = k8sCli.Create(context.TODO(), si)
		require.NoError(t, err)

		op := fixture.FixSuspensionOperationAsOperation(fixOperationID, fixInstanceID)
		op.ProvisioningParameters.PlanID = broker.AWSPlanID
		op.State = "in progress"
		step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op, entry)

		// then
		assert.NoError(t, err)

		// given
		emptySI := &unstructured.Unstructured{}
		emptySI.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "services.cloud.sap.com",
			Version: "v1",
			Kind:    "ServiceInstance",
		})

		// then
		assert.False(t, k8sCli.cleanupInstances)
		assert.False(t, k8sCli.cleanupBindings)
	})

	t.Run("should skip btp-cleanup if not suspension", func(t *testing.T) {
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		si := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "ServiceInstance",
			"metadata": map[string]interface{}{
				"name":      "test-instance",
				"namespace": "kyma-system",
			},
		}}
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)
		decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
		obj, gvk, err := decoder.Decode(siCRD, nil, nil)
		fmt.Println(gvk)
		require.NoError(t, err)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(obj, ns).Build()}
		err = k8sCli.Create(context.TODO(), si)
		require.NoError(t, err)

		op := fixture.FixSuspensionOperationAsOperation(fixOperationID, fixInstanceID)
		op.State = "in progress"
		op.Temporary = false
		step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op, entry)

		// then
		assert.NoError(t, err)

		// given
		emptySI := &unstructured.Unstructured{}
		emptySI.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "services.cloud.sap.com",
			Version: "v1",
			Kind:    "ServiceInstance",
		})

		// then
		assert.False(t, k8sCli.cleanupInstances)
		assert.False(t, k8sCli.cleanupBindings)
	})
}

func TestBTPOperatorCleanupStep_SoftDelete(t *testing.T) {
	t.Run("should skip resources deletion when CRDs are missing", func(t *testing.T) {
		// given
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(ns).Build()}
		require.NoError(t, err)

		op := fixture.FixDeprovisioningOperation(fixOperationID, fixInstanceID)
		op.UserAgent = broker.AccountCleanupJob
		op.State = "in progress"
		step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op.Operation, entry)

		// then
		assert.NoError(t, err)
		assert.False(t, k8sCli.cleanupInstances)
		assert.False(t, k8sCli.cleanupBindings)
	})

	t.Run("should delete SI and skip SB deletion ", func(t *testing.T) {
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		si := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "ServiceInstance",
			"metadata": map[string]interface{}{
				"name":      "test-instance",
				"namespace": "kyma-system",
			},
		}}
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)
		decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
		obj, gvk, err := decoder.Decode(siCRD, nil, nil)
		fmt.Println(gvk)
		require.NoError(t, err)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(obj, ns).Build()}
		err = k8sCli.Create(context.TODO(), si)
		require.NoError(t, err)

		op := fixture.FixDeprovisioningOperation(fixOperationID, fixInstanceID)
		op.Temporary = true
		op.ProvisioningParameters.PlanID = broker.TrialPlanID
		op.UserAgent = broker.AccountCleanupJob
		op.State = "in progress"
		step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op.Operation, entry)

		// then
		assert.NoError(t, err)
		assert.True(t, k8sCli.cleanupInstances)
		assert.False(t, k8sCli.cleanupBindings)
	})

	t.Run("should delete SB and skip SI deletion ", func(t *testing.T) {
		log := logrus.New()
		ms := storage.NewMemoryStorage()
		sb := &unstructured.Unstructured{Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "ServiceBinding",
			"metadata": map[string]interface{}{
				"name":      "test-binding",
				"namespace": "kyma-system",
			},
		}}
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "kyma-system"}}

		scheme := internal.NewSchemeForTests(t)
		err := apiextensionsv1.AddToScheme(scheme)
		decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
		obj, gvk, err := decoder.Decode(sbCRD, nil, nil)
		fmt.Println(gvk)
		require.NoError(t, err)

		k8sCli := &fakeK8sClientWrapper{fake: fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(obj, ns).Build()}
		err = k8sCli.Create(context.TODO(), sb)
		require.NoError(t, err)

		op := fixture.FixDeprovisioningOperation(fixOperationID, fixInstanceID)
		op.UserAgent = broker.AccountCleanupJob
		op.State = "in progress"
		op.Temporary = true
		op.ProvisioningParameters.PlanID = broker.TrialPlanID
		step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))

		// when
		entry := log.WithFields(logrus.Fields{"step": "TEST"})
		_, _, err = step.Run(op.Operation, entry)

		// then
		assert.NoError(t, err)
		assert.False(t, k8sCli.cleanupInstances)
		assert.True(t, k8sCli.cleanupBindings)
	})
}

func TestBTPOperatorCleanupStep_NoKubeconfig(t *testing.T) {
	// given
	ms := storage.NewMemoryStorage()
	scheme := internal.NewSchemeForTests(t)
	// k8s client to an "empty" K8s
	k8sCli := fake.NewClientBuilder().WithScheme(scheme).Build()
	step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewK8sClientFromSecretProvider(k8sCli))
	op := fixture.FixDeprovisioningOperation(fixOperationID, fixInstanceID)
	op.State = "in progress"

	// when
	log := logrus.New()
	_, backoff, err := step.Run(op.Operation, log)

	// then
	assert.NoError(t, err)
	assert.Zero(t, backoff)
}

func TestBTPOperatorCleanupStep_NoRuntimeID(t *testing.T) {
	// given
	ms := storage.NewMemoryStorage()
	scheme := internal.NewSchemeForTests(t)
	// k8s client to an "empty" K8s
	k8sCli := fake.NewClientBuilder().WithScheme(scheme).Build()
	step := NewBTPOperatorCleanupStep(ms.Operations(), kubeconfig.NewFakeK8sClientProvider(k8sCli))
	op := fixture.FixDeprovisioningOperation(fixOperationID, fixInstanceID)
	op.State = "in progress"
	op.RuntimeID = ""

	// when
	log := logrus.New()
	_, backoff, err := step.Run(op.Operation, log)

	// then
	assert.NoError(t, err)
	assert.Zero(t, backoff)
}

type fakeK8sClientWrapper struct {
	fake             client.Client
	cleanupInstances bool
	cleanupBindings  bool
}

func (f *fakeK8sClientWrapper) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	return f.fake.Get(ctx, key, obj)
}

func (f *fakeK8sClientWrapper) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	if u, ok := list.(*unstructured.UnstructuredList); ok {
		switch u.Object["kind"] {
		case "ServiceBindingList":
			f.cleanupBindings = true
		case "ServiceInstanceList":
			f.cleanupInstances = true
		}
	}
	return f.fake.List(ctx, list, opts...)
}

func (f *fakeK8sClientWrapper) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return f.fake.Create(ctx, obj, opts...)
}

func (f *fakeK8sClientWrapper) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return f.fake.Delete(ctx, obj, opts...)
}

func (f *fakeK8sClientWrapper) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return f.fake.Update(ctx, obj, opts...)
}

func (f *fakeK8sClientWrapper) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return f.fake.Patch(ctx, obj, patch, opts...)
}

func (f *fakeK8sClientWrapper) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	if u, ok := obj.(*unstructured.Unstructured); ok {
		switch u.Object["kind"] {
		case "ServiceBinding":
			f.cleanupBindings = true
		case "ServiceInstance":
			f.cleanupInstances = true
		}
	}
	return f.fake.DeleteAllOf(ctx, obj, opts...)
}

func (f *fakeK8sClientWrapper) Status() client.StatusWriter {
	return f.fake.Status()
}

func (f *fakeK8sClientWrapper) Scheme() *runtime.Scheme {
	return f.fake.Scheme()
}

func (f *fakeK8sClientWrapper) RESTMapper() meta.RESTMapper {
	return f.fake.RESTMapper()
}

func (f *fakeK8sClientWrapper) SubResource(subresource string) client.SubResourceClient {
	return f.fake.SubResource(subresource)
}

type fakeProvisionerClient struct {
	empty bool
}

func newEmptyProvisionerClient() fakeProvisionerClient {
	return fakeProvisionerClient{true}
}

func (f fakeProvisionerClient) ProvisionRuntime(accountID, subAccountID string, config gqlschema.ProvisionRuntimeInput) (gqlschema.OperationStatus, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) DeprovisionRuntime(accountID, runtimeID string) (string, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) UpgradeRuntime(accountID, runtimeID string, config gqlschema.UpgradeRuntimeInput) (gqlschema.OperationStatus, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) UpgradeShoot(accountID, runtimeID string, config gqlschema.UpgradeShootInput) (gqlschema.OperationStatus, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) ReconnectRuntimeAgent(accountID, runtimeID string) (string, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) RuntimeOperationStatus(accountID, operationID string) (gqlschema.OperationStatus, error) {
	panic("not implemented")
}

func (f fakeProvisionerClient) RuntimeStatus(accountID, runtimeID string) (gqlschema.RuntimeStatus, error) {
	if f.empty {
		return gqlschema.RuntimeStatus{}, fmt.Errorf("not found")
	}
	kubeconfig := "sample fake kubeconfig"
	return gqlschema.RuntimeStatus{
		RuntimeConfiguration: &gqlschema.RuntimeConfig{
			Kubeconfig: &kubeconfig,
		},
	}, nil
}
