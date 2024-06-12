package broker

import (
	"fmt"
)

//go:generate mockery --name=RegionReader --output=automock --outpkg=automock --case=underscore
type RegionReader interface {
	Read(filename string) (map[string][]string, error)
}

type ConvergedCloudRegionProvider interface {
	GetRegions() []string
}

type PathBasedConvergedCloudRegionsProvider struct {
	// placeholder
	regionConfiguration map[string][]string
}

func NewPathBasedConvergedCloudRegionsProvider(regionConfigurationPath string, reader RegionReader) (*PathBasedConvergedCloudRegionsProvider, error) {
	regionConfiguration, err := reader.Read(regionConfigurationPath)
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling a file with sap-converged-cloud region mappings: %w", err)
	}

	return &PathBasedConvergedCloudRegionsProvider{
		regionConfiguration: regionConfiguration,
	}, nil
}

func (c *PathBasedConvergedCloudRegionsProvider) GetRegions(region string) []string {
	item, found:=c.regionConfiguration[region]

	if  !found {
		return []string{}
	}
	
	return item
}

type OneForAllConvergedCloudRegionsProvider struct {
}

func (c *OneForAllConvergedCloudRegionsProvider) GetRegions() []string {
	return []string{"eu-de-1"}
}
