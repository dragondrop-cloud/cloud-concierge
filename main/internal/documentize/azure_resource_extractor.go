package documentize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

// azureResourceDetails is a struct for packaging all relevant information for an Azure cloud resource.
type azureResourceDetails struct {
	// terraformName is the name of the azure resource within Terraform configuration.
	terraformName string

	// terraformType is the name of the azure resource type within Terraform configuration.
	terraformType string

	// terraformModule is the name of the module where the azure resource resides.
	terraformModule string

	// azureInstanceName is the name of the resource as it resides in AZURE.
	azureInstanceName string

	// azureInstanceLocation is the name of the region where the resource resides in AZURE.
	azureInstanceLocation string

	// azureInstanceTags are the tags on the resource
	azureInstanceTags map[string]string
}

// azureResourceExtractor implements the ResourceExtractor interface for
// azure cloud resources.
type azureResourceExtractor struct {
	// currentResourceDetails is a struct containing information about a resource necessary
	// for generating a document about it.
	currentResourceDetails *azureResourceDetails

	// typeToCategory is a map between azure resource types and their categories.
	typeToCategory TypeToCategory
}

// NewAzureResourceExtractor returns an instance of azureResourceExtractor.
func NewAzureResourceExtractor() ResourceExtractor {
	return &azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{},
		typeToCategory:         azureResourceCategories(),
	}
}

// GetCurrentResourceDetails returns the details for a azureResourceExtractor instance.
func (arx *azureResourceExtractor) GetCurrentResourceDetails() *azureResourceDetails {
	return arx.currentResourceDetails
}

// ExtractResourceDetails extracts relevant data points from a terraform state resource.
func (arx *azureResourceExtractor) ExtractResourceDetails(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) error {
	attribute := "attributes"
	if isAttributesFlat {
		attribute += "_flat"
	}

	resourcesArray := tfStateParsed.Path("resources").Data().([]interface{})
	if resourceIndex >= len(resourcesArray) {
		return fmt.Errorf("resourceIndex out of bounds")
	}

	resource := resourcesArray[resourceIndex].(map[string]interface{})
	instances := resource["instances"].([]interface{})
	if instanceIndex >= len(instances) {
		return fmt.Errorf("instanceIndex out of bounds")
	}

	instance := instances[instanceIndex].(map[string]interface{})
	attributes := instance[attribute].(map[string]interface{})

	arx.currentResourceDetails.terraformName = resource["name"].(string)
	arx.currentResourceDetails.terraformType = resource["type"].(string)

	arx.currentResourceDetails.terraformModule = "none"
	if instance["module"] != nil {
		arx.currentResourceDetails.terraformModule = instance["module"].(string)
	}

	arx.currentResourceDetails.azureInstanceName = resource["name"].(string)
	if instance["name"] != nil {
		arx.currentResourceDetails.azureInstanceName = attributes["name"].(string)
	}

	arx.currentResourceDetails.azureInstanceLocation = "global"
	if attributes["location"] != nil {
		arx.currentResourceDetails.azureInstanceLocation = attributes["location"].(string)
	}

	tags := make(map[string]string)
	if isAttributesFlat {
		for key, value := range attributes {
			if strings.HasPrefix(key, "tags.") && key != "tags.%" {
				tagKey := strings.TrimPrefix(key, "tags.")
				tags[tagKey] = value.(string)
			}
		}
	} else if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "tags") {
		tagsChildMap := tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "tags",
		).ChildrenMap()

		for key, value := range tagsChildMap {
			tags[key] = value.Data().(string)
		}

		arx.currentResourceDetails.azureInstanceTags = tags
	}

	arx.currentResourceDetails.azureInstanceTags = tags

	return nil
}

// ResourceDetailsToSentence converts resource details to an english sentence format.
func (arx *azureResourceExtractor) ResourceDetailsToSentence() string {
	// Base sentence structure
	sentence := fmt.Sprintf(
		"terraform name of %s and type %s",
		stringToWords(arx.currentResourceDetails.terraformName),
		stringToWords(arx.currentResourceDetails.terraformType),
	)

	if arx.currentResourceDetails.terraformModule != "none" {
		sentence = fmt.Sprintf(
			"%s within module %s",
			sentence,
			stringToWords(arx.currentResourceDetails.terraformModule),
		)
	}

	sentence = fmt.Sprintf("%v resource at location %v",
		sentence,
		stringToWords(arx.currentResourceDetails.azureInstanceLocation),
	)

	if arx.currentResourceDetails.azureInstanceName != "none" {
		sentence = fmt.Sprintf("%v resource name of %v",
			sentence,
			stringToWords(arx.currentResourceDetails.azureInstanceName),
		)
	}

	// Add tags if they exist
	if len(arx.currentResourceDetails.azureInstanceTags) > 0 {
		for key, value := range arx.currentResourceDetails.azureInstanceTags {
			sentence += fmt.Sprintf(" with tag key of %s and value of %s", key, value)
		}
	}

	// Add category info
	if category, ok := arx.typeToCategory[ResourceType(arx.currentResourceDetails.terraformType)]; ok {
		sentence += fmt.Sprintf(" with primary category of %s", category.primaryCat)
		if category.secondaryCat != "" {
			sentence += fmt.Sprintf(" and secondary category of %s", category.secondaryCat)
		}
	}

	// End the sentence
	sentence += "."

	return sentence
}

// OutputResourceDetailsSentence coordinates ExtractResourceDetails and ResourceDetailsToSentence in order
// to extract and format as a sentence a resource's details from within a state file.
func (arx *azureResourceExtractor) OutputResourceDetailsSentence(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) (string, error) {
	err := arx.ExtractResourceDetails(tfStateParsed, isAttributesFlat, resourceIndex, instanceIndex)
	if err != nil {
		return "", fmt.Errorf("[arx.ExtractResourceDetails] %v", err)
	}

	resourceSentenceDetails := arx.ResourceDetailsToSentence()

	return resourceSentenceDetails, nil
}
