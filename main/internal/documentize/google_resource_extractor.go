package documentize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

// googleResourceDetails is a struct for packaging all relevant information for a GCP cloud resource.
type googleResourceDetails struct {

	// terraformName is the name of the google resource within Terraform configuration.
	terraformName string

	// terraformType is the name of the google resource type within Terraform configuration.
	terraformType string

	// terraformModule is the name of the module where the google resource resides.
	terraformModule string

	// gcpInstanceName is the name of the resource as it resides in GCP.
	gcpInstanceName string

	// gcpInstanceProject is the name of the project where the resource resides in GCP.
	gcpInstanceProject string

	// gcpInstanceLocation is the name of the region where the resource resides in GCP.
	gcpInstanceLocation string
}

// googleResourceExtractor implements the ResourceExtractor interface for
// google cloud resources.
type googleResourceExtractor struct {

	// currentResourceDetails is a struct containing information about a resource necessary
	// for generating a document about it.
	currentResourceDetails *googleResourceDetails

	// typeToCategory is a map between google resource types and their categories.
	typeToCategory TypeToCategory
}

// NewGoogleResourceExtractor returns an instance of googleResourceExtractor.
func NewGoogleResourceExtractor() ResourceExtractor {
	return &googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{},
		typeToCategory:         googleResourceCategories(),
	}
}

// GetCurrentResourceDetails returns the details for a googleResourceExtractor instance.
func (gre *googleResourceExtractor) GetCurrentResourceDetails() *googleResourceDetails {
	return gre.currentResourceDetails
}

// ExtractResourceDetails extracts relevant data points from a terraform state resource.
func (gre *googleResourceExtractor) ExtractResourceDetails(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) error {

	attribute := "attributes"
	if isAttributesFlat {
		attribute += "_flat"
	}

	gre.currentResourceDetails.terraformType = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "type").Data().(string)

	gre.currentResourceDetails.terraformName = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "name").Data().(string)

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "module") {
		gre.currentResourceDetails.terraformModule = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "module").Data().(string)
	} else {
		gre.currentResourceDetails.terraformModule = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "project") {
		gre.currentResourceDetails.gcpInstanceProject = tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "project",
		).Data().(string)
	} else {
		gre.currentResourceDetails.gcpInstanceProject = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "name") {
		gre.currentResourceDetails.gcpInstanceName = tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "name",
		).Data().(string)
	} else {
		gre.currentResourceDetails.gcpInstanceName = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "location") {
		gre.currentResourceDetails.gcpInstanceLocation = tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "location",
		).Data().(string)
	} else {
		gre.currentResourceDetails.gcpInstanceLocation = "global"
	}

	return nil
}

// ResourceDetailsToSentence converts resource details to an english sentence format.
func (gre *googleResourceExtractor) ResourceDetailsToSentence() string {
	outputSentence := fmt.Sprintf("terraform name of %v and type %v",
		stringToWords(gre.currentResourceDetails.terraformName),
		stringToWords(gre.currentResourceDetails.terraformType),
	)

	if gre.currentResourceDetails.terraformModule != "none" {
		outputSentence = fmt.Sprintf("%v within module %v",
			outputSentence,
			stringToWords(gre.currentResourceDetails.terraformModule),
		)
	}

	outputSentence = fmt.Sprintf("%v resource at location %v",
		outputSentence,
		gre.currentResourceDetails.gcpInstanceLocation,
	)

	if gre.currentResourceDetails.gcpInstanceName != "none" {
		outputSentence += fmt.Sprintf(" resource name of %v",
			stringToWords(gre.currentResourceDetails.gcpInstanceName),
		)
	}

	if gre.currentResourceDetails.gcpInstanceProject != "none" {
		outputSentence += fmt.Sprintf(" resource project of %v",
			stringToWords(gre.currentResourceDetails.gcpInstanceProject),
		)
	}

	if catMap, ok := gre.typeToCategory[ResourceType(gre.currentResourceDetails.terraformType)]; ok {
		outputSentence += fmt.Sprintf(" with primary category of %v", catMap.primaryCat)

		if catMap.secondaryCat != "" {
			outputSentence += fmt.Sprintf(" and secondary category of %v", catMap.secondaryCat)
		}
	}

	return strings.ToLower(outputSentence) + ". "
}

// OutputResourceDetailsSentence coordinates ExtractResourceDetails and ResourceDetailsToSentence in order
// to extract and format as a sentence a resource's details from within a state file.
func (gre *googleResourceExtractor) OutputResourceDetailsSentence(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) (string, error) {
	err := gre.ExtractResourceDetails(tfStateParsed, isAttributesFlat, resourceIndex, instanceIndex)

	if err != nil {
		return "", fmt.Errorf("[gre.ExtractResourceDetails] %v", err)
	}

	resourceSentenceDetails := gre.ResourceDetailsToSentence()

	return resourceSentenceDetails, nil
}
