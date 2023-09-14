package hclcreate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

// CreateTFMigrate coordinates CreateTFMigrateConfiguration and CreateTFMigrateMigration to create the needed
// components for TFMigrate to operate successfully.
func (h *hclCreate) CreateTFMigrate(uniqueID string, workspaceToDirectory map[string]string) error {
	// load in resource to import location map
	resourceToImportLoc, err := os.ReadFile("outputs/resources-to-import-location.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile] outputs/resources-to-import-location.json error: %v", err)
	}
	logrus.Debugf("[hclcreate][write_import_blocks] resourceToImportLoc: %v", string(resourceToImportLoc))

	resourceToImportDataPair := ResourceToImportDataPair{}
	err = json.Unmarshal(resourceToImportLoc, &resourceToImportDataPair)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] error unmarshalling `resourceToImportLoc`: %v", err)
	}

	// load in resource to workspace map
	resourceToWorkspace, err := os.ReadFile("outputs/new-resources-to-workspace.json")
	if err != nil {
		return fmt.Errorf("[os.ReadFile] outputs/new-resources-to-workspace.json error: %v", err)
	}
	logrus.Debugf("[hclcreate][write_import_blocks] resourceToWorkspace: %v", string(resourceToWorkspace))

	newResourceToWorkspace := NewResourceToWorkspace{}
	err = json.Unmarshal(resourceToWorkspace, &newResourceToWorkspace)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] error unmarshalling `resourceToWorkspace`: %v", err)
	}

	err = h.CreateTFMigrateConfiguration(workspaceToDirectory)
	if err != nil {
		return fmt.Errorf("[h.CreateTFMigrateConfiguration] %v", err)
	}

	err = h.CreateTFMigrateMigration(
		uniqueID,
		resourceToImportDataPair,
		newResourceToWorkspace,
		workspaceToDirectory,
	)
	if err != nil {
		return fmt.Errorf("[h.CreateTFMigrateMigration] %v", err)
	}

	return nil
}

// CreateTFMigrateConfiguration saves HCL which defines TFMigrate configuration.
func (h *hclCreate) CreateTFMigrateConfiguration(workspaceToDirectory map[string]string) error {
	for workspace, directory := range workspaceToDirectory {
		logrus.Debugf("[hclcreate][write_import_blocks] workspace: %v", workspace)

		err := os.MkdirAll(fmt.Sprintf("repo%vcloud-concierge/tfmigrate", directory), 0400)
		if err != nil {
			return fmt.Errorf("[os.MkdirAll] cloud-concierge/tfmigrate within %v: %v", directory, err)
		}

		newFilePath := fmt.Sprintf("repo%vcloud-concierge/tfmigrate/.tfmigrate.hcl", directory)

		currentTfMigrateConfig, err := h.individualTFMigrateConfig(workspace)
		if err != nil {
			return fmt.Errorf("[h.individualTFMigrateConfig] %v", err)
		}

		err = os.WriteFile(newFilePath, currentTfMigrateConfig, 0400)
		if err != nil {
			return fmt.Errorf("[os.writeFile] %v", err)
		}

	}
	return nil
}

// individualTFMigrateConfig creates the []byte representing a tf migrate configuration for an individual
func (h *hclCreate) individualTFMigrateConfig(workspace string) ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	tfmigrateBlock := body.AppendNewBlock("tfmigrate", nil)

	tfmigrateBlockBody := tfmigrateBlock.Body()

	tfmigrateBlockBody.SetAttributeValue("migration_dir", cty.StringVal("./cloud-concierge/tfmigrate/"))
	tfmigrateBlockBody.SetAttributeValue("is_backend_terraform_cloud", cty.BoolVal(true))

	historyBlock := tfmigrateBlockBody.AppendNewBlock("history", nil)
	historyBlockBody := historyBlock.Body()

	storageType := h.config.MigrationHistoryStorage.StorageType
	storageBlock := historyBlockBody.AppendNewBlock(
		"storage", []string{storageType},
	)
	storageBlockBody := storageBlock.Body()

	historyKey := fmt.Sprintf("%v/history.json", workspace)

	switch storageType {
	case "gcs":
		storageBlockBody.SetAttributeValue("bucket", cty.StringVal(h.config.MigrationHistoryStorage.Bucket))
		storageBlockBody.SetAttributeValue("name", cty.StringVal(historyKey))
	case "s3":
		storageBlockBody.SetAttributeValue("bucket", cty.StringVal(h.config.MigrationHistoryStorage.Bucket))
		storageBlockBody.SetAttributeValue("key", cty.StringVal(historyKey))
		storageBlockBody.SetAttributeValue("region", cty.StringVal(h.config.MigrationHistoryStorage.Region))
	default:
		return nil, fmt.Errorf("tfmigrate storage type of %v passed, only s3 is currently supported", storageType)
	}

	return f.Bytes(), nil
}

// CreateTFMigrateMigration saves HCL which defines a TFMigrate migration.
func (h *hclCreate) CreateTFMigrateMigration(
	uniqueID string,
	resourceToImportDataPair ResourceToImportDataPair,
	newResourceToWorkspace NewResourceToWorkspace,
	workspaceToDirectory map[string]string,
) error {
	workspacesWithMigrations := h.setOfWorkspacesWithMigrationsStruct(newResourceToWorkspace)

	// complete one workspace migration file at a time
	for workspace, directory := range workspaceToDirectory {
		if !workspacesWithMigrations[workspace] {
			continue
		}
		// create the workspace migration file
		migrationFileBytes, err := h.individualTFMigrateMigration(
			directory,
			workspace,
			resourceToImportDataPair,
			newResourceToWorkspace,
		)
		if err != nil {
			return fmt.Errorf("[h.individualTFMigrateMigration] %v", err)
		}

		// outputting the file
		outputPath := fmt.Sprintf("repo%vcloud-concierge/tfmigrate/%v_migrations.hcl", directory, uniqueID)
		err = os.WriteFile(outputPath, migrationFileBytes, 0400)
		if err != nil {
			return fmt.Errorf("[os.WriteFile] Error writing %v: %v", outputPath, err)
		}
	}

	return nil
}

// individualTFMigrateMigration creates a TFMigrateMigration file for the specified workspace.
func (h *hclCreate) individualTFMigrateMigration(
	directory string,
	workspace string,
	resourceToImportDataPair ResourceToImportDataPair,
	newResourceToWorkspace NewResourceToWorkspace,
) ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	fBody := f.Body()

	// set up the base of the file
	migrationBlock := fBody.AppendNewBlock("migration", []string{"state", "import"})
	migrationBlockBody := migrationBlock.Body()

	dirVal := fmt.Sprintf("/github/workspace%v", directory)
	migrationBlockBody.SetAttributeValue("dir", cty.StringVal(dirVal))
	migrationBlockBody.SetAttributeValue("workspace", cty.StringVal(workspace))

	var importStatementSlice []cty.Value

	// Generate a list of import statements for resources.
	for resource, workspaceName := range newResourceToWorkspace {
		if workspaceName == workspace {
			importStatement, err := h.generateImportStatement(
				resource,
				resourceToImportDataPair,
			)
			if err != nil {
				return nil, fmt.Errorf("[h.generateImportStatement] Error with resource %v: %v", resource, err)
			}

			importStatementSlice = append(importStatementSlice, cty.StringVal(importStatement))
		}
	}

	migrationBlockBody.SetAttributeValue("actions", cty.ListVal(importStatementSlice))

	return f.Bytes(), nil
}

// generateImportStatement generates the text for an import statement for the specified resource from
// resourceToImportLocation.
func (h *hclCreate) generateImportStatement(
	resource string,
	resourceToImportDataPair ResourceToImportDataPair,
) (string, error) {
	resourceIDStruct := h.resourceToIdentifierStruct(resource)

	resourceImportData := resourceToImportDataPair[fmt.Sprintf("%v.%v", resourceIDStruct.resourceType, resourceIDStruct.resourceName)]

	importText := h.generateImportStatementText(resourceImportData.RemoteCloudReference, resourceIDStruct)
	logrus.Debugf("Import statement for %v: %v", resource, importText)
	return importText, nil
}

// generateImportStatementText generates the final input statement text for a given cloud resource needing to be
// imported into terraform control
func (h *hclCreate) generateImportStatementText(remoteCloudReference string, resourceIDStruct ResourceIdentifier) string {
	cleanedResourceName := ConvertTerraformerResourceName(resourceIDStruct.resourceName)

	return fmt.Sprintf("import %v.%v %v", resourceIDStruct.resourceType, cleanedResourceName, remoteCloudReference)
}

// resourceToIdentifierStruct structures the information found within the resource string
// within a ResourceIdentifier struct.
func (h *hclCreate) resourceToIdentifierStruct(resource string) ResourceIdentifier {
	resourceComponents := strings.Split(resource, ".")
	return ResourceIdentifier{
		resourceType: resourceComponents[0],
		resourceName: resourceComponents[1],
	}
}
