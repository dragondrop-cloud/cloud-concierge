package documentize

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

// ResourceData contains data on a Terraform resource type and the id for that resource
type ResourceData struct {
	// tfType is the Terraform resource's type
	tfType ResourceType

	// id is the Terraform resource's id
	id ResourceID

	// name is th the name of the Terraform resource within the Terraform configuration
	name ResourceName
}

// ConvertNewResourcesToJSON converts the output of NewResourceDocuments to a json-format byte array.
func (d *documentize) ConvertNewResourcesToJSON(resourceDocMap map[ResourceName]string) ([]byte, error) {
	jsonObj := gabs.New()

	for resourceName, doc := range resourceDocMap {
		_, err := jsonObj.Set(doc, string(resourceName))

		if err != nil {
			return nil, fmt.Errorf("[ConvertNewResourcesToJSON] error in jsonObj.Set() for resource %v: %v", resourceName, err)
		}
	}
	return jsonObj.Bytes(), nil
}

// IdentifyNewResources determines which resources in the remote cloud environment state files from
// terraformer are not present in the workspace state files. Returns a map of new resources to their
// corresponding provider.
func (d *documentize) IdentifyNewResources(workspaceToDirectory map[string]string) (map[ResourceData]bool, error) {
	workspaceToIDMap, err := d.pullWorkspaceResourceIdentifiers(workspaceToDirectory)
	if err != nil {
		return nil, fmt.Errorf("[d.pullWorkspaceResourceIdentifiers] %v", err)
	}

	resourceData, err := d.pullTerraformerResourceIdentifiers()
	if err != nil {
		return nil, fmt.Errorf("[d.pullTerraformerResourceIdentifiers] %v", err)
	}

	divToResourceData := selectNewResources(workspaceToIDMap, resourceData)

	return divToResourceData, nil
}

// selectNewResources determines which of the resources in divToID are not seen within typeToIDMap.
func selectNewResources(workspaceToIDMap map[Workspace]map[ResourceData]bool, resourceData map[ResourceData]bool) map[ResourceData]bool {
	newResourcesMap := map[ResourceData]bool{}
	resourceBlackList := newResourceBlackList()

	for tfrResourceData := range resourceData {
		if isValidNewResource(resourceBlackList, tfrResourceData, workspaceToIDMap) {
			newResourcesMap[tfrResourceData] = true
		}
	}

	return newResourcesMap
}

// isValidNewResource checks to see if the resource tfrResource is present in the typeToIDMap and is
// not a black-listed terraformer-generate resource.
func isValidNewResource(resourceBlackList map[string]bool, tfrResource ResourceData, typeToIDMap map[Workspace]map[ResourceData]bool) bool {
	for _, definedResourceSet := range typeToIDMap {
		for definedResource := range definedResourceSet {
			if (tfrResource.tfType == definedResource.tfType && tfrResource.id == definedResource.id) || resourceBlackList[string(tfrResource.tfType)] {
				return false
			}
		}
	}
	return true
}

// pullTerraformerResourceIdentifiers extracts identifiers for each unique resource instance within pulled terraformer generated state files.
func (d *documentize) pullTerraformerResourceIdentifiers() (map[ResourceData]bool, error) {
	tfStateBytes, err := os.ReadFile("current_cloud/terraform.tfstate")
	if err != nil {
		return nil, fmt.Errorf("[os.ReadFile] Error reading in state file: %v", err)
	}

	tfStateParsed, err := gabs.ParseJSON(tfStateBytes)
	if err != nil {
		return nil, fmt.Errorf("[gabs.ParseJSON] Error reading in state file: %v", err)
	}

	resourceDataSet, err := extractResourceIdsFromTerraformerState(tfStateParsed)

	if err != nil {
		return nil, fmt.Errorf("[extractResourceIdsFromWorkspaceState] Error in reading resource ids from workspace state: %v", err)
	}

	return resourceDataSet, nil
}

// extractResourceIdsFromWorkspaceState extracts identifying information for all resource instances
// within the current gabs-parsed terraformer-generated state json.
func extractResourceIdsFromTerraformerState(tfStateParsed *gabs.Container) (map[ResourceData]bool, error) {
	outputMap := map[ResourceData]bool{}

	i := 0
	for tfStateParsed.Exists("resources", strconv.Itoa(i)) {
		j := 0
		currentType := tfStateParsed.Search("resources", strconv.Itoa(i), "type").Data().(string)
		currentName := tfStateParsed.Search("resources", strconv.Itoa(i), "name").Data().(string)

		for tfStateParsed.Exists("resources", strconv.Itoa(i), "instances", strconv.Itoa(j), "attributes_flat", "id") {
			currentInstanceID := tfStateParsed.Search("resources", strconv.Itoa(i), "instances", strconv.Itoa(j), "attributes_flat", "id").Data().(string)

			currentResourceData := ResourceData{
				id:     ResourceID(currentInstanceID),
				name:   ResourceName(currentName),
				tfType: ResourceType(currentType),
			}

			outputMap[currentResourceData] = true
			j++
		}
		i++
	}
	return outputMap, nil
}

// pullWorkspaceResourceIdentifiers extracts identifiers for each unique resource instance within pulled
// workspace state files.
func (d *documentize) pullWorkspaceResourceIdentifiers(workspaceToDirectory map[string]string) (map[Workspace]map[ResourceData]bool, error) {
	outputMap := map[Workspace]map[ResourceData]bool{}

	for w := range workspaceToDirectory {
		tfStateBytes, err := os.ReadFile(fmt.Sprintf("state_files/%v.json", w))

		if err != nil {
			return nil, fmt.Errorf("[os.ReadFile] Error reading in state for workspace %v: %v", w, err)
		}

		tfStateParsed, err := gabs.ParseJSON(tfStateBytes)

		if err != nil {
			return nil, fmt.Errorf("[gabs.ParseJSON] Error parsing state for workspace %v: %v", w, err)
		}

		resourceTypeToID, err := extractResourceIdsFromWorkspaceState(tfStateParsed)

		if err != nil {
			return nil, fmt.Errorf("[extractResourceIdsFromWorkspaceState] %v", err)
		}

		outputMap[Workspace(w)] = resourceTypeToID
	}

	return outputMap, nil
}

// extractResourceIdsFromWorkspaceState extracts ids for all resource instances within the current
// gabs-parsed workspace state json.
func extractResourceIdsFromWorkspaceState(tfStateParsed *gabs.Container) (map[ResourceData]bool, error) {
	outputMap := map[ResourceData]bool{}

	i := 0
	for tfStateParsed.Exists("resources", strconv.Itoa(i)) {
		j := 0
		currentType := tfStateParsed.Search("resources", strconv.Itoa(i), "type").Data().(string)

		for tfStateParsed.Exists("resources", strconv.Itoa(i), "instances", strconv.Itoa(j), "attributes", "id") {
			currentInstanceID := tfStateParsed.Search("resources", strconv.Itoa(i), "instances", strconv.Itoa(j), "attributes", "id").Data().(string)
			outputMap[ResourceData{id: ResourceID(currentInstanceID), tfType: ResourceType(currentType)}] = true
			j++
		}
		i++
	}
	return outputMap, nil
}

// NewResourceDocuments creates a map between new resources and a document extracted from that
// resource definition.
func (d *documentize) NewResourceDocuments(resourceSet map[ResourceData]bool) (map[ResourceName]string, error) {
	outputMap := map[ResourceName]string{}

	tfrStateBytes, err := os.ReadFile("current_cloud/terraform.tfstate")
	if err != nil {
		return nil, fmt.Errorf("[os.ReadFile] Error reading in current_cloud/terraform.tfstate: %v", err)
	}

	tfrStateParsed, err := gabs.ParseJSON(tfrStateBytes)
	if err != nil {
		return nil, fmt.Errorf("[gabs.ParseJSON] Error parsing tfrStateBytes %v", err)
	}

	for resource := range resourceSet {
		resourceName, resourceDoc, err := d.pullResourceDocumentFromDiv(tfrStateParsed, resource)

		if err != nil {
			return nil, fmt.Errorf("[d.pullResourceDocumentFromDiv] Error: %v", err)
		}

		outputMap[resourceName] = resourceDoc
	}

	return outputMap, nil
}

// pullResourceDocumentFromDiv determines which resource from which extract the document definition of
// a particular resource identified by terraformer.
func (d *documentize) pullResourceDocumentFromDiv(tfrStateParsed *gabs.Container, resource ResourceData) (ResourceName, string, error) {
	i := 0
	for tfrStateParsed.Exists("resources", strconv.Itoa(i)) {
		if string(resource.tfType) != tfrStateParsed.Search("resources", strconv.Itoa(i), "type").Data().(string) {
			i++
			continue
		}
		currentMode := tfrStateParsed.Search("resources", strconv.Itoa(i), "mode").Data().(string)
		if currentMode != "managed" {
			i++
			continue
		}

		j := 0
		for tfrStateParsed.Exists("resources", strconv.Itoa(i), "instances", strconv.Itoa(j)) {
			currentID := tfrStateParsed.Search("resources", strconv.Itoa(i), "instances", strconv.Itoa(j), "attributes_flat", "id").Data().(string)

			if currentID == string(resource.id) {
				doc, err := d.extractResourceDocument(tfrStateParsed, true, i, j)
				if err != nil {
					return "", "", fmt.Errorf("[d.extractResourceDocument] Error: %v", err)
				}

				return ResourceName(fmt.Sprintf("%v.%v", resource.tfType, resource.name)), doc, nil
			}
			j++
		}
		i++
	}

	return "", "", fmt.Errorf("could not find resource %v, %v", resource.tfType, resource.name)
}

// extractResourceDocument extracts the document definition for a resource at a given resource and instance location
// within tfStateParsed.
func (d *documentize) extractResourceDocument(tfStateParsed *gabs.Container, isAttributesFlat bool, i int, j int) (string, error) {
	currentTFProvider := tfStateParsed.Search("resources", strconv.Itoa(i), "provider").Data().(string)

	currentTFProvider, err := regexProviderName(currentTFProvider)
	if err != nil {
		return "", fmt.Errorf("[regexProviderName] %v", err)
	}

	switch currentTFProvider {
	case "provider[\"registry.terraform.io/hashicorp/google\"]", "provider[\"registry.terraform.io/hashicorp/google-beta\"]":
		currentResourceSentence, err := d.resourceExtractors["google"].OutputResourceDetailsSentence(tfStateParsed, isAttributesFlat, i, j)
		if err != nil {
			return "", fmt.Errorf("[google.OutputResourceDetailsSentence] Error pulling details: %v", err)
		}
		return currentResourceSentence, nil

	case "provider[\"registry.terraform.io/hashicorp/aws\"]":
		currentResourceSentence, err := d.resourceExtractors["aws"].OutputResourceDetailsSentence(tfStateParsed, isAttributesFlat, i, j)
		if err != nil {
			return "", fmt.Errorf("[aws.OutputResourceDetailsSentence] Error pulling details: %v", err)
		}
		return currentResourceSentence, nil

	case "provider[\"registry.terraform.io/hashicorp/azurerm\"]":
		currentResourceSentence, err := d.resourceExtractors["azurerm"].OutputResourceDetailsSentence(tfStateParsed, isAttributesFlat, i, j)
		if err != nil {
			return "", fmt.Errorf("[aws.OutputResourceDetailsSentence] Error pulling details: %v", err)
		}
		return currentResourceSentence, nil

	default:
		resourceName := tfStateParsed.Search("resources", strconv.Itoa(i), "name").Data().(string)
		fmt.Printf("Currently unsupported provider %v, skipping the resource: %v", currentTFProvider, resourceName)
	}

	return "", nil
}
