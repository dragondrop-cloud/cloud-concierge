package hclcreate

import (
	"fmt"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// ExtractResourceDefinitions outputs a bytes slice which defines needed Terraform resources extracted from
// another configuration file.
func (h *hclCreate) ExtractResourceDefinitions(noNewResources bool, workspaceToDirectory map[string]string) error {
	logrus.Debugf("[hclcreate][ExtractResourceDefinitions] noNewResources: %v", noNewResources)

	// Mapping between workspace and the new hclwrite document for that workspace
	workspaceToHCLFile := WorkspaceToHCL{}

	for workspace := range workspaceToDirectory {
		workspaceToHCLFile[workspace] = hclwrite.NewEmptyFile()
	}

	rawCloudActions, err := os.ReadFile("outputs/resources-to-cloud-actions.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile resources-to-cloud-actions.json]%v", err)
	}

	parsedCloudActions, err := gabs.ParseJSON(rawCloudActions)
	if err != nil {
		return fmt.Errorf("[gabs.ParseJSON rawCloudActions]%v", err)
	}

	rawCloudCosts, err := os.ReadFile("outputs/cost-estimates.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile cost-estimates.json]%v", err)
	}

	parsedCloudCosts, err := gabs.ParseJSON(rawCloudCosts)
	if err != nil {
		return fmt.Errorf("[gabs.ParseJSON rawCloudCosts]%v", err)
	}

	costEstimates, err := gabsContainerToCostsStruct(parsedCloudCosts)
	if err != nil {
		return fmt.Errorf("[gabsContainerToAllCostsStruct]%v", err)
	}

	hclBytes, err := os.ReadFile("current_cloud/resources.tf")
	if err != nil {
		return fmt.Errorf("[os.ReadFile()] Error reading in resources.tf")
	}

	terraformerResources, hclDiagnostics := hclwrite.ParseConfig(
		hclBytes,
		"cloud-resources.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)

	if hclDiagnostics != nil {
		return fmt.Errorf("[hclwrite.ParseConfig]%v", hclDiagnostics)
	}

	resourceActions, err := h.cloudActionsToResourceActionMap(parsedCloudActions)
	if err != nil {
		return fmt.Errorf("[h.subsetCloudActionsToCurrentDivision]%v", err)
	}

	// Read in new-resources-to-workspace.json, parse as gabs file
	newResourcesToWorkspace := []byte("{}")
	if !noNewResources {
		newResourcesToWorkspace, err = os.ReadFile("outputs/new-resources-to-workspace.json")
		if err != nil {
			return fmt.Errorf("[os.ReadFile()] Error reading in new-resources-to-workspace.json: %v", err)
		}
	}

	parsedNewResourceToWorkspace, err := gabs.ParseJSON(newResourcesToWorkspace)

	if err != nil {
		return fmt.Errorf("[gabs.ParseJSON] Error parsing new-resources-to-workspace.json")
	}

	completeWorkspaceToHCLFile, err := h.placeHCLIntoNewFileDef(
		resourceActions,
		costEstimates,
		terraformerResources,
		parsedNewResourceToWorkspace,
		workspaceToHCLFile,
	)

	if err != nil {
		return fmt.Errorf("[h.placeHCLIntoNewFileDef] %v", err)
	}

	err = h.writeNewResourceFiles(
		workspaceToDirectory,
		completeWorkspaceToHCLFile,
	)
	if err != nil {
		return fmt.Errorf("[h.writeNewHCLFiles] %v", err)
	}

	return nil
}

// cloudActionsToResourceActionMap takes in a gabs.Container, and converts to a ResourceActionMap
func (h *hclCreate) cloudActionsToResourceActionMap(parsedCloudActions *gabs.Container) (
	terraformValueObjects.ResourceActionMap, error,
) {
	resourceActionMap := terraformValueObjects.ResourceActionMap{}

	for resourceName, resourceActions := range parsedCloudActions.ChildrenMap() {
		currentResourceActions := terraformValueObjects.ResourceActions{}
		if resourceActions.Exists("creation") {
			currentResourceActions.Creator = &terraformValueObjects.CloudActorTimeStamp{
				Actor:     terraformValueObjects.CloudActor(resourceActions.Search("creation", "actor").Data().(string)),
				Timestamp: terraformValueObjects.Timestamp(resourceActions.Search("creation", "timestamp").Data().(string)),
			}
		}
		if resourceActions.Exists("modified") {
			currentResourceActions.Modifier = &terraformValueObjects.CloudActorTimeStamp{
				Actor:     terraformValueObjects.CloudActor(resourceActions.Search("modified", "actor").Data().(string)),
				Timestamp: terraformValueObjects.Timestamp(resourceActions.Search("modified", "timestamp").Data().(string)),
			}
		}

		resourceActionMap[terraformValueObjects.ResourceName(resourceName)] = &currentResourceActions
	}

	return resourceActionMap, nil
}

// placeHCLIntoNewFileDef transfers the relevant HCL created by terraformer
// into the new file definition.
func (h *hclCreate) placeHCLIntoNewFileDef(
	cloudActions terraformValueObjects.ResourceActionMap,
	costEstimates costs,
	terraformerResources *hclwrite.File,
	parsedNewResourceToWorkspace *gabs.Container,
	workspaceToHCLFile WorkspaceToHCL,
) (WorkspaceToHCL, error) {
	for resource, workspaceName := range parsedNewResourceToWorkspace.ChildrenMap() {
		resourceID := h.splitResourceIdentifier(resource)

		cleanResourceName := ConvertTerraformerResourceName(resourceID.resourceName)
		extractedBlock, err := h.extractResourceBlockDefinition(
			terraformerResources,
			cleanResourceName,
			resourceID,
		)
		if err != nil {
			return nil, fmt.Errorf("[h.extractResourceBlockDefinition] %v", err)
		}

		cloudIdentifierComment := h.generateHCLCloudActorsComment(resourceID.resourceType, cleanResourceName, cloudActions)

		cloudCostComment := h.generateHCLCloudCostComment(resourceID.resourceType, cleanResourceName, costEstimates)

		// place resource within the corresponding workspace's file.
		workspaceNameString := workspaceName.Data().(string)

		workspaceToHCLFile[workspaceNameString] = h.writeBlockToWorkspaceHCL(
			workspaceToHCLFile[workspaceNameString],
			cloudIdentifierComment,
			cloudCostComment,
			extractedBlock,
		)
	}

	return workspaceToHCLFile, nil
}

// generateHCLCloudActorsComment generates data on Cloud Actor actions for the specified resource.
func (h *hclCreate) generateHCLCloudActorsComment(
	resourceType string, resourceName string,
	resourceToCloudActions terraformValueObjects.ResourceActionMap,
) hclwrite.Tokens {
	logrus.Debugf("[hclcreate][generateHCLCloudActorsComment] resourceType: %v, resourceName: %v", resourceType, resourceName)

	completeResourceName := fmt.Sprintf("%v.%v", resourceType, resourceName)
	cloudActorActionStatement := ""
	cloudActions, ok := resourceToCloudActions[terraformValueObjects.ResourceName(completeResourceName)]
	if ok {
		creatorLine := ""
		modifierLine := ""
		if cloudActions.Creator != nil && cloudActions.Creator.Actor != "" {
			creatorLine = fmt.Sprintf("\n# Created at %v by %v", cloudActions.Creator.Timestamp, cloudActions.Creator.Actor)
		}
		if cloudActions.Modifier != nil && cloudActions.Modifier.Actor != "" {
			modifierLine = fmt.Sprintf("\n# Last Modified at %v by %v", cloudActions.Modifier.Timestamp, cloudActions.Modifier.Actor)
		}
		cloudActorActionStatement = fmt.Sprintf("%v%v", creatorLine, modifierLine)
	}
	return hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(cloudActorActionStatement),
			SpacesBefore: 0,
		},
	}
}

// extractResourceBlockDefinition pulls the resource block from the specified hclFile and renames it
// to have the name specified by the cleanResourceName variable.
func (h *hclCreate) extractResourceBlockDefinition(
	hclFile *hclwrite.File,
	cleanResourceName string,
	resourceID ResourceIdentifier,
) (*hclwrite.Block, error) {
	body := hclFile.Body()

	typeName := "resource"
	labels := []string{resourceID.resourceType, resourceID.resourceName}

	extractBlock := body.FirstMatchingBlock(typeName, labels)

	if extractBlock == nil {
		return nil, fmt.Errorf("could not find block matching %v, although it was expected", resourceID)
	}

	labels[1] = cleanResourceName

	extractBlock.SetLabels(labels)

	return extractBlock, nil
}

// writeBlockToWorkspaceHCL writes the needed information for a new block to the
// current workspace's HCL Body.
func (h *hclCreate) writeBlockToWorkspaceHCL(
	hclFile *hclwrite.File,
	cloudIdentifierComment hclwrite.Tokens,
	cloudCostComment hclwrite.Tokens,
	extractedBlock *hclwrite.Block,
) *hclwrite.File {
	fileBody := hclFile.Body()
	fileBody.AppendUnstructuredTokens(cloudCostComment)
	fileBody.AppendUnstructuredTokens(cloudIdentifierComment)
	fileBody.AppendNewline()
	fileBody.AppendBlock(extractedBlock)
	fileBody.AppendNewline()

	return hclFile
}

// splitResourceIdentifier takes information from the resourceIdentifier string and outputs it
// organized within the ResourceIdentifier struct.
func (h *hclCreate) splitResourceIdentifier(resourceIdentifier string) ResourceIdentifier {
	resourceIDSlice := strings.Split(resourceIdentifier, ".")

	return ResourceIdentifier{
		resourceType: resourceIDSlice[0],
		resourceName: resourceIDSlice[1],
	}
}

// writeNewResourceFiles takes the hcl files from workspaceToHCLFile and outputs each to the appropriate directory
// as informed by completeWorkspaceToHCLFile.
func (h *hclCreate) writeNewResourceFiles(
	workspaceToDirectoryMap map[string]string, completeWorkspaceToHCLFile WorkspaceToHCL,
) error {
	for workspace, hclFile := range completeWorkspaceToHCLFile {
		fileContent := hclwrite.Format(hclFile.Bytes())

		subDirectory := workspaceToDirectoryMap[workspace]

		if string(fileContent) != "" {
			filePath := fmt.Sprintf("repo%vnew-resources.tf", subDirectory)

			err := os.WriteFile(filePath, fileContent, 0400)

			if err != nil {
				return fmt.Errorf(
					"[os.WriteFile] Error for repo%vnew-resources.tf:  %v",
					subDirectory,
					err,
				)
			}
		}
	}
	return nil
}
