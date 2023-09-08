package driftdetector

import (
	"testing"

	"github.com/stretchr/testify/require"

	terraformvalueobjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestManagedResourcesDriftDetector_identifyDeletedResources_SameResources(t *testing.T) {
	// Given
	detector := &ManagedResourcesDriftDetector{}

	// Given
	terraformerResources := TerraformerResourceIDToData{
		"dragondrop-modules": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules",
			},
		},
	}

	terraformStateResources := TerraformStateResourceIDToData{
		"dragondrop-modules": TerraformStateUniqueResourceData{
			Provider:  "google",
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules",
			Module:    "module_name",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules",
			},
		},
	}

	// When
	deletedResources, err := detector.identifyDeletedResources(terraformerResources, terraformStateResources)

	// Then
	require.NoError(t, err)
	require.NotNil(t, deletedResources)
	require.Len(t, deletedResources, 0)
}

func TestManagedResourcesDriftDetector_identifyDeletedResources_NewResources(t *testing.T) {
	// Given
	detector := &ManagedResourcesDriftDetector{}

	// Given
	terraformerResources := TerraformerResourceIDToData{
		"dragondrop-modules": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules",
			},
		},
		"dragondrop-modules-2": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules-new",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules-2",
			},
		},
	}

	terraformStateResources := TerraformStateResourceIDToData{
		"dragondrop-modules": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules",
			},
		},
	}

	// When
	deletedResources, err := detector.identifyDeletedResources(terraformerResources, terraformStateResources)

	// Then
	require.NoError(t, err)
	require.NotNil(t, deletedResources)
	require.Len(t, deletedResources, 0)
}

func TestManagedResourcesDriftDetector_identifyDeletedResources_OneDeletedResource_WhiteList(t *testing.T) {
	// Given
	detector := &ManagedResourcesDriftDetector{
		config: ManagedResourceDriftDetectorConfig{
			ResourcesWhiteList: terraformvalueobjects.ResourceNameList{"google_storage_bucket"},
		},
	}

	// Given
	terraformerResources := TerraformerResourceIDToData{
		"dragondrop-modules": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules",
			},
		},
	}

	terraformStateResources := TerraformStateResourceIDToData{
		"dragondrop-modules": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules",
			},
		},
		"dragondrop-modules-2": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules_old",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules-2",
			},
		},
	}

	// When
	deletedResources, err := detector.identifyDeletedResources(terraformerResources, terraformStateResources)

	// Then
	require.NoError(t, err)
	require.NotNil(t, deletedResources)
	require.Len(t, deletedResources, 1)

	require.Contains(t, deletedResources, DeletedResource{
		InstanceID:    "dragondrop-modules-2",
		StateFileName: "example2.tfstate",
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "dragondrop_modules_old",
	})
}

func TestManagedResourcesDriftDetector_identifyDeletedResources_OneDeletedResource_BlackList(t *testing.T) {
	// Given
	detector := &ManagedResourcesDriftDetector{
		config: ManagedResourceDriftDetectorConfig{
			ResourcesBlackList: terraformvalueobjects.ResourceNameList{"google_storage_bucket"},
		},
	}

	// Given
	terraformerResources := TerraformerResourceIDToData{
		"dragondrop-modules": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules",
			},
		},
	}

	terraformStateResources := TerraformStateResourceIDToData{
		"dragondrop-modules": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules",
			},
		},
		"dragondrop-modules-2": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules_old",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules-2",
			},
		},
	}

	// When
	deletedResources, err := detector.identifyDeletedResources(terraformerResources, terraformStateResources)

	// Then
	require.NoError(t, err)
	require.NotNil(t, deletedResources)
	require.Len(t, deletedResources, 0)
}

func TestManagedResourcesDriftDetector_identifyDeletedResources_OneDeletedResource_NoResourceModeDefined(t *testing.T) {
	// Given
	detector := &ManagedResourcesDriftDetector{
		config: ManagedResourceDriftDetectorConfig{
			ResourcesWhiteList: terraformvalueobjects.ResourceNameList{"google_storage_bucket"},
		},
	}

	// Given
	terraformerResources := TerraformerResourceIDToData{
		"dragondrop-modules": TerraformerUniqueResourceData{
			Type:     "google_storage_bucket",
			Name:     "tfer--dragondrop-modules",
			Module:   "module_name",
			Provider: "google",
			AttributesFlat: map[string]string{
				"id": "dragondrop-modules",
			},
		},
	}

	terraformStateResources := TerraformStateResourceIDToData{
		"dragondrop-modules": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules",
			},
		},
		"dragondrop-modules-2": TerraformStateUniqueResourceData{
			StateFile: "example2.tfstate",
			Type:      "google_storage_bucket",
			Name:      "dragondrop_modules_old",
			Module:    "module_name",
			Provider:  "google",
			Attributes: map[string]interface{}{
				"id": "dragondrop-modules-2",
			},
		},
	}

	// When
	deletedResources, err := detector.identifyDeletedResources(terraformerResources, terraformStateResources)

	// Then
	require.NoError(t, err)
	require.NotNil(t, deletedResources)
	require.Len(t, deletedResources, 1)

	require.Contains(t, deletedResources, DeletedResource{
		InstanceID:    "dragondrop-modules-2",
		StateFileName: "example2.tfstate",
		ModuleName:    "module_name",
		ResourceType:  "google_storage_bucket",
		ResourceName:  "dragondrop_modules_old",
	})
}
