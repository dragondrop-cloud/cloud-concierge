package terraformerCLI

import (
	"testing"

	"github.com/stretchr/testify/require"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func Test_getValidRegions_notFoundRegions(t *testing.T) {
	// Given
	cloudRegions := []terraformValueObjects.CloudRegion{
		terraformValueObjects.CloudRegion("westus2"),
		terraformValueObjects.CloudRegion("eastus"),
	}
	providerRegions := map[string]bool{
		"us-east-1": true,
		"us-east-2": true,
	}
	defaultRegions := []string{"us-east-1"}

	// When
	regions := getValidRegions(cloudRegions, providerRegions, defaultRegions)

	// Then
	require.Equal(t, []string{"us-east-1"}, regions)
}

func Test_getValidRegions_providerRegionsLimitedToSingleRegion(t *testing.T) {
	// Given
	cloudRegions := []terraformValueObjects.CloudRegion{
		terraformValueObjects.CloudRegion("us-east-1"),
		terraformValueObjects.CloudRegion("us-east-2"),
	}
	providerRegions := map[string]bool{
		"us-east-1": true,
		"us-east-2": true,
	}
	defaultRegions := []string{"us-east-1"}

	// When
	regions := getValidRegions(cloudRegions, providerRegions, defaultRegions)

	// Then
	require.Equal(t, []string{"us-east-1"}, regions)
}

func Test_getValidRegions_emptyCloudRegions(t *testing.T) {
	// Given
	cloudRegions := make([]terraformValueObjects.CloudRegion, 0)
	providerRegions := map[string]bool{
		"us-east-1": true,
		"us-east-2": true,
	}
	defaultRegions := []string{"us-east-1"}

	// When
	regions := getValidRegions(cloudRegions, providerRegions, defaultRegions)

	// Then
	require.Equal(t, []string{"us-east-1"}, regions)
}
