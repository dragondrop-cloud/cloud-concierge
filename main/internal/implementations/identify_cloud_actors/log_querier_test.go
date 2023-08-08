package identifyCloudActors

import (
	"reflect"
	"testing"

	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
)

func TestCreateDivisionUniqueDriftedResources(t *testing.T) {
	// Given
	inputDifferences := []driftDetector.AttributeDifference{
		{
			InstanceID:    "dragondrop",
			AttributeName: "1",
			AttributeDetail: driftDetector.AttributeDetail{
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
		{
			InstanceID:    "dragondrop",
			AttributeName: "2",
			AttributeDetail: driftDetector.AttributeDetail{
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
		{
			InstanceID:    "dragondrop",
			AttributeName: "5",
			AttributeDetail: driftDetector.AttributeDetail{
				ResourceType:  "type-3",
				ResourceName:  "name-3",
				StateFileName: "state-3",
			},
		},
	}

	expectedOutput := UniqueDriftedResources{
		"state-1.type-1.name-1.dragondrop": {
			ResourceType:  "type-1",
			ResourceName:  "name-1",
			StateFileName: "state-1",
			InstanceID:    "dragondrop",
		},
		"state-3.type-3.name-3.dragondrop": {

			ResourceType:  "type-3",
			ResourceName:  "name-3",
			StateFileName: "state-3",
			InstanceID:    "dragondrop",
		},
	}

	// When
	output, err := createUniqueDriftedResources(inputDifferences)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then
	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v\n", output, expectedOutput)
	}
}

func TestDetermineActionClass(t *testing.T) {
	output := determineActionClass("cloud.run.service.update")
	if output != "modification" {
		t.Errorf("got:\n%vexpected:\n%v", output, "modification")
	}

	output = determineActionClass("cloud.run.service.replace")
	if output != "modification" {
		t.Errorf("got:\n%vexpected:\n%v", output, "modification")
	}

	output = determineActionClass("cloud.run.service.delete")
	if output != "deletion" {
		t.Errorf("got:\n%vexpected:\n%v", output, "deletion")
	}

	output = determineActionClass("cloud.run.service.create")
	if output != "creation" {
		t.Errorf("got:\n%vexpected:\n%v", output, "creation")
	}

	output = determineActionClass("cloud.run.service")
	if output != "not_classified" {
		t.Errorf("got:\n%vexpected:\n%v", output, "not_classified")
	}
}

func TestUniqueDriftedResourceToName(t *testing.T) {
	// Given
	input := UniqueDriftedResource{
		StateFileName: "state-1",
		ResourceType:  "type-1",
		ResourceName:  "name-1",
		InstanceID:    "dragondrop",
	}
	expectedOutput := "state-1.type-1.name-1.dragondrop"

	// When
	output := uniqueDriftedResourceToName(input)

	// Then
	if string(output) != expectedOutput {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestAttributeDriftedResourceToName(t *testing.T) {
	// Given
	input := driftDetector.AttributeDifference{
		AttributeDetail: driftDetector.AttributeDetail{
			StateFileName: "state-1",
			ResourceType:  "type-1",
			ResourceName:  "name-1",
		},
		InstanceID: "dragondrop",
	}
	expectedOutput := "state-1.type-1.name-1.dragondrop"

	// When
	output := attributeDifferenceToResourceName(input)

	// Then
	if string(output) != expectedOutput {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}
