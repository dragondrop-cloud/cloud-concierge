package hclcreate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// WorkspaceToHCL is a map between workspace names and the corresponding hclwrite File object.
type WorkspaceToHCL map[string]*hclwrite.File

// ResourceIdentifier is a struct comprised of the fields necessary to uniquely identify a
// terraformer-generated resource.
type ResourceIdentifier struct {

	// division is the division within which a cloud resource was identified.
	division string

	// resourceType is the Terraform resource-type.
	resourceType string

	// resourceName is the Terraform resource name.
	resourceName string
}

// MigrationHistory is a map containing information needed for specifying tfmigrate
// history storage appropriately.
type MigrationHistory struct {
	// StorageType is the name of the storage resource used to store migration history files.
	StorageType string

	// Bucket is the name of the bucket-like storage location used to store migration history files.
	Bucket string

	// Region is the region of the bucket for storage used to store migration history files. Only needed
	// when storageType is 's3'
	Region string
}

// Config is a struct comprising the configuration needed for hclCreate
type Config struct {

	// MigrationHistoryStorage is a map containing information needed for specifying tfmigrate
	//// history storage appropriately.
	MigrationHistoryStorage MigrationHistory `required:"true"`

	// TerraformVersion is the version of Terraform used.
	TerraformVersion string `required:"true"`
}

// NewResourceToWorkspace is a map of resource unique id to workspace name
type NewResourceToWorkspace map[string]string

// ResourceImportsByDivision is a map of division name to a DivisionToImportDataPairs map
type ResourceImportsByDivision map[string]DivisionToImportDataPairs

// DivisionToImportDataPairs is a map of resource unique id within a division to ImportDataPair
type DivisionToImportDataPairs map[string]ImportDataPair

// ImportDataPair is a struct that holds the data needed to write an individual import block
type ImportDataPair struct {
	TerraformConfigLocation string

	RemoteCloudReference string
}

// HCLCreate is an interface that provides pre-built methods
// for generating and manipulating common HCL configuration.
type HCLCreate interface {
	// CreateMainTF outputs a bytes slice which defines a baseline main.tf file.
	CreateMainTF(providers map[string]string) ([]byte, error)

	// CreateImports creates either import blocks or tfmigrate configuration to import resources
	// into Terraform state.
	CreateImports(uniqueID string, workspaceToDirectory map[string]string) error

	// CreateTFMigrate coordinates CreateTFMigrateConfiguration and CreateTFMigrateMigration to create the needed
	// components for TFMigrate to operate successfully.
	CreateTFMigrate(uniqueID string, workspaceToDirectory map[string]string) error

	// CreateTFMigrateConfiguration saves HCL which defines TFMigrate configuration.
	CreateTFMigrateConfiguration(workspaceToDirectory map[string]string) error

	// CreateTFMigrateMigration saves HCL which defines a TFMigrate migration.
	CreateTFMigrateMigration(
		uniqueID string,
		resourceImportsByDivision ResourceImportsByDivision,
		newResourceToWorkspace NewResourceToWorkspace,
		workspaceToDirectory map[string]string,
	) error

	// ExtractResourceDefinitions outputs a bytes slice which defines needed Terraform resources extracted from
	// another configuration file.
	ExtractResourceDefinitions(noNewResources bool, workspaceToDirectory map[string]string) error

	// WriteImportBlocks writes import blocks to .tf files for configurations using Terraform version 1.5.0 or higher.
	WriteImportBlocks(uniqueID string, workspaceToDirectory map[string]string) error
}

// hclCreate implements the HCLCreate interface.
type hclCreate struct {
	// config comprises the configuration needed for hclCreate
	config Config

	// divisionToProvider is a mapping between a division and the provider that is responsible
	// for that division.
	provider terraformValueObjects.Provider `required:"true"`
}

// NewHCLCreate creates and returns a struct which implements the HCLCreate interface.
func NewHCLCreate(config Config, provider terraformValueObjects.Provider) (HCLCreate, error) {
	return &hclCreate{
		config:   config,
		provider: provider,
	}, nil
}

// Decode is a custom decoder of the MigrationHistoryDataMap for use with the envconfig library.
func (mhd *MigrationHistory) Decode(value string) error {
	var currentMap MigrationHistory

	if value == "" {
		return nil
	}

	err := json.Unmarshal([]byte(value), &currentMap)
	if err != nil {
		return fmt.Errorf("Error parsing specified json string: %v", err)
	}

	if currentMap.StorageType == "S3" {
		currentMap.StorageType = "s3"
	}
	if currentMap.StorageType == `Google Storage Bucket` {
		currentMap.StorageType = "gcs"
	}

	if currentMap.StorageType != "s3" && currentMap.StorageType != "gcs" {
		return fmt.Errorf("only types of 's3' and 'gcs' are currently supported. Attempted %v", currentMap.StorageType)
	}

	if currentMap.Region == "" {
		return fmt.Errorf("region variable cannot be empty")
	}

	if currentMap.Bucket == "" {
		return fmt.Errorf("the required field `bucket` is not present")
	}

	*mhd = currentMap
	return nil
}

// CreateImports creates either import blocks or tfmigrate configuration to import resources into Terraform state.
func (h *hclCreate) CreateImports(uniqueID string, workspaceToDirectory map[string]string) error {
	if h.config.TerraformVersion >= "1.5.0" {
		err := h.WriteImportBlocks(uniqueID, workspaceToDirectory)
		if err != nil {
			return fmt.Errorf("error creating import blocks: %v", err)
		}
	} else {
		err := h.CreateTFMigrate(uniqueID, workspaceToDirectory)
		if err != nil {
			return fmt.Errorf("error creating tfmigrate configuration: %v", err)
		}
	}
	return nil
}

// ConvertTerraformerResourceName takes an input resource name as output by the Terraformer
// package and outputs the resource in gold-standard Terraform name format.
func ConvertTerraformerResourceName(name string) string {
	intermediateString := strings.Replace(name, "tfer--", "", -1)

	finalResourceName := strings.Replace(intermediateString, "-", "_", -1)

	return finalResourceName
}
