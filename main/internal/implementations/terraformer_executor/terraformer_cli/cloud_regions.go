package terraformercli

import (
	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func getValidRegions(cloudRegions []terraformValueObjects.CloudRegion, providerRegions map[string]bool, defaultRegions []string) []string {
	if len(cloudRegions) == 0 {
		return defaultRegions
	}

	regions := make([]string, 0)
	for _, region := range cloudRegions {
		if providerRegions[string(region)] {
			regions = append(regions, string(region))
			break
		}
	}

	if len(regions) == 0 {
		return defaultRegions
	}

	logrus.Debugf("[terraformer_executor][getValidRegions] Valid regions: %v", regions)
	return regions
}
