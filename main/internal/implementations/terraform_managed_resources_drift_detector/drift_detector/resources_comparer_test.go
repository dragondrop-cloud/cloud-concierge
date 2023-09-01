package driftdetector

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExistentResourcesHaveChanged(t *testing.T) {
	detector := &ManagedResourcesDriftDetector{}

	// Given
	terraformerResourcesIDToData := TerraformerResourceIDToData{
		"google_example.id_1": TerraformerUniqueResourceData{
			Module:         "root",
			Type:           "google_example",
			Name:           "my_resource",
			Provider:       "google",
			AttributesFlat: map[string]string{"id": "id_1"},
		},
	}

	stateFileResourcesIDToData := TerraformStateResourceIDToData{
		"google_example.id_1": TerraformStateUniqueResourceData{
			StateFile:  "My State File",
			Module:     "root",
			Type:       "google_example",
			Name:       "my_resource",
			Provider:   "google",
			Attributes: map[string]interface{}{"id": "id_1"},
		},
	}

	// When
	differences, err := detector.identifyResourceDifferences(terraformerResourcesIDToData, stateFileResourcesIDToData)

	// Then
	require.Len(t, differences, 0)
	require.NoError(t, err)

	// Modify the remote resource instance attribute to simulate a change
	stateFileResourcesIDToData["google_example.id_1"].Attributes["id"] = "modified-dragondrop-modules"

	// When
	differences, err = detector.identifyResourceDifferences(terraformerResourcesIDToData, stateFileResourcesIDToData)

	// Then
	require.NoError(t, err)
	require.Len(t, differences, 1)
	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "id",
		TerraformValue: "modified-dragondrop-modules",
		CloudValue:     "id_1",
		InstanceID:     "id_1",
		AttributeDetail: AttributeDetail{
			StateFileName: "My State File",
			ModuleName:    "root",
			ResourceType:  "google_example",
			ResourceName:  "my_resource",
		},
	},
	)
}

func TestGetExistentResourcesHaveChanged_MoreThanOneAttribute(t *testing.T) {
	detector := &ManagedResourcesDriftDetector{}

	// Given
	terraformerResourcesIDToData := TerraformerResourceIDToData{
		"id_1": TerraformerUniqueResourceData{
			Module:         "root",
			Type:           "google_example",
			Name:           "my_resource",
			Provider:       "google",
			AttributesFlat: map[string]string{"id": "id_1", "abc": "123"},
		},
	}

	stateFileResourcesIDToData := TerraformStateResourceIDToData{
		"id_1": TerraformStateUniqueResourceData{
			StateFile:  "My State File",
			Module:     "root",
			Type:       "google_example",
			Name:       "my_resource",
			Provider:   "google",
			Attributes: map[string]interface{}{"id": "id_2", "abc": "456"},
		},
	}

	// When
	differences, err := detector.identifyResourceDifferences(terraformerResourcesIDToData, stateFileResourcesIDToData)

	// Then
	require.NoError(t, err)
	require.Len(t, differences, 2)
	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "id",
		TerraformValue: "id_2",
		CloudValue:     "id_1",
		InstanceID:     "id_1",
		AttributeDetail: AttributeDetail{
			StateFileName: "My State File",
			ModuleName:    "root",
			ResourceType:  "google_example",
			ResourceName:  "my_resource",
		},
	},
	)

	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "abc",
		TerraformValue: "456",
		CloudValue:     "123",
		InstanceID:     "id_1",
		AttributeDetail: AttributeDetail{
			StateFileName: "My State File",
			ModuleName:    "root",
			ResourceType:  "google_example",
			ResourceName:  "my_resource",
		},
	},
	)
}

func TestConvertNestedMapToFlatAttributes(t *testing.T) {
	// Given
	input := map[string]interface{}{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule": []interface{}{
			map[string]interface{}{
				"action": []interface{}{
					map[string]interface{}{
						"storage_class": "",
						"type":          "Delete",
					},
				},
				"condition": []interface{}{
					map[string]interface{}{
						"age":                        "0",
						"created_before":             "",
						"custom_time_before":         "",
						"days_since_custom_time":     "0",
						"days_since_noncurrent_time": "0",
						"noncurrent_time_before":     "",
						"num_newer_versions":         "3",
						"with_state":                 "ARCHIVED",
					},
				},
			},
			map[string]interface{}{
				"action": []interface{}{
					map[string]interface{}{
						"storage_class": "",
						"type":          "Delete",
					},
				},
				"condition": []interface{}{
					map[string]interface{}{
						"age":                        "0",
						"created_before":             "",
						"custom_time_before":         "",
						"days_since_custom_time":     "0",
						"days_since_noncurrent_time": "7",
						"noncurrent_time_before":     "",
						"num_newer_versions":         "0",
						"with_state":                 "ANY",
					},
				},
			},
		},
		"location":                    "US",
		"name":                        "dragondrop-modules",
		"project":                     "dragondrop-dev",
		"requester_pays":              "false",
		"self_link":                   "https://www.googleapis.com/storage/v1/b/dragondrop-modules",
		"storage_class":               "STANDARD",
		"uniform_bucket_level_access": "true",
		"url":                         "gs://dragondrop-modules",
		"versioning": []interface{}{
			map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	expectedOutput := map[string]string{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule.0.action.0.storage_class":                 "",
		"lifecycle_rule.0.action.0.type":                          "Delete",
		"lifecycle_rule.0.condition.0.age":                        "0",
		"lifecycle_rule.0.condition.0.created_before":             "",
		"lifecycle_rule.0.condition.0.custom_time_before":         "",
		"lifecycle_rule.0.condition.0.days_since_custom_time":     "0",
		"lifecycle_rule.0.condition.0.days_since_noncurrent_time": "0",
		"lifecycle_rule.0.condition.0.noncurrent_time_before":     "",
		"lifecycle_rule.0.condition.0.num_newer_versions":         "3",
		"lifecycle_rule.0.condition.0.with_state":                 "ARCHIVED",
		"lifecycle_rule.1.action.0.storage_class":                 "",
		"lifecycle_rule.1.action.0.type":                          "Delete",
		"lifecycle_rule.1.condition.0.age":                        "0",
		"lifecycle_rule.1.condition.0.created_before":             "",
		"lifecycle_rule.1.condition.0.custom_time_before":         "",
		"lifecycle_rule.1.condition.0.days_since_custom_time":     "0",
		"lifecycle_rule.1.condition.0.days_since_noncurrent_time": "7",
		"lifecycle_rule.1.condition.0.noncurrent_time_before":     "",
		"lifecycle_rule.1.condition.0.num_newer_versions":         "0",
		"lifecycle_rule.1.condition.0.with_state":                 "ANY",
		"location":                    "US",
		"name":                        "dragondrop-modules",
		"project":                     "dragondrop-dev",
		"requester_pays":              "false",
		"self_link":                   "https://www.googleapis.com/storage/v1/b/dragondrop-modules",
		"storage_class":               "STANDARD",
		"uniform_bucket_level_access": "true",
		"url":                         "gs://dragondrop-modules",
		"versioning.0.enabled":        "true",
	}

	// When
	output, err := convertNestedMapToFlatAttributes(input)

	// Then
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

}

func TestConvertNestedMapToFlatAttributes_Two(t *testing.T) {
	m := ManagedResourcesDriftDetector{}
	inputTerraformStateFile, err := m.parseRemoteStateFile([]byte(`{
  "version": 4,
  "terraform_version": "1.5.0",
  "serial": 24,
  "lineage": "dfff0d5f-0f62-bb39-6f7d-0cabe520755c",
  "outputs": {
    "container_path_api": {
      "value": "us-east4-docker.pkg.dev/dragondrop-dev/artifact-repository-dev/api",
      "type": "string"
    }
  },
  "resources": [
    {
      "mode": "managed",
      "type": "google_storage_bucket",
      "name": "dragondrop_modules",
      "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "autoclass": [],
            "cors": [],
            "custom_placement_config": [],
            "default_event_based_hold": false,
            "encryption": [],
            "force_destroy": false,
            "id": "dragondrop-modules",
            "labels": {},
            "lifecycle_rule": [
              {
                "action": [
                  {
                    "storage_class": "",
                    "type": "Delete"
                  }
                ],
                "condition": [
                  {
                    "age": 0,
                    "created_before": "",
                    "custom_time_before": "",
                    "days_since_custom_time": 0,
                    "days_since_noncurrent_time": 0,
                    "matches_prefix": [],
                    "matches_storage_class": [],
                    "matches_suffix": [],
                    "noncurrent_time_before": "",
                    "num_newer_versions": 3,
                    "with_state": "ARCHIVED"
                  }
                ]
              },
              {
                "action": [
                  {
                    "storage_class": "",
                    "type": "Delete"
                  }
                ],
                "condition": [
                  {
                    "age": 0,
                    "created_before": "",
                    "custom_time_before": "",
                    "days_since_custom_time": 0,
                    "days_since_noncurrent_time": 7,
                    "matches_prefix": [],
                    "matches_storage_class": [],
                    "matches_suffix": [],
                    "noncurrent_time_before": "",
                    "num_newer_versions": 0,
                    "with_state": "ANY"
                  }
                ]
              }
            ],
            "location": "US",
            "logging": [],
            "name": "dragondrop-modules",
            "project": "dragondrop-dev",
            "public_access_prevention": "enforced",
            "requester_pays": false,
            "retention_policy": [],
            "self_link": "https://www.googleapis.com/storage/v1/b/dragondrop-modules",
            "storage_class": "STANDARD",
            "timeouts": null,
            "uniform_bucket_level_access": true,
            "url": "gs://dragondrop-modules",
            "versioning": [
              {
                "enabled": true
              }
            ],
            "website": []
          },
          "sensitive_attributes": [],
          "private": "eyJlMmJmYjczMC1lY2FhLTExZTYtOGY4OC0zNDM2M2JjN2M0YzAiOnsiY3JlYXRlIjoyNDAwMDAwMDAwMDAsInJlYWQiOjI0MDAwMDAwMDAwMCwidXBkYXRlIjoyNDAwMDAwMDAwMDB9fQ=="
        }
      ]
    }
  ],
  "check_results": null
}`))

	if err != nil {
		t.Errorf("unexpected error parsing remote state file: %s", err)
	}

	expectedOutput := map[string]string{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule.0.action.0.storage_class":                 "",
		"lifecycle_rule.0.action.0.type":                          "Delete",
		"lifecycle_rule.0.condition.0.age":                        "0",
		"lifecycle_rule.0.condition.0.created_before":             "",
		"lifecycle_rule.0.condition.0.custom_time_before":         "",
		"lifecycle_rule.0.condition.0.days_since_custom_time":     "0",
		"lifecycle_rule.0.condition.0.days_since_noncurrent_time": "0",
		"lifecycle_rule.0.condition.0.noncurrent_time_before":     "",
		"lifecycle_rule.0.condition.0.num_newer_versions":         "3",
		"lifecycle_rule.0.condition.0.with_state":                 "ARCHIVED",
		"lifecycle_rule.1.action.0.storage_class":                 "",
		"lifecycle_rule.1.action.0.type":                          "Delete",
		"lifecycle_rule.1.condition.0.age":                        "0",
		"lifecycle_rule.1.condition.0.created_before":             "",
		"lifecycle_rule.1.condition.0.custom_time_before":         "",
		"lifecycle_rule.1.condition.0.days_since_custom_time":     "0",
		"lifecycle_rule.1.condition.0.days_since_noncurrent_time": "7",
		"lifecycle_rule.1.condition.0.noncurrent_time_before":     "",
		"lifecycle_rule.1.condition.0.num_newer_versions":         "0",
		"lifecycle_rule.1.condition.0.with_state":                 "ANY",
		"location":                    "US",
		"name":                        "dragondrop-modules",
		"project":                     "dragondrop-dev",
		"public_access_prevention":    "enforced",
		"requester_pays":              "false",
		"self_link":                   "https://www.googleapis.com/storage/v1/b/dragondrop-modules",
		"storage_class":               "STANDARD",
		"uniform_bucket_level_access": "true",
		"url":                         "gs://dragondrop-modules",
		"versioning.0.enabled":        "true",
	}

	// When
	resources := TerraformStateResourceIDToData{}

	resourcesFromStateFile := m.terraformStateExtractUniqueResourceIDToData("", inputTerraformStateFile)

	for resourceID, resourceData := range resourcesFromStateFile {
		resources[resourceID] = resourceData
	}

	// Then
	for _, data := range resources {
		flatAttributes, err := convertNestedMapToFlatAttributes(data.Attributes)
		if err != nil {
			t.Errorf("unexpected error converting nested map to flat attributes: %s", err)
		}
		if !reflect.DeepEqual(flatAttributes, expectedOutput) {
			formattedExpectation, _ := json.MarshalIndent(expectedOutput, "", "  ")
			formattedFlatAttributes, _ := json.MarshalIndent(flatAttributes, "", "  ")

			t.Errorf("expected:\n%v\ngot:\n%v", string(formattedExpectation), string(formattedFlatAttributes))
		}
		break
	}
}

func TestStringMapsEqual(t *testing.T) {
	// Given
	testCases := []struct {
		m1       map[string]string
		m2       map[string]string
		expected bool
	}{
		{
			m1:       map[string]string{"a": "1", "b": "2"},
			m2:       map[string]string{"a": "1", "b": "2"},
			expected: true,
		},
		{
			m1:       map[string]string{"a": "1", "b": "2"},
			m2:       map[string]string{"a": "1", "b": "3"},
			expected: false,
		},
		{
			m1:       map[string]string{"a": "1", "b": "2"},
			m2:       map[string]string{"a": "1"},
			expected: false,
		},
		{
			m1:       map[string]string{},
			m2:       map[string]string{},
			expected: true,
		},
	}

	for i, tc := range testCases {
		// When
		result := stringMapsEqual(tc.m1, tc.m2)

		// Then
		assert.Equal(t, tc.expected, result, "Test case %d failed", i+1)
	}
}

func TestStringSlicesEqual(t *testing.T) {
	// Given
	testCases := []struct {
		s1       []string
		s2       []string
		expected bool
	}{
		{
			s1:       []string{"a", "b", "c"},
			s2:       []string{"a", "b", "c"},
			expected: true,
		},
		{
			s1:       []string{"a", "b", "c"},
			s2:       []string{"a", "c", "b"},
			expected: true,
		},
		{
			s1:       []string{"a", "b", "c"},
			s2:       []string{"a", "b", "d"},
			expected: false,
		},
		{
			s1:       []string{"a", "b", "c"},
			s2:       []string{"a", "b"},
			expected: false,
		},
		{
			s1:       []string{},
			s2:       []string{},
			expected: true,
		},
	}

	for i, tc := range testCases {
		// When
		result := stringSlicesEqual(tc.s1, tc.s2)

		// Then
		assert.Equal(t, tc.expected, result, "Test case %d failed", i+1)
	}
}

func Test_compareAttributesAndGetDrifted(t *testing.T) {
	// Given
	terraformerResourceAttributes := map[string]string{
		"cors.#":                   "0",
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"labels.%":                 "0",
		"lifecycle_rule.#":         "1",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "Delete",
		"lifecycle_rule.0.condition.0.age":        "0",
	}

	remoteResourceAttributes := map[string]string{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "Delete",
		"lifecycle_rule.0.condition.0.age":        "0",
	}

	attributeComplement := &AttributeDetail{
		StateFileName: StateFileName("state_file_name"),
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "my_storage_bucket",
	}

	// When
	differences, resourcesChanged, err := compareFlatAttributesAndGetDrifted(remoteResourceAttributes, terraformerResourceAttributes, attributeComplement)
	if err != nil {
		t.Errorf("Error should be nil, got: %v", err)
	}

	// Then
	require.False(t, resourcesChanged)
	require.Len(t, differences, 0)
}

func Test_compareAttributesAndGetDrifted_OneAttribute(t *testing.T) {
	// Given
	terraformerResourceAttributes := map[string]string{
		"cors.#":                   "0",
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"labels.%":                 "0",
		"lifecycle_rule.#":         "1",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "NewValue",
		"lifecycle_rule.0.condition.0.age":        "0",
	}

	remoteResourceAttributes := map[string]string{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "OldValue",
		"lifecycle_rule.0.condition.0.age":        "0",
	}

	attributeComplement := &AttributeDetail{
		StateFileName: StateFileName("state_file_name"),
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "my_storage_bucket",
	}

	// When
	differences, resourcesChanged, err := compareFlatAttributesAndGetDrifted(remoteResourceAttributes, terraformerResourceAttributes, attributeComplement)
	if err != nil {
		t.Errorf("Error should be nil, got: %v", err)
	}

	// Then
	require.True(t, resourcesChanged)
	require.Len(t, differences, 1)
	require.Equal(t, AttributeDifference{
		AttributeName:  "lifecycle_rule.0.action.0.type",
		TerraformValue: "OldValue",
		CloudValue:     "NewValue",
		InstanceID:     "projects/_/buckets/dragondrop-modules",
		AttributeDetail: AttributeDetail{
			StateFileName: "state_file_name",
			ModuleName:    "module_name",
			ResourceType:  "google_storage_bucket",
			ResourceName:  "my_storage_bucket",
		},
	}, differences[0])
}

func Test_compareAttributesAndGetDrifted_MoreThanOneAttributes(t *testing.T) {
	// Given
	terraformerResourceAttributes := map[string]string{
		"cors.#":                   "0",
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"labels.%":                 "0",
		"lifecycle_rule.#":         "1",
		"lifecycle_rule.0.action.0.storage_class": "NewValue",
		"lifecycle_rule.0.action.0.type":          "NewValue",
		"lifecycle_rule.0.condition.0.age":        "NewValue",
	}

	remoteResourceAttributes := map[string]string{
		"default_event_based_hold":         "false",
		"force_destroy":                    "false",
		"id":                               "dragondrop-modules",
		"lifecycle_rule.0.action.0.type":   "OldValue",
		"lifecycle_rule.0.condition.0.age": "0",
	}

	attributeComplement := &AttributeDetail{
		StateFileName: StateFileName("state_file_name"),
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "my_storage_bucket",
	}

	// When
	differences, resourcesChanged, err := compareFlatAttributesAndGetDrifted(remoteResourceAttributes, terraformerResourceAttributes, attributeComplement)
	if err != nil {
		t.Errorf("Error should be nil, got: %v", err)
	}

	// Then
	require.True(t, resourcesChanged)
	require.Len(t, differences, 3)
	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "lifecycle_rule.0.action.0.storage_class",
		TerraformValue: "",
		CloudValue:     "NewValue",
		InstanceID:     "projects/_/buckets/dragondrop-modules",
		AttributeDetail: AttributeDetail{
			StateFileName: "state_file_name",
			ModuleName:    "module_name",
			ResourceType:  "google_storage_bucket",
			ResourceName:  "my_storage_bucket",
		},
	})
	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "lifecycle_rule.0.action.0.type",
		TerraformValue: "OldValue",
		CloudValue:     "NewValue",
		InstanceID:     "projects/_/buckets/dragondrop-modules",
		AttributeDetail: AttributeDetail{
			StateFileName: "state_file_name",
			ModuleName:    "module_name",
			ResourceType:  "google_storage_bucket",
			ResourceName:  "my_storage_bucket",
		},
	})
	require.Contains(t, differences, AttributeDifference{
		AttributeName:  "lifecycle_rule.0.condition.0.age",
		TerraformValue: "0",
		CloudValue:     "NewValue",
		InstanceID:     "projects/_/buckets/dragondrop-modules",
		AttributeDetail: AttributeDetail{
			StateFileName: "state_file_name",
			ModuleName:    "module_name",
			ResourceType:  "google_storage_bucket",
			ResourceName:  "my_storage_bucket",
		},
	})
}

func Test_compareAttributesAndGetDrifted_DeleteAttribute(t *testing.T) {
	// Given
	terraformerResourceAttributes := map[string]string{
		"cors.#":                   "0",
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"labels.%":                 "0",
		"lifecycle_rule.#":         "1",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "1",
		"lifecycle_rule.0.condition.0.age":        "0",
	}

	remoteResourceAttributes := map[string]string{
		"default_event_based_hold": "false",
		"force_destroy":            "false",
		"id":                       "dragondrop-modules",
		"lifecycle_rule.0.action.0.storage_class": "",
		"lifecycle_rule.0.action.0.type":          "1",
		"lifecycle_rule.0.condition.0.age":        "0",
		"attribute.deleted":                       "deleted_value",
	}

	attributeComplement := &AttributeDetail{
		StateFileName: StateFileName("state_file_name"),
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "my_storage_bucket",
	}

	// When
	differences, resourcesChanged, err := compareFlatAttributesAndGetDrifted(remoteResourceAttributes, terraformerResourceAttributes, attributeComplement)
	if err != nil {
		t.Errorf("Error should be nil, got: %v", err)
	}

	// Then
	require.True(t, resourcesChanged)
	require.Len(t, differences, 1)
	require.Equal(t, AttributeDifference{
		AttributeName:  "attribute.deleted",
		TerraformValue: "deleted_value",
		CloudValue:     "",
		InstanceID:     "projects/_/buckets/dragondrop-modules",
		AttributeDetail: AttributeDetail{
			StateFileName: "state_file_name",
			ModuleName:    "module_name",
			ResourceType:  "google_storage_bucket",
			ResourceName:  "my_storage_bucket",
		},
	}, differences[0])
}

func TestConvertNestedMapToFlatAttributes_EdgeCasesWithOtherTypes(t *testing.T) {
	// Given
	nestedAttributes := map[string]interface{}{
		"autoclass":                []interface{}{},
		"cors":                     []interface{}{},
		"custom_placement_config":  []interface{}{},
		"default_event_based_hold": false,
		"encryption":               []interface{}{},
		"force_destroy":            false,
		"id":                       "dragondrop-modules",
		"labels":                   map[string]interface{}{"": ""},
		"lifecycle_rule": []interface{}{
			map[string]interface{}{
				"action": []interface{}{
					map[string]interface{}{
						"storage_class": "",
						"type":          "Delete",
					},
				},
				"condition": []interface{}{
					map[string]interface{}{
						"age":                        0,
						"created_before":             "",
						"custom_time_before":         "",
						"days_since_custom_time":     0,
						"days_since_noncurrent_time": 0,
						"matches_prefix":             []interface{}{},
						"matches_storage_class":      []interface{}{},
						"matches_suffix":             []interface{}{},
						"noncurrent_time_before":     "",
						"num_newer_versions":         3,
						"with_state":                 "ARCHIVED",
					},
				},
			},
			map[string]interface{}{
				"action": []interface{}{
					map[string]interface{}{
						"storage_class": "",
						"type":          "Delete",
					},
				},
				"condition": []interface{}{
					map[string]interface{}{
						"age":                        float64(0),
						"created_before":             "",
						"custom_time_before":         "",
						"days_since_custom_time":     float32(0),
						"days_since_noncurrent_time": 7,
						"matches_prefix":             []interface{}{},
						"matches_storage_class":      []interface{}{},
						"matches_suffix":             []interface{}{},
						"noncurrent_time_before":     "",
						"num_newer_versions":         0,
						"with_state":                 "ANY",
					},
				},
			},
		},
		"location":                 "US",
		"logging":                  []interface{}{},
		"name":                     "dragondrop-modules",
		"project":                  "dragondrop-dev",
		"public_access_prevention": "enforced",
		"requester_pays":           false,
		"retention_policy":         []interface{}{},
		"self_link":                "https://www.googleapis.com/storage/v1/b/dragondrop-modules",
		"storage_class":            "STANDARD",
		"timeouts": map[string]interface{}{
			"create": nil,
			"read":   nil,
			"update": nil,
		},
		"uniform_bucket_level_access": true,
		"url":                         "gs://dragondrop-modules",
		"versioning": []interface{}{
			map[string]interface{}{
				"enabled": true,
			},
		},
		"website": []interface{}{},
	}

	// When
	attributes, err := convertNestedMapToFlatAttributes(nestedAttributes)

	// Then
	require.NoError(t, err)
	require.NotNil(t, attributes)

	require.Equal(t, "false", attributes["default_event_based_hold"])
	require.Equal(t, "false", attributes["force_destroy"])
	require.Equal(t, "0", attributes["lifecycle_rule.1.condition.0.age"])
	require.Equal(t, "0", attributes["lifecycle_rule.1.condition.0.days_since_custom_time"])
	require.Equal(t, "0", attributes["lifecycle_rule.1.condition.0.num_newer_versions"])
	require.NotContains(t, "timeouts.create", attributes)
	require.NotContains(t, "timeouts.read", attributes)
	require.NotContains(t, "timeouts.update", attributes)
}

func TestExtractRegionFromAWSAttributes(t *testing.T) {
	// Given, baseline region variable exists in map
	awsAttributes := map[string]string{
		"region": "eu-west-1",
	}

	// When
	region, err := extractRegionFromAWSAttributes(awsAttributes)

	// Then
	require.NoError(t, err)
	require.Equal(t, "eu-west-1", region)

	// Given, no region key in map and arn variable exists
	awsAttributes = map[string]string{
		"arn": "arn:aws:ec2:eu-west-2:123456789012:instance/i-1234567890abcdef0",
	}

	// When
	region, err = extractRegionFromAWSAttributes(awsAttributes)

	// Then
	require.NoError(t, err)
	require.Equal(t, "eu-west-2", region)

	// Given no region key in map and arn variable for an iam role exists
	awsAttributes = map[string]string{
		"arn": "arn:aws:iam::123456789012:role/test-role",
	}

	// When
	region, err = extractRegionFromAWSAttributes(awsAttributes)

	// Then
	require.NoError(t, err)
	require.Equal(t, "us-east-1", region)

	// Given no region key or arn variable in map
	awsAttributes = map[string]string{}

	// When
	region, err = extractRegionFromAWSAttributes(awsAttributes)

	// Then
	require.NoError(t, err)
	require.Equal(t, "us-east-1", region)
}
