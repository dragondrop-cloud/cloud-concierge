package documentize

import "github.com/Jeffail/gabs/v2"

// ResourceExtractor defines an interface for extracting relevant details from
// a cloud resource.
type ResourceExtractor interface {
	// ExtractResourceDetails extracts relevant data points from a terraform state resource.
	ExtractResourceDetails(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) error

	// ResourceDetailsToSentence converts resource details to an english sentence format.
	ResourceDetailsToSentence() string

	// OutputResourceDetailsSentence coordinates ExtractResourceDetails and ResourceDetailsToSentence in order
	// to extract and format as a sentence a resource's details from within a state file.
	OutputResourceDetailsSentence(tfStateParsed *gabs.Container, isAttributesFlat bool, resourceIndex int, instanceIndex int) (string, error)
}
