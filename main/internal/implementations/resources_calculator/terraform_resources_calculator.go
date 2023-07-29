package resourcesCalculator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/documentize"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/pyscriptexec"
)

var ErrNoNewResources = errors.New("[no new resources identified]")

// TerraformResourcesCalculator is a struct that implements the interfaces.ResourcesCalculator interface for
// running within a "live" dragondrop job.
type TerraformResourcesCalculator struct {
	// documentize implements the Document
	documentize *documentize.Documentize

	// pyScriptExec is an implementation of the PyScriptExec python scripts
	pyScriptExec pyscriptexec.PyScriptExec

	// dragonDrop interface implementation for sending requests to the dragondrop API.
	dragonDrop interfaces.DragonDrop
}

// ResourceID is a string that represents a resource id for a cloud resource within a terraform state file.
type ResourceID string

// DivisionToNewResources is a mapping of a division
// to a map of resource ids to defining resource data.
type DivisionToNewResources map[terraformValueObjects.Division]map[ResourceID]NewResourceData

// NewResourceData is a struct that contains key fields defining a Terraform resource.
type NewResourceData struct {
	ResourceType            string `json:"ResourceType"`
	ResourceTerraformerName string `json:"ResourceTerraformerName"`
	Region                  string `json:"Region"`
}

// NewTerraformResourcesCalculator creates and returns an instance of the TerraformResourcesCalculator.
func NewTerraformResourcesCalculator(documentize *documentize.Documentize, pyScriptExec pyscriptexec.PyScriptExec, dragonDrop interfaces.DragonDrop) interfaces.ResourcesCalculator {
	return &TerraformResourcesCalculator{documentize: documentize, pyScriptExec: pyScriptExec, dragonDrop: dragonDrop}
}

// Execute calculates the association between resources and a state file.
func (c *TerraformResourcesCalculator) Execute(ctx context.Context, workspaceToDirectory map[string]string) error {
	_, err := c.calculateResourceToWorkspaceMapping(ctx, *c.documentize, workspaceToDirectory)
	if err != nil {
		if errors.Unwrap(err) == ErrNoNewResources {
			err := c.dragonDrop.InformNoResourcesFound(ctx)
			if err != nil {
				return fmt.Errorf("[resources_calculator][error informing no new resources identified]%w", err)
			}
		}

		return fmt.Errorf("[resources_calculator][error calculating resources to workspace]%w", err)
	}
	return nil
}

// calculateResourceToWorkspaceMapping determines which resources need to be added
// and to which workspaces.
func (c *TerraformResourcesCalculator) calculateResourceToWorkspaceMapping(ctx context.Context, docu documentize.Documentize, workspaceToDirectory map[string]string) (string, error) {

	message, err := c.createWorkspaceDocuments(ctx, docu, workspaceToDirectory)
	if err != nil {
		return message, fmt.Errorf("[calculate_resource_to_workspace_mapping][error creating workspace documents]%w", err)
	}

	newResources, err := c.identifyNewResources(ctx, docu, workspaceToDirectory)
	if err != nil {
		return message, err
	}

	if len(newResources) == 0 {
		fmt.Println("No new resources identified")
		return "no new resources", fmt.Errorf("[calculate_resource_to_workspace][error identifying new resources]%w", ErrNoNewResources)
	}

	err = c.createNewResourceDocuments(ctx, docu, newResources)
	if err != nil {
		return message, err
	}

	err = c.getResourceToWorkspaceMapping(ctx)
	if err != nil {
		return message, err
	}

	return "", nil
}

// getResourceToWorkspaceMapping runs the NLPEngine python script to produce a mapping of new resources to suggested workspace.
func (c *TerraformResourcesCalculator) getResourceToWorkspaceMapping(ctx context.Context) error {
	c.dragonDrop.PostLog(ctx, "Beginning to calculate recommended placement of resources to workspace.")
	err := c.pyScriptExec.RunNLPEngine()

	if err != nil {
		return fmt.Errorf("[get_resource_to_workspace][pse.RunNLPEngine]%w", err)
	}

	c.dragonDrop.PostLog(ctx, "Done making a map of workspaces to documents.")
	return nil
}

// createNewResourceDocuments defines documents out of new resources to be used in downstream processing
// like NLP modeling and cloud actor action querying.
func (c *TerraformResourcesCalculator) createNewResourceDocuments(ctx context.Context, docu documentize.Documentize, newResources map[terraformValueObjects.Division]map[documentize.ResourceData]bool) error {
	c.dragonDrop.PostLog(ctx, "Beginning to create new resource documents.")

	newResourceDocs, err := docu.NewResourceDocuments(newResources)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][docu.NewResourceDocuments]%w", err)
	}

	resourceDocsJSON, err := docu.ConvertNewResourcesToJSON(newResourceDocs)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][docu.ConvertNewResourcesToJSON] Error: %v", err)
	}

	err = os.WriteFile("mappings/new-resources-to-documents.json", resourceDocsJSON, 0400)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][write mappings/new-resources-to-documents.json] Error: %v", err)
	}

	gabsContainer, terraformerBytes, err := c.createDivisionToTerraformerStateMap(resourceDocsJSON)
	if err != nil {
		return fmt.Errorf("[createDivisionToTerraformerStateMap]%v", err)
	}

	newResourceData, err := c.createNewResourceData(gabsContainer, terraformerBytes)
	if err != nil {
		return fmt.Errorf("[createDivisionToNewResourceData]%v", err)
	}

	newResourceDataJSON, err := json.MarshalIndent(newResourceData, "", "  ")
	if err != nil {
		return fmt.Errorf("[json.MarshalIndent]%v", err)
	}

	err = os.WriteFile("mappings/new-resources.json", newResourceDataJSON, 0400)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][write mappings/new-resources.json] Error: %v", err)
	}

	c.dragonDrop.PostLog(ctx, "Done creating new resource documents.")
	return nil
}

// createDivisionToNewResourceData loads Terraformer state file bytes
// along with a gabs container of the resource to documents JSON.
func (c *TerraformResourcesCalculator) createDivisionToTerraformerStateMap(resourceDocsJSON []byte) (
	*gabs.Container, driftDetector.TerraformerStateFile, error,
) {
	divisionToTerraformerByteArray := driftDetector.TerraformerStateFile{}

	container, err := gabs.ParseJSON(resourceDocsJSON)
	if err != nil {
		return nil, divisionToTerraformerByteArray, fmt.Errorf("[gabs.ParseJSON]%v", err)
	}

	terraformerContent, err := os.ReadFile(
		fmt.Sprintf("current_cloud/terraform.tfstate"),
	)
	if err != nil {
		return nil, divisionToTerraformerByteArray, fmt.Errorf("[os.ReadFile]%v", err)
	}

	parsedStateFile, err := driftDetector.ParseTerraformerStateFile(terraformerContent)
	if err != nil {
		return nil, divisionToTerraformerByteArray, fmt.Errorf("[driftDetector.ParseTerraformerStateFile]%v", err)
	}

	return container, parsedStateFile, nil
}

// createNewResourceData converts the resourceDocsJSON to a newResources struct.
// This data is saved in downstream operations for subsequent use with cloud actor identification.
func (c *TerraformResourcesCalculator) createNewResourceData(
	container *gabs.Container,
	terraformerStateFile driftDetector.TerraformerStateFile,
) (map[ResourceID]NewResourceData, error) {
	var err error

	newResources := map[ResourceID]NewResourceData{}

	for key := range container.ChildrenMap() {
		divisionTypeNameSlice := strings.Split(key, ".")
		resourceType := divisionTypeNameSlice[1]
		resourceName := divisionTypeNameSlice[2]

		resourceID := ""
		region := ""

		for _, resource := range terraformerStateFile.Resources {
			if resource.Type == resourceType && resource.Name == resourceName {
				cloudProvider := strings.Split(resource.Type, "_")[0]
				attributesFlat := resource.Instances[0].AttributesFlat
				resourceID, err = driftDetector.ResourceIDCalculator(attributesFlat, cloudProvider, resourceType)
				if err != nil {
					return nil, fmt.Errorf("[driftDetector.ResourceIDCalculator]%v", err)
				}
				region, err = driftDetector.ParseRegionFromTfStateMap(
					resource.Instances[0].AttributesFlat,
					cloudProvider,
				)
				if err != nil {
					return nil, fmt.Errorf("[driftDetector.ParseRegionFromTfStateMap]%v", err)
				}
			}
		}

		newResources[ResourceID(resourceID)] = NewResourceData{
			ResourceType:            resourceType,
			ResourceTerraformerName: resourceName,
			Region:                  region,
		}
	}

	return newResources, nil
}

// identifyNewResources compares Terraformer output with workspace state files to determine which
// cloud resources will be new to Terraform control.
func (c *TerraformResourcesCalculator) identifyNewResources(ctx context.Context, docu documentize.Documentize, workspaceToDirectory map[string]string) (
	map[terraformValueObjects.Division]map[documentize.ResourceData]bool, error) {
	c.dragonDrop.PostLog(ctx, "Beginning to identify new Resources.")

	newResources, err := docu.IdentifyNewResources(workspaceToDirectory)

	if err != nil {
		return nil, fmt.Errorf("[identify_new_resources][docu.IdentifyNewResources]%w", err)
	}

	c.dragonDrop.PostLog(ctx, "Done identifying new Resources.")
	return newResources, nil
}

// createWorkspaceDocuments defines documents out of remote workspace state to be used in NLP modeling.
func (c *TerraformResourcesCalculator) createWorkspaceDocuments(ctx context.Context, docu documentize.Documentize, workspaceToDirectory map[string]string) (string, error) {
	c.dragonDrop.PostLog(ctx, "Beginning to make map of workspaces to documents.")

	workspaceToDocument, err := docu.AllWorkspaceStatesToDocuments(workspaceToDirectory)

	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}

	outputBytes, err := docu.ConvertWorkspaceDocumentsToJSON(workspaceToDocument)

	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}

	err = os.WriteFile("mappings/workspace-to-documents.json", outputBytes, 0400)

	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}

	c.dragonDrop.PostLog(ctx, "Done with map between workspaces to documents.")
	return "", nil
}
