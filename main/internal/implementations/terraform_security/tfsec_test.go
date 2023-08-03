package terraformSecurity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
)

func TestGetResultsWithResourceID(t *testing.T) {
	// Given
	tfsec := TFSec{}

	resources := make(map[driftDetector.ResourceIdentifier]string)
	resources["resourceType1.resourceName1"] = "uniqueID1"
	resources["resourceType2.resourceName2"] = "uniqueID2"
	resources["resourceType3.resourceName3"] = "uniqueID3"

	results := TFSecFile{
		[]Result{
			{Resource: "resourceType1.resourceName1"},
			{Resource: "resourceType2.resourceName2"},
			{Resource: "resourceType3.resourceName3"},
		},
	}
	// When
	newResults := tfsec.getResultsWithResourceID(results, resources)

	// Then
	expectedResults := TFSecFile{
		[]Result{
			{ID: "uniqueID1", Resource: "resourceType1.resourceName1"},
			{ID: "uniqueID2", Resource: "resourceType2.resourceName2"},
			{ID: "uniqueID3", Resource: "resourceType3.resourceName3"},
		},
	}
	assert.Equal(t, expectedResults, newResults, "The expected and actual results should match")
}

func TestMapResourceIDsFromStateFile(t *testing.T) {
	// Given
	tfsec := TFSec{}

	file := driftDetector.TerraformerStateFile{
		Resources: []*driftDetector.TerraformerResource{
			{
				Type: "type1",
				Name: "name1",
				Instances: []driftDetector.TerraformerInstance{
					{
						AttributesFlat: map[string]string{
							"id": "id1",
						},
					},
				},
			},
			{
				Type: "type2",
				Name: "name2",
				Instances: []driftDetector.TerraformerInstance{
					{
						AttributesFlat: map[string]string{
							"id": "id2",
						},
					},
				},
			},
		},
	}

	expectedResourcesMap := map[driftDetector.ResourceIdentifier]string{
		"type1.name1": "id1",
		"type2.name2": "id2",
	}

	// When
	resourcesMap := tfsec.mapResourceIDsFromStateFile(file)

	// Then
	assert.Equal(t, expectedResourcesMap, resourcesMap, "The expected and actual resource maps should match")
}
