package hclcreate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// WriteImportBlocks writes import blocks to .tf files for
// configurations using Terraform version 1.5.0 or higher.
func (h *hclCreate) WriteImportBlocks(uniqueID string, workspaceToDirectory map[string]string) error {
	// load in resource to import location map
	resourceToImportLoc, err := os.ReadFile("outputs/resources-to-import-location.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile] outputs/resources-to-import-location.json error: %v", err)
	}

	resourcetoImportDataPair := ResourceToImportDataPair{}
	err = json.Unmarshal(resourceToImportLoc, &resourcetoImportDataPair)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] error unmarshalling `resourceToImportLoc`: %v", err)
	}

	// load in resource to workspace map
	resourceToWorkspace, err := os.ReadFile("outputs/new-resources-to-workspace.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile] outputs/new-resources-to-workspace.json error: %v", err)
	}

	newResourceToWorkspace := NewResourceToWorkspace{}
	err = json.Unmarshal(resourceToWorkspace, &newResourceToWorkspace)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] error unmarshalling `resourceToWorkspace`: %v", err)
	}

	workspacesWithMigrations := h.setOfWorkspacesWithMigrationsStruct(newResourceToWorkspace)

	for workspace, directory := range workspaceToDirectory {
		if _, ok := workspacesWithMigrations[workspace]; !ok {
			continue
		}

		importBlockFileBytes, err := h.generateImportBlockFile(
			workspace,
			resourcetoImportDataPair,
			newResourceToWorkspace,
		)
		if err != nil {
			return fmt.Errorf("[h.generateImportBlockFile]%v", err)
		}

		err = os.MkdirAll(fmt.Sprintf("repo%vcloud-concierge/imports", directory), 0o400)
		if err != nil {
			return fmt.Errorf("[os.MkdirAll] error making directory: %v", err)
		}
		// outputting the file
		outputPath := fmt.Sprintf("repo%vcloud-concierge/imports/%v_imports.tf", directory, uniqueID)
		err = os.WriteFile(outputPath, importBlockFileBytes, 0o400)
		if err != nil {
			return fmt.Errorf("[os.WriteFile] Error writing %v: %v", outputPath, err)
		}
	}

	return nil
}

// setOfWorkspacesWithMigrations returns a set of workspaces that have associated new resources
func (h *hclCreate) setOfWorkspacesWithMigrationsStruct(resourceToWorkspace NewResourceToWorkspace) map[string]bool {
	workspacesWithMigration := map[string]bool{}

	for _, workspace := range resourceToWorkspace {
		workspacesWithMigration[workspace] = true
	}

	return workspacesWithMigration
}

// generateImportBlockFile generates a .tf file containing import blocks for
// all resources within a workspace that are to be imported.
func (h *hclCreate) generateImportBlockFile(
	workspace string,
	resourceToImportLocation ResourceToImportDataPair,
	resourceToWorkspace NewResourceToWorkspace,
) ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	fBody := f.Body()

	for resource, currentWorkspace := range resourceToWorkspace {
		if currentWorkspace == workspace {
			currentResource := h.resourceToIdentifierStruct(resource)
			resourceID := fmt.Sprintf("%v.%v", currentResource.resourceType, currentResource.resourceName)
			currentImportDataPair := resourceToImportLocation[resourceID]
			fBody = h.hclImportBlock(fBody, currentImportDataPair)
		}
	}

	return f.Bytes(), nil
}

// hclImportBlock writes an import block to the passed-in hclwrite body.
func (h *hclCreate) hclImportBlock(body *hclwrite.Body, importDataPair ImportDataPair) *hclwrite.Body {
	importBlock := body.AppendNewBlock(
		"import", nil)
	importBlock.Body().SetAttributeValue(
		"to",
		cty.StringVal(strings.Replace(importDataPair.TerraformConfigLocation, "tfer--", "", -1)),
	)
	importBlock.Body().SetAttributeValue("id", cty.StringVal(importDataPair.RemoteCloudReference))
	return body
}
