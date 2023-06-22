package documentize

import terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"

// Directory is the path to a Workspace's terraform configuration within a code repository.
type Directory string

// ResourceCategory is a struct containing information on a resource category
type ResourceCategory struct {

	// primaryCat is the resource's primary category.
	primaryCat string

	// secondaryCat is the resource's secondary category.
	secondaryCat string
}

// ResourceName is the Terraform resource name.
type ResourceName string

// ResourceID is the ID of the Terraform resource.
type ResourceID string

// ResourceType is the type of the Terraform resource.
type ResourceType string

// TypeToCategory is a map between resource type and resource categorization.
type TypeToCategory map[ResourceType]ResourceCategory

// Workspace is the name of a Terraform Cloud workspace.
type Workspace string

// Documentize interface for making "documents" out of Terraform state files and
// performing light comparisons between state files
type Documentize interface {

	// AllWorkspaceStatesToDocuments converts all workspace states to documents of non-sensitive strings.
	AllWorkspaceStatesToDocuments(workspaceToDirectory map[string]string) (map[Workspace][]byte, error)

	// ConvertWorkspaceDocumentsToJSON converts the output of AllWorkspaceStatesToDocuments to a json-format byte array.
	ConvertWorkspaceDocumentsToJSON(workspaceDocMap map[Workspace][]byte) ([]byte, error)

	// ConvertNewResourcesToJSON converts the output of NewResourceDocuments to a json-format byte array.
	ConvertNewResourcesToJSON(resourceDocMap map[ResourceName]string) ([]byte, error)

	// IdentifyNewResources determines which resources in the remote cloud environment state files from
	// terraformer are not present in the workspace state files. Returns a map of new resources to their
	// corresponding provider.
	IdentifyNewResources(workspaceToDirectory map[string]string) (map[terraformValueObjects.Division]map[ResourceData]bool, error)

	// NewResourceDocuments creates a map between new resources and a document extracted from that
	// resource definition.
	NewResourceDocuments(divisionToResource map[terraformValueObjects.Division]map[ResourceData]bool) (map[ResourceName]string, error)

	// WorkspaceStateToDocument converts a workspace state to a document of non-sensitive strings.
	WorkspaceStateToDocument(workspace Workspace) ([]byte, error)
}

// documentize is a struct that implements the Documentize interface.
type documentize struct {
	// DivisionToProvider is a map between the division name and the provider name
	divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider

	// resourceExtractors is a map between a provider name and the logic needed to extract
	// resource information for the provider.
	resourceExtractors map[terraformValueObjects.Provider]ResourceExtractor
}

// NewDocumentize creates a new instance that implements the Documentize interface.
func NewDocumentize(divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) (Documentize, error) {

	resourceExtractors := map[terraformValueObjects.Provider]ResourceExtractor{
		"aws":     NewAWSResourceExtractor(),
		"google":  NewGoogleResourceExtractor(),
		"azurerm": NewAzureResourceExtractor(),
	}

	return &documentize{
		divisionToProvider: divisionToProvider,
		resourceExtractors: resourceExtractors,
	}, nil
}
