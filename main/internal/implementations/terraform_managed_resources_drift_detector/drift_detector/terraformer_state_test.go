package driftDetector

import (
	"reflect"
	"testing"
)

func Test_ExtractUniqueResourceIdToData(t *testing.T) {
	// Given
	inputStateFile := TerraformerStateFile{
		Resources: []*TerraformerResource{
			{
				Mode:     "managed",
				Module:   "root",
				Type:     "google_example",
				Name:     "my_resource",
				Provider: "google",
				Instances: []TerraformerInstance{
					{
						SchemaVersion: 1,
						AttributesFlat: map[string]string{
							"id": "id_1",
						},
					},
					{
						SchemaVersion: 1,
						AttributesFlat: map[string]string{
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
				Instances: []TerraformerInstance{
					{
						SchemaVersion: 1,
						AttributesFlat: map[string]string{
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
	output := m.extractUniqueResourceIDToData(inputStateFile)

	expectedOutput := TerraformerResourceIDToData{
		"google_example.id_1": TerraformerUniqueResourceData{
			Module:         "root",
			Type:           "google_example",
			Name:           "my_resource",
			Provider:       "google",
			AttributesFlat: map[string]string{"id": "id_1"},
		},
		"google_example.id_2": TerraformerUniqueResourceData{
			Module:         "root",
			Type:           "google_example",
			Name:           "my_resource",
			Provider:       "google",
			AttributesFlat: map[string]string{"id": "id_2"},
		},
		"aws_example.id_3": TerraformerUniqueResourceData{
			Module:         "my_module",
			Type:           "aws_example",
			Name:           "my_resource",
			Provider:       "aws",
			AttributesFlat: map[string]string{"id": "id_3", "xyz": "abc"},
		},
	}

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}
