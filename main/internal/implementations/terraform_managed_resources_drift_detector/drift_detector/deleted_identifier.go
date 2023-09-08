package driftdetector

// DeletedResource represents a drifted deleted resource
type DeletedResource struct {
	InstanceID    string
	StateFileName StateFileName
	ModuleName    string
	ResourceType  string
	ResourceName  string
}

// identifyDeletedResources identifies the deleted resources from the current TerraformState TerraformStateResourceIDToData
// compared with the current cloud state obtained with terraformer TerraformerResourceIDToData
func (m *ManagedResourcesDriftDetector) identifyDeletedResources(terraformerResources TerraformerResourceIDToData, terraformResources TerraformStateResourceIDToData) ([]DeletedResource, error) {
	deletedResources := make([]DeletedResource, 0)

	for id, data := range terraformResources {
		if _, found := terraformerResources[id]; !found && m.isValidDeletedResource(data.Type) {
			deletedResource := DeletedResource{
				InstanceID:    data.Attributes["id"].(string),
				StateFileName: StateFileName(data.StateFile),
				ModuleName:    data.Module,
				ResourceType:  data.Type,
				ResourceName:  data.Name,
			}
			deletedResources = append(deletedResources, deletedResource)
		}
	}

	return deletedResources, nil
}

func (m *ManagedResourcesDriftDetector) isValidDeletedResource(resourceType string) bool {
	if m.config.ResourcesWhiteList != nil && len(m.config.ResourcesWhiteList) > 0 {
		for _, resource := range m.config.ResourcesWhiteList {
			if string(resource) == resourceType {
				return true
			}
		}

		return false
	} else if m.config.ResourcesBlackList != nil && len(m.config.ResourcesBlackList) > 0 {
		for _, resource := range m.config.ResourcesBlackList {
			if string(resource) == resourceType {
				return false
			}
		}

		return true
	}

	return true
}
