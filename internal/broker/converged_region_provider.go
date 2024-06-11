package broker

type ConvergedCloudRegionProvider interface {
	GetRegions() []string
}

type PathBasedCCEERegionProvider struct {
	// placeholder
}

func (c *PathBasedCCEERegionProvider) GetRegions() []string {
	return []string{"eu-de-1"}
}

type OneForAllConvergedCloudRegionsProvider struct {
}

func (c *OneForAllConvergedCloudRegionsProvider) GetRegions() []string {
	return []string{"eu-de-1"}
}
