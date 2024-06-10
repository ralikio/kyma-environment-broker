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

type OneForAllCCEERegionProvider struct {
}

func (c *OneForAllCCEERegionProvider) GetRegions() []string {
	return []string{"eu-de-1"}
}
