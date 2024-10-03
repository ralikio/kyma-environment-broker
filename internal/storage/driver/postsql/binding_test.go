package postsql_test

import (
	"testing"

	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinding(t *testing.T) {

	t.Run("Should create binding and delete binding", func(t *testing.T) {
		storageCleanup, brokerStorage, err := GetStorageForDatabaseTests()
		require.NoError(t, err)
		require.NotNil(t, brokerStorage)
		defer func() {
			err := storageCleanup()
			assert.NoError(t, err)
		}()

		// given
		testBindingId := "test"
		fixedBinding := fixture.FixBinding(testBindingId)

		err = brokerStorage.Bindings().Insert(&fixedBinding)
		assert.NoError(t, err)

		// when
		createdBinding, err := brokerStorage.Bindings().Get(testBindingId)

		// then
		assert.NoError(t, err)
		assert.Equal(t, fixedBinding.ID, createdBinding.ID)
		assert.NotNil(t, createdBinding.InstanceID)
		assert.Equal(t, fixedBinding.InstanceID, createdBinding.InstanceID)
		assert.NotNil(t, createdBinding.ExpirationSeconds)
		assert.Equal(t, fixedBinding.ExpirationSeconds, createdBinding.ExpirationSeconds)
		assert.NotNil(t, createdBinding.Kubeconfig)
		assert.Equal(t, fixedBinding.Kubeconfig, createdBinding.Kubeconfig)

		// when
		err = brokerStorage.Bindings().Delete(testBindingId)

		// then
		nonExisting, err := brokerStorage.Bindings().Get(testBindingId)
		assert.Error(t, err)
		assert.Nil(t, nonExisting)
	})
}
