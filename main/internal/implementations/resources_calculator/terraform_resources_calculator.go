package resourcescalculator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/documentize"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

var ErrNoNewResources = errors.New("[no new resources identified]")

// TerraformResourcesCalculator is a struct that implements the interfaces.ResourcesCalculator interface for
// running within a "live" dragondrop job.
type TerraformResourcesCalculator struct {
	// documentize provides functionality to process Terraform state files into structured data
	// that the NLPEngine can handle.
	documentize *documentize.Documentize

	// nlpEngine is the implementation of interfaces.NLPEngine for interacting with
	// the NLP engine endpoint.
	nlpEngine interfaces.NLPEngine
}

// ResourceID is a string that represents a resource id for a cloud resource within a terraform state file.
type ResourceID string

// NewResourceMap is a map of resource ids to defining resource data.
type NewResourceMap map[ResourceID]NewResourceData

// NewResourceData is a struct that contains key fields defining a Terraform resource.
type NewResourceData struct {
	ResourceType            string `json:"ResourceType"`
	ResourceTerraformerName string `json:"ResourceTerraformerName"`
	Region                  string `json:"Region"`
}

// NewTerraformResourcesCalculator creates and returns an instance of the TerraformResourcesCalculator.
func NewTerraformResourcesCalculator(documentize *documentize.Documentize, nlpEngine interfaces.NLPEngine) interfaces.ResourcesCalculator {
	return &TerraformResourcesCalculator{documentize: documentize, nlpEngine: nlpEngine}
}

// Execute calculates the association between resources and a state file.
func (c *TerraformResourcesCalculator) Execute(ctx context.Context, workspaceToDirectory map[string]string) error {
	logrus.Debugf("[resources_calculator][Execute] Starting to calculate resources to workspace mapping.")
	_, err := c.calculateResourceToWorkspaceMapping(ctx, *c.documentize, workspaceToDirectory)
	if err != nil {
		if errors.Unwrap(err) == ErrNoNewResources {
			fmt.Println("No new resources identified")
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
	logrus.Debugf("[resources_calculator][calculateResourceToWorkspaceMapping] Identified %v new resources.", len(newResources))

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

// getResourceToWorkspaceMapping hits the NLPEngine endpoint to receive a mapping of new resources to suggested workspace.
func (c *TerraformResourcesCalculator) getResourceToWorkspaceMapping(ctx context.Context) error {
	err := c.nlpEngine.PostNLPEngine(ctx)
	if err != nil {
		return fmt.Errorf("[postNLPEngine]%w", err)
	}

	return nil
}

// createNewResourceDocuments defines documents out of new resources to be used in downstream processing
// like NLP modeling and cloud actor action querying.
func (c *TerraformResourcesCalculator) createNewResourceDocuments(ctx context.Context, docu documentize.Documentize, newResources map[documentize.ResourceData]bool) error {
	newResourceDocs, err := docu.NewResourceDocuments(newResources)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][docu.NewResourceDocuments]%w", err)
	}
	logrus.Debugf("[resources_calculator][createNewResourceDocuments] Created %v new resource documents.", len(newResourceDocs))

	resourceDocsJSONBytes, err := docu.ConvertNewResourcesToJSON(newResourceDocs)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][docu.ConvertNewResourcesToJSON] Error: %v", err)
	}
	logrus.Debugf("[resources_calculator][createNewResourceDocuments] Created new resource documents JSON.")

	err = os.WriteFile("outputs/new-resources-to-documents.json", resourceDocsJSONBytes, 0o400)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][write outputs/new-resources-to-documents.json] Error: %v", err)
	}

	terraformerParsed, err := c.parseTerraformStateFile()
	if err != nil {
		return fmt.Errorf("[createDivisionToTerraformerStateMap]%v", err)
	}

	newResourceData, err := c.createNewResourceData(resourceDocsJSONBytes, terraformerParsed)
	if err != nil {
		return fmt.Errorf("[createDivisionToNewResourceData]%v", err)
	}

	newResourceDataJSON, err := json.MarshalIndent(newResourceData, "", "  ")
	if err != nil {
		return fmt.Errorf("[json.MarshalIndent]%v", err)
	}

	err = os.WriteFile("outputs/new-resources.json", newResourceDataJSON, 0o400)
	if err != nil {
		return fmt.Errorf("[create_new_resource_documents][write outputs/new-resources.json] Error: %v", err)
	}

	return nil
}

// parseTerraformStateFile parses the terraform state file to a TerraformerStateFile struct.
func (c *TerraformResourcesCalculator) parseTerraformStateFile() (
	driftDetector.TerraformerStateFile, error,
) {
	terraformerByteArray := driftDetector.TerraformerStateFile{}

	terraformerContent, err := os.ReadFile("current_cloud/terraform.tfstate")
	if err != nil {
		return terraformerByteArray, fmt.Errorf("[os.ReadFile]%v", err)
	}

	parsedStateFile, err := driftDetector.ParseTerraformerStateFile(terraformerContent)
	if err != nil {
		return terraformerByteArray, fmt.Errorf("[driftDetector.ParseTerraformerStateFile]%v", err)
	}

	return parsedStateFile, nil
}

// createNewResourceData converts the resourceDocsJSON to a newResources struct.
// This data is saved in downstream operations for subsequent use with cloud actor identification.
func (c *TerraformResourcesCalculator) createNewResourceData(
	resourceDocsJSON []byte,
	terraformerStateFile driftDetector.TerraformerStateFile,
) (map[ResourceID]NewResourceData, error) {
	var err error

	newResources := map[ResourceID]NewResourceData{}

	container, err := gabs.ParseJSON(resourceDocsJSON)
	if err != nil {
		return nil, fmt.Errorf("[gabs.ParseJSON]%v", err)
	}

	for key := range container.ChildrenMap() {
		typeNameSlice := strings.Split(key, ".")
		resourceType := typeNameSlice[0]
		resourceName := typeNameSlice[1]

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
	map[documentize.ResourceData]bool, error,
) {
	newResources, err := docu.IdentifyNewResources(workspaceToDirectory)
	if err != nil {
		return nil, fmt.Errorf("[identify_new_resources][docu.IdentifyNewResources]%w", err)
	}

	return newResources, nil
}

// createWorkspaceDocuments defines documents out of remote workspace state to be used in NLP modeling.
func (c *TerraformResourcesCalculator) createWorkspaceDocuments(ctx context.Context, docu documentize.Documentize, workspaceToDirectory map[string]string) (string, error) {
	workspaceToDocument, err := docu.AllWorkspaceStatesToDocuments(workspaceToDirectory)
	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}

	outputBytes, err := docu.ConvertWorkspaceDocumentsToJSON(workspaceToDocument)
	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}
	logrus.Debugf("[resources_calculator][createWorkspaceDocuments] Created workspace documents JSON.")

	err = os.WriteFile("outputs/workspace-to-documents.json", outputBytes, 0o400)
	if err != nil {
		return "[createWorkspacesToDocuments] %v", err
	}

	return "", nil
}
