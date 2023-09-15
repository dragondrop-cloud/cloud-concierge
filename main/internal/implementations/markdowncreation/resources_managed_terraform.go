package markdowncreation

import (
	"fmt"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// StateFileName represents a state file name
type StateFileName string

// ResourcePath represents a resource path
type ResourcePath string

// InstanceID represents an instance id from a resource
type InstanceID string

// InstanceDriftedResources represents a map of instance id to a list of drifted resources
type InstanceDriftedResources map[InstanceID][]ManagedDriftResource

// ResourcePathDriftedResources represents a map of resource path to a list of drifted resources
type ResourcePathDriftedResources map[ResourcePath]InstanceDriftedResources

// StateFileDriftedResources represents a map of state file name to a list of drifted resources
type StateFileDriftedResources map[StateFileName]ResourcePathDriftedResources

// setDriftedResourcesManagedByTerraformData sets the drifted resources managed by terraform data in the markdown report
func (m *MarkdownCreator) setDriftedResourcesManagedByTerraformData(report *doc.MarkDownDoc) {
	report.Write("# Drifted Resources Managed By Terraform").Writeln().Writeln()

	if len(m.managedDrift) == 0 {
		report.Write("No controlled resources have drifted!").Writeln()
		return
	}

	stateFileDriftedResources := StateFileDriftedResources{}
	for _, driftedResource := range m.managedDrift {
		resourcePath := fmt.Sprintf("%s (module) \"%s\" \"%s\"", driftedResource.ModuleName, driftedResource.ResourceType, driftedResource.ResourceName)

		if stateFileDriftedResources[StateFileName(driftedResource.StateFileName)] == nil {
			stateFileDriftedResources[StateFileName(driftedResource.StateFileName)] = ResourcePathDriftedResources{
				ResourcePath(resourcePath): InstanceDriftedResources{
					InstanceID(driftedResource.InstanceID): []ManagedDriftResource{},
				},
			}
		}
		if stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)] == nil {
			stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)] = InstanceDriftedResources{
				InstanceID(driftedResource.InstanceID): []ManagedDriftResource{},
			}
		}
		if stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)][InstanceID(driftedResource.InstanceID)] == nil {
			stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)][InstanceID(driftedResource.InstanceID)] = []ManagedDriftResource{}
		}

		stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)][InstanceID(driftedResource.InstanceID)] = append(
			stateFileDriftedResources[StateFileName(driftedResource.StateFileName)][ResourcePath(resourcePath)][InstanceID(driftedResource.InstanceID)],
			driftedResource,
		)
	}

	for stateFileName, resourcePathDriftedResources := range stateFileDriftedResources {
		report.Write(fmt.Sprintf("## State File `%s`", stateFileName)).Writeln().Writeln()

		for resourcePath, instanceDriftedResources := range resourcePathDriftedResources {
			report.Write(fmt.Sprintf("### Resource: %s", resourcePath)).Writeln().Writeln()

			for instanceID, driftedResources := range instanceDriftedResources {
				report.Write(fmt.Sprintf("**Instance ID**: `%s`", instanceID)).Writeln().Writeln()
				report.Write(fmt.Sprintf("**Most Recent Non-Terraform Actor**: `%s`", driftedResources[0].RecentActor)).Writeln()
				report.Write(fmt.Sprintf("**Most Recent Action Date**: `%s`", driftedResources[0].RecentActionTimestamp)).Writeln().Writeln()
				report.Write("- [ ] Completed").Writeln().Writeln()

				report.Write("|Attribute|Terraform Value|Cloud Value|\n| :---: | :---: | :---: |\n")

				for _, driftedResource := range driftedResources {
					report.Write(fmt.Sprintf("|%s", driftedResource.AttributeName))
					report.Write(fmt.Sprintf("|%s", driftedResource.TerraformValue))
					report.Write(fmt.Sprintf("|%s|", driftedResource.CloudValue)).Writeln()
				}
				report.Writeln()
			}
		}
	}
}
