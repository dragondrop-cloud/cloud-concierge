package driftdetector

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// AttributeDifference contains data on the specific differences between a Cloud resource
// and its Terraform representation.
type AttributeDifference struct {
	RecentActor           terraformValueObjects.CloudActor
	RecentActionTimestamp terraformValueObjects.Timestamp
	AttributeName         string
	TerraformValue        string
	CloudValue            string
	InstanceID            string
	InstanceRegion        string
	AttributeDetail
}

// AttributeDetail contains information about the resource that the attribute belongs to.
type AttributeDetail struct {
	StateFileName StateFileName
	ModuleName    string
	ResourceType  string
	ResourceName  string
}

// identifyResourceDifferences compares TerraformerResources and RemoteResources,
// and returns a slice of AttributeDifference with any differences found between the two maps of resources.
// It also returns an error if there's any issue during the comparison process.
func (m *ManagedResourcesDriftDetector) identifyResourceDifferences(
	terraformerResources TerraformerResourceIDToData,
	terraformResources TerraformStateResourceIDToData,
) ([]AttributeDifference, error) {
	attributeDifferences := make([]AttributeDifference, 0)

	for id, data := range terraformResources {

		if terraformerResource, ok := terraformerResources[id]; ok {
			terraformInstanceConverted, err := convertNestedMapToFlatAttributes(data.Attributes)
			if err != nil {
				return nil, fmt.Errorf("[convertNestedMapToFlatAttributes]%v", err)
			}

			attributeComplement := &AttributeDetail{
				StateFileName: StateFileName(data.StateFile),
				ModuleName:    data.Module,
				ResourceType:  data.Type,
				ResourceName:  data.Name,
			}

			driftedResources, resourcesChanged, err := compareFlatAttributesAndGetDrifted(terraformInstanceConverted, terraformerResource.AttributesFlat, attributeComplement)
			if err != nil {
				return nil, fmt.Errorf("[compareFlatAttributesAndGetDrifted]%v", err)
			}

			if resourcesChanged {
				attributeDifferences = append(attributeDifferences, driftedResources...)
			}
		}
	}

	return attributeDifferences, nil
}

// compareFlatAttributesAndGetDrifted compares the attributes of remoteResourceAttributes and terraformerAttributes,
// and returns a slice of AttributeDifference with any differences found between the two attribute maps.
// It also returns a boolean value indicating if any differences were found.
func compareFlatAttributesAndGetDrifted(terraformResourceAttributes map[string]string, terraformerAttributes map[string]string, complement *AttributeDetail) ([]AttributeDifference, bool, error) {
	var differences []AttributeDifference
	resourcesChanged := false

	cloudProvider := strings.Split(complement.ResourceType, "_")[0]
	region, err := ParseRegionFromTfStateMap(terraformerAttributes, cloudProvider)
	if err != nil {
		return nil, true, fmt.Errorf("[parseRegionFromTfStateMap]%v", err)
	}

	id, err := ResourceIDCalculator(terraformerAttributes, cloudProvider, complement.ResourceType)
	if err != nil {
		return nil, true, fmt.Errorf("[resourcesCalculator.ResourceIDCalculator]%v", err)
	}

	// case where the cloud representation of the resource has an attribute that is different from terraform
	for attribute, value := range terraformerAttributes {
		attributePathSeparated := strings.Split(attribute, ".")
		attributeName := attributePathSeparated[len(attributePathSeparated)-1]
		if attributeName == "#" || attributeName == "%" {
			continue
		}

		terraformValue, ok := terraformResourceAttributes[attribute]
		if !ok || terraformValue != value {
			resourcesChanged = true
			differences = append(differences, AttributeDifference{
				AttributeName:   attribute,
				TerraformValue:  terraformValue,
				CloudValue:      value,
				InstanceID:      id,
				InstanceRegion:  region,
				AttributeDetail: *complement,
			})
		}
	}

	// case where terraform has an attribute that the cloud representation of the resource does not
	for attribute, terraformValue := range terraformResourceAttributes {
		if _, ok := terraformerAttributes[attribute]; !ok && !strings.ContainsAny(attribute, "#%") {
			resourcesChanged = true

			differences = append(differences, AttributeDifference{
				AttributeName:   attribute,
				TerraformValue:  terraformValue,
				CloudValue:      "",
				InstanceID:      id,
				InstanceRegion:  region,
				AttributeDetail: *complement,
			})
		}
	}

	return differences, resourcesChanged, nil
}

// ParseRegionFromTfStateMap extracts the region value from a terraform state file attributes map.
func ParseRegionFromTfStateMap(attributes map[string]string, cloudProvider string) (string, error) {
	switch cloudProvider {
	case "aws":
		return extractRegionFromAWSAttributes(attributes)
	case "azurerm":
		return "", nil
	case "google":
		return attributes["location"], nil
	default:
		return "", fmt.Errorf("unknown cloud provider: %s", cloudProvider)
	}
}

// extractRegionFromAWSAttributes extracts the region value from a map of AWS attributes.
// The value is either in the region attribute or extracted from the arn of the resource.
func extractRegionFromAWSAttributes(attributes map[string]string) (string, error) {
	if region, ok := attributes["region"]; ok {
		return region, nil
	} else if arn, ok := attributes["arn"]; ok {
		arnSplit := strings.Split(arn, ":")
		arnRegion := arnSplit[3]
		if arnRegion != "" {
			return arnRegion, nil
		}
	}
	return "us-east-1", nil
}

// stringMapsEqual compares two maps with string keys and string values for equality.
// It returns true if both maps have the same length and corresponding key-value pairs,
// and false otherwise.
func stringMapsEqual(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if v != m2[k] {
			return false
		}
	}
	return true
}

// stringSlicesEqual compares two string slices for equality.
// It returns true if both slices have the same length and sorted elements,
// and false otherwise.
func stringSlicesEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// resourcesEqual compares two TerraformerResource pointers for equality.
// It returns true if the resources have the same Type, Name, Provider, Module,
// number of Instances, AttributesFlat, and SensitiveAttributes,
// and false otherwise.
func resourcesEqual(r1, r2 *TerraformerResource) bool {
	if !(r1.Type == r2.Type && r1.Name == r2.Name && r1.Provider == r2.Provider && r1.Module == r2.Module && len(r1.Instances) == len(r2.Instances)) {
		return false
	}

	for index, instance := range r1.Instances {
		if !stringMapsEqual(instance.AttributesFlat, r2.Instances[index].AttributesFlat) {
			return false
		}
	}

	return true
}

// convertNestedMapToFlatAttributes converts a dynamic nested map to a map of flat attributes as found in terraformer
// state file outputs.
func convertNestedMapToFlatAttributes(nestedMap map[string]interface{}) (map[string]string, error) {
	output := map[string]string{}

	err := recursiveToFlatAttributes(output, "", true, nil, nestedMap)
	if err != nil {
		return nil, fmt.Errorf("[recursiveToFlatAttributes]%v", err)
	}

	return output, nil
}

// recursiveToFlatAttributes is a recursive helper function used to add values to the output map of strings.
// Unlike most recursive functions, there is no "base case" return for the function unless an error is thrown at one stage.
func recursiveToFlatAttributes(
	output map[string]string, currentBase string, isMap bool,
	currentSlice []interface{}, currentMap map[string]interface{},
) error {
	if currentBase != "" {
		currentBase += "."
	}
	if isMap {
		for key, value := range currentMap {
			if value != nil {
				switch t := value.(type) {
				case float32:
					output[currentBase+key] = strconv.Itoa(int(value.(float32)))
				case float64:
					output[currentBase+key] = strconv.Itoa(int(value.(float64)))
				case bool:
					output[currentBase+key] = strconv.FormatBool(value.(bool))
				case int:
					output[currentBase+key] = strconv.Itoa(value.(int))
				case string:
					output[currentBase+key] = value.(string)
				case []interface{}:
					err := recursiveToFlatAttributes(output, currentBase+key, false, value.([]interface{}), nil)
					if err != nil {
						return err
					}
				case map[string]interface{}:
					err := recursiveToFlatAttributes(output, currentBase+key, true, nil, value.(map[string]interface{}))
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("unhandled type of %v", t)
				}
			}
		}
	} else {
		for i, value := range currentSlice {
			switch t := value.(type) {
			case string:
				output[currentBase+strconv.Itoa(i)] = value.(string)
			case []interface{}:
				err := recursiveToFlatAttributes(output, currentBase+strconv.Itoa(i), false, value.([]interface{}), nil)
				if err != nil {
					return err
				}
			case map[string]interface{}:
				err := recursiveToFlatAttributes(output, currentBase+strconv.Itoa(i), true, nil, value.(map[string]interface{}))
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unhandled type of %v", t)
			}
		}
	}
	return nil
}
