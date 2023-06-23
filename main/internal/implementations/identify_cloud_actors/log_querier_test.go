package identifyCloudActors

import (
	"reflect"
	"testing"

	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestCreateDivisionUniqueDriftedResources(t *testing.T) {
	// Given
	inputDifferences := []driftDetector.AttributeDifference{
		{
			InstanceID:    "dragondrop",
			AttributeName: "1",
			AttributeDetail: driftDetector.AttributeDetail{
				CloudDivision: "div-1",
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
		{
			InstanceID:    "dragondrop",
			AttributeName: "2",
			AttributeDetail: driftDetector.AttributeDetail{
				CloudDivision: "div-1",
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
		{
			InstanceID:    "dragondrop",
			AttributeName: "5",
			AttributeDetail: driftDetector.AttributeDetail{
				CloudDivision: "div-1",
				ResourceType:  "type-3",
				ResourceName:  "name-3",
				StateFileName: "state-3",
			},
		},
		{
			InstanceID:    "dragondrop",
			AttributeName: "4",
			AttributeDetail: driftDetector.AttributeDetail{
				CloudDivision: "div-2",
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
			},
		},
	}

	expectedOutput := DivisionToUniqueDriftedResources{
		"div-1": {
			"state-1.type-1.name-1.dragondrop": {
				CloudDivision: "div-1",
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
				InstanceID:    "dragondrop",
			},
			"state-3.type-3.name-3.dragondrop": {
				CloudDivision: "div-1",
				ResourceType:  "type-3",
				ResourceName:  "name-3",
				StateFileName: "state-3",
				InstanceID:    "dragondrop",
			},
		},
		"div-2": {
			"state-1.type-1.name-1.dragondrop": {
				CloudDivision: "div-2",
				ResourceType:  "type-1",
				ResourceName:  "name-1",
				StateFileName: "state-1",
				InstanceID:    "dragondrop",
			},
		},
	}

	// When
	output, err := createDivisionUniqueDriftedResources(inputDifferences)
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

func TestFilterDivisionCloudCredentialsForProvider(t *testing.T) {
	// Testing base filtering case
	divisionToProvider := map[terraformValueObjects.Division]terraformValueObjects.Provider{
		"div-1": "google",
		"div-2": "google",
		"div-3": "aws",
	}

	inputGlobalConfig := Config{
		DivisionCloudCredentials: map[terraformValueObjects.Division]terraformValueObjects.Credential{
			"div-1": "asd",
			"div-2": "qwe",
			"div-3": "plo",
		},
	}

	expectedOutput := terraformValueObjects.DivisionCloudCredentialDecoder{
		"div-1": "asd",
		"div-2": "qwe",
	}

	output := filterDivisionCloudCredentialsForProvider("google", divisionToProvider, inputGlobalConfig)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}

	// Testing empty case once filtered
	expectedOutput = terraformValueObjects.DivisionCloudCredentialDecoder{}

	output = filterDivisionCloudCredentialsForProvider("azurerm", divisionToProvider, inputGlobalConfig)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
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
