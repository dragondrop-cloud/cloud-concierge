package documentize

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

// awsResourceDetails is a struct for packaging all relevant information for a AWS cloud resource.
type awsResourceDetails struct {

	// terraformName is the name of the aws resource within Terraform configuration.
	terraformName string

	// terraformType is the name of the aws resource type within Terraform configuration.
	terraformType string

	// terraformModule is the name of the module where the aws resource resides.
	terraformModule string

	// awsInstanceName is the name of the resource as it resides in AWS.
	awsInstanceName string

	// awsInstanceAccountID is the AWS Account ID where the resource resides.
	awsInstanceAccountID string

	// awsInstanceRegion is the region where the AWS resource resides.
	awsInstanceRegion string

	// awsInstanceTags are the tags on the resource
	awsInstanceTags map[string]string
}

// awsResourceExtractor implements the ResourceExtractor interface for
// aws cloud resources.
type awsResourceExtractor struct {

	// currentResourceDetails is a struct containing information about a resource necessary
	// for generating a document about it.
	currentResourceDetails *awsResourceDetails

	// typeToCategory is a map between aws resource types and their categories.
	typeToCategory TypeToCategory
}

// NewAWSResourceExtractor returns an instance of awsResourceExtractor.
func NewAWSResourceExtractor() ResourceExtractor {
	return &awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{},
		typeToCategory:         awsResourceCategories(),
	}
}

// GetCurrentResourceDetails returns the details for a awsResourceExtractor instance.
func (are *awsResourceExtractor) GetCurrentResourceDetails() *awsResourceDetails {
	return are.currentResourceDetails
}

// ExtractResourceDetails extracts relevant data points from a terraform state resource.
func (are *awsResourceExtractor) ExtractResourceDetails(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) error {

	attribute := "attributes"
	if isAttributesFlat {
		attribute += "_flat"
	}

	are.currentResourceDetails.terraformType = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "type").Data().(string)

	are.currentResourceDetails.terraformName = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "name").Data().(string)

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "module") {
		are.currentResourceDetails.terraformModule = tfStateParsed.Search("resources", strconv.Itoa(resourceIndex), "module").Data().(string)
	} else {
		are.currentResourceDetails.terraformModule = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "name") {
		are.currentResourceDetails.awsInstanceName = tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "name",
		).Data().(string)
	} else {
		are.currentResourceDetails.awsInstanceName = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "arn") {
		arn := tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "arn",
		).Data().(string)

		arnArray := strings.Split(arn, ":")

		are.currentResourceDetails.awsInstanceAccountID = arnArray[4]
		are.currentResourceDetails.awsInstanceRegion = arnArray[3]
	} else {
		are.currentResourceDetails.awsInstanceAccountID = "none"
		are.currentResourceDetails.awsInstanceRegion = "none"
	}

	if tfStateParsed.Exists("resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "tags_all") {
		tagsChildMap := tfStateParsed.Search(
			"resources", strconv.Itoa(resourceIndex), "instances", strconv.Itoa(instanceIndex), attribute, "tags_all",
		).ChildrenMap()

		tagMap := map[string]string{}
		for key, value := range tagsChildMap {
			tagMap[key] = value.Data().(string)
		}

		are.currentResourceDetails.awsInstanceTags = tagMap
	} else {
		are.currentResourceDetails.awsInstanceTags = map[string]string{}
	}

	return nil
}

// ResourceDetailsToSentence converts resource details to an english sentence format.
func (are *awsResourceExtractor) ResourceDetailsToSentence() string {
	outputSentence := fmt.Sprintf("terraform name of %v and type %v",
		stringToWords(are.currentResourceDetails.terraformName),
		stringToWords(are.currentResourceDetails.terraformType),
	)

	if are.currentResourceDetails.terraformModule != "none" {
		outputSentence = fmt.Sprintf("%v within module %v",
			outputSentence,
			stringToWords(are.currentResourceDetails.terraformModule),
		)
	}

	outputSentence = fmt.Sprintf("%v resource at location %v",
		outputSentence,
		are.currentResourceDetails.awsInstanceRegion,
	)

	if are.currentResourceDetails.awsInstanceName != "none" {
		outputSentence += fmt.Sprintf(" resource name of %v",
			stringToWords(are.currentResourceDetails.awsInstanceName),
		)
	}

	if are.currentResourceDetails.awsInstanceAccountID != "none" {
		outputSentence += fmt.Sprintf(" resource account of %v",
			stringToWords(are.currentResourceDetails.awsInstanceAccountID),
		)
	}

	if len(are.currentResourceDetails.awsInstanceTags) != 0 {
		for key, value := range are.currentResourceDetails.awsInstanceTags {
			outputSentence += fmt.Sprintf(
				" with tag key of %v and value of %v",
				stringToWords(key),
				stringToWords(value),
			)
		}
	}

	if catMap, ok := are.typeToCategory[ResourceType(are.currentResourceDetails.terraformType)]; ok {
		outputSentence += fmt.Sprintf(" with primary category of %v", catMap.primaryCat)

		if catMap.secondaryCat != "" {
			outputSentence += fmt.Sprintf(" and secondary category of %v", catMap.secondaryCat)
		}
	}

	return strings.ToLower(outputSentence) + ". "
}

// OutputResourceDetailsSentence coordinates ExtractResourceDetails and ResourceDetailsToSentence in order
// to extract and format as a sentence a resource's details from within a state file.
func (are *awsResourceExtractor) OutputResourceDetailsSentence(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) (string, error) {
	err := are.ExtractResourceDetails(tfStateParsed, isAttributesFlat, resourceIndex, instanceIndex)

	if err != nil {
		return "", fmt.Errorf("[are.ExtractResourceDetails] %v", err)
	}

	resourceSentenceDetails := are.ResourceDetailsToSentence()

	return resourceSentenceDetails, nil
}
