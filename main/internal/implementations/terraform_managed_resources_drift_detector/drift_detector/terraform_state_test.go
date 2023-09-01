package driftdetector

import (
	"reflect"
	"testing"
)

func Test_TerraformStateExtractUniqueResourceIdToData(t *testing.T) {
	// Given
	inputStateFileName := "My State File"

	inputStateFile := TerraformStateFile{
		Resources: []*Resource{
			{
				Mode:     "managed",
				Module:   "root",
				Type:     "google_example",
				Name:     "my_resource",
				Provider: "google",
				Instances: []ResourceInstance{
					ResourceInstance{
						SchemaVersion: 1,
						Attributes: map[string]interface{}{
							"id": "id_1",
						},
					},
					ResourceInstance{
						SchemaVersion: 1,
						Attributes: map[string]interface{}{
							"id": "id_2",
						},
					},
				},
			},
			{
				Mode:     "managed",
				Module:   "my_module",
				Type:     "aws_example",
				Name:     "my_resource",
				Provider: "aws",
				Instances: []ResourceInstance{
					ResourceInstance{
						SchemaVersion: 1,
						Attributes: map[string]interface{}{
							"id":  "id_3",
							"xyz": "abc",
						},
					},
				},
			},
		},
	}

	m := ManagedResourcesDriftDetector{}

	// When
	output := m.terraformStateExtractUniqueResourceIDToData(inputStateFileName, inputStateFile)

	expectedOutput := TerraformStateResourceIDToData{
		"google_example.id_1": TerraformStateUniqueResourceData{
			StateFile:  "My State File",
			Module:     "root",
			Type:       "google_example",
			Name:       "my_resource",
			Provider:   "google",
			Attributes: map[string]interface{}{"id": "id_1"},
		},
		"google_example.id_2": TerraformStateUniqueResourceData{
			StateFile:  "My State File",
			Module:     "root",
			Type:       "google_example",
			Name:       "my_resource",
			Provider:   "google",
			Attributes: map[string]interface{}{"id": "id_2"},
		},
		"aws_example.id_3": TerraformStateUniqueResourceData{
			StateFile:  "My State File",
			Module:     "my_module",
			Type:       "aws_example",
			Name:       "my_resource",
			Provider:   "aws",
			Attributes: map[string]interface{}{"id": "id_3", "xyz": "abc"},
		},
	}

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}
