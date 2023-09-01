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
		if _, ok := terraformerResources[id]; !ok {
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
