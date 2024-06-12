package broker

import (
	"testing"

	"github.com/kyma-project/kyma-environment-broker/internal/broker/automock"
	"github.com/stretchr/testify/assert"
)

func TestOneForAllConvergedCloudRegionsProvider_GetDefaultRegions(t *testing.T) {
	// given
	c := &OneForAllConvergedCloudRegionsProvider{}

	// when
	result := c.GetRegions()

	// then
	assert.Equal(t, []string{"eu-de-1"}, result)
}

func TestPathBasedConvergedCloudRegionsProvider_FactoryMethod(t *testing.T) {
	// given
	configLocation := "path-to-config"
	regions := map[string][]string{
		"eu-de-1": {"eu-de-1"},
	}

	mockReader := automock.NewRegionReader(t)
	mockReader.On("Read", configLocation).Return(regions, nil)
	
	// when
	provider, err := NewPathBasedConvergedCloudRegionsProvider(configLocation, mockReader)

	// then
	assert.NoError(t, err)
	assert.Equal(t, regions, provider.regionConfiguration)
}
