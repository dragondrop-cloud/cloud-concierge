package costEstimation

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// InfracostResourceData is a struct for handling individual data values for a resource entry in
// infracost cost estimation output.
type InfracostResourceData struct {

	// resourceID is the terraformer resource name in Terraform configuration: {resource_type}.{resource_name}.
	resourceID string

	// monthlyCost is the estimated monthly cost for the resource strictly as defined.
	monthlyCost string

	// isPrimarilyUsageBased is a boolean representing whether a resource's price estimation is usage based.
	isPrimarilyUsageBased bool

	// costComponentList is a list of CostComponent structs for the resource.
	costComponentList []CostComponent

	// subResources is a list of sub-resources for the current resource
	subResources []SubResource
}

// SubResource is a billing mechanism that is owned by a primary resource
// (like network egress billing on a storage bucket)
type SubResource struct {

	// name is the name of the sub-resource
	name string

	// monthlyCost is the estimated monthly cost for the resource strictly as defined.
	monthlyCost string

	// isPrimarilyUsageBased is a boolean representing whether a resource's price estimation is usage based.
	isPrimarilyUsageBased bool

	// costComponentList is a list of CostComponent structs for the resource.
	costComponentList []CostComponent
}

// CostComponent is a struct for organizing cost output for a single CostComponent in Infracost.
type CostComponent struct {

	// name is the cost component's name
	name string

	// unit is the billable unit for the cost component
	unit string

	// unitPrice is the billable price per unit for the cost component
	unitPrice string

	// monthlyQuantity is the number of units seen in a month with the current configuration.
	monthlyQuantity string

	// monthlyCost is the monthly estimated cost for the resource
	monthlyCost string
}

// TODO: Major refactor needed here
// FormatAllCostEstimates processes infracost-generated cost estimation data into a more concise format for
// downstream usage for all cloud divisions.
func (ce *CostEstimator) FormatAllCostEstimates() error {
	for division := range ce.config.DivisionCloudCredentials {
		gabsJSONString, err := ce.FormatCostEstimate(division)
		if err != nil {
			return fmt.Errorf("[ce.FormatCostEstimate for division %v]%v", division, err)
		}

		filePath := "current_cloud/infracost-formatted.json"

		err = os.WriteFile(filePath, []byte(gabsJSONString), 0400)
		if err != nil {
			return fmt.Errorf("[os.WriteFile]%v", err)
		}
	}
	return nil
}

// FormatCostEstimate processes infracost-generated cost estimates.
func (ce *CostEstimator) FormatCostEstimate(division terraformValueObjects.Division) (string, error) {
	filePath := "current_cloud/infracost.json"

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("[os.ReadFile]%v", err)
	}

	resourceDataList, err := ce.ParseJSONToStruct(fileBytes)
	if err != nil {
		return "", fmt.Errorf("[ce.parseJSONGABSContainerToStruct]%v", err)
	}

	output, err := ce.StructToJSONString(resourceDataList)
	if err != nil {
		return "", fmt.Errorf("[ce.structToJSONString]%v", err)
	}

	return output, nil
}

// ParseJSONToStruct takes a JSON file byte list and converts to a list of InfracostResourceData.
func (ce *CostEstimator) ParseJSONToStruct(rawJSONBytes []byte) ([]InfracostResourceData, error) {
	parsedJSON, err := gabs.ParseJSON(rawJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("[gabs.ParseJSON]%v", err)
	}

	infracostDataList := []InfracostResourceData{}

	resourcesContainer := parsedJSON.Search("projects", "0", "breakdown", "resources")

	i := 0
	for resourcesContainer.Exists(strconv.Itoa(i)) {
		currentResource := resourcesContainer.Search(strconv.Itoa(i))
		resourceName := currentResource.Search("name").Data().(string)
		monthlyCost := ""
		isPrimarilyUsageBased := true

		if currentResource.Search("monthlyCost").Data() != nil {
			monthlyCost = currentResource.Search("monthlyCost").Data().(string)
			isPrimarilyUsageBased = false
		}

		costComponentList := extractCostComponentList(currentResource)

		subResourceList := []SubResource{}

		k := 0
		for currentResource.Exists("subresources", strconv.Itoa(k)) {
			subResourceContainer := currentResource.Search("subresources", strconv.Itoa(k))

			monthlyCostSubResource := ""
			if subResourceContainer.Search("monthlyCost").Data() != nil {
				monthlyCostSubResource = subResourceContainer.Search("monthlyCost").Data().(string)
			}

			isPrimarilyUsageBasedSubResource := true
			if monthlyCostSubResource != "" {
				isPrimarilyUsageBasedSubResource = false
			}

			subResource := SubResource{
				name:                  subResourceContainer.Search("name").Data().(string),
				monthlyCost:           monthlyCostSubResource,
				isPrimarilyUsageBased: isPrimarilyUsageBasedSubResource,
				costComponentList:     extractCostComponentList(subResourceContainer),
			}
			subResourceList = append(subResourceList, subResource)

			k++
		}

		currentInfracostData := InfracostResourceData{
			resourceID:            resourceName,
			monthlyCost:           monthlyCost,
			isPrimarilyUsageBased: isPrimarilyUsageBased,
			costComponentList:     costComponentList,
			subResources:          subResourceList,
		}
		infracostDataList = append(infracostDataList, currentInfracostData)

		i++
	}

	return infracostDataList, nil
}

// extractCostComponentList extracts into a []CostComponent data type all cost component data from
// a gabs container.
func extractCostComponentList(currentResource *gabs.Container) []CostComponent {
	j := 0
	costComponentList := []CostComponent{}
	for currentResource.Exists("costComponents", strconv.Itoa(j)) {
		componentContainer := currentResource.Search("costComponents", strconv.Itoa(j))

		monthlyQuantity := ""
		if componentContainer.Search("monthlyQuantity").Data() != nil {
			monthlyQuantity = componentContainer.Search("monthlyQuantity").Data().(string)
		}

		monthlyCost := ""
		if componentContainer.Search("monthlyCost").Data() != nil {
			monthlyCost = componentContainer.Search("monthlyCost").Data().(string)
		}

		costComponentList = append(costComponentList, CostComponent{
			name:            componentContainer.Search("name").Data().(string),
			unit:            componentContainer.Search("unit").Data().(string),
			unitPrice:       componentContainer.Search("price").Data().(string),
			monthlyQuantity: monthlyQuantity,
			monthlyCost:     monthlyCost,
		})
		j++
	}

	return costComponentList
}

// StructToJSONString converts structured data into a pandas-compatible json string via the gabs library.
func (ce *CostEstimator) StructToJSONString(inputData []InfracostResourceData) (string, error) {
	jsonObj := gabs.New()
	_, err := jsonObj.Array()
	if err != nil {
		return "", fmt.Errorf("[jsonObj.Array()]%v", err)
	}

	for _, resource := range inputData {
		for _, costComponent := range resource.costComponentList {
			currentCostComponentRow, err := buildCostComponentRow(
				resource.resourceID, resource.isPrimarilyUsageBased, "", costComponent,
			)
			if err != nil {
				return "", fmt.Errorf("[buildCostComponentRow]%v", err)
			}

			err = jsonObj.ArrayAppend(currentCostComponentRow)
			if err != nil {
				return "", fmt.Errorf("[jsonObj.ArrayAppend()]%v", err)
			}
		}

		for _, subResource := range resource.subResources {
			for _, subCostComponent := range subResource.costComponentList {
				currentCostComponentRow, err := buildCostComponentRow(
					resource.resourceID, subResource.isPrimarilyUsageBased, subResource.name, subCostComponent,
				)
				if err != nil {
					return "", fmt.Errorf("[buildCostComponentRow]%v", err)
				}

				err = jsonObj.ArrayAppend(currentCostComponentRow)
				if err != nil {
					return "", fmt.Errorf("[jsonObj.ArrayAppend()]%v", err)
				}
			}
		}
	}

	return jsonObj.String(), nil
}

// buildCostComponentRow builds a unified row of cost component data for later export to json using GABS.
func buildCostComponentRow(resourceID string, isPrimarilyUsageBased bool, subResourceName string, costComponent CostComponent) (*gabs.Container, error) {
	currentCostComponentRow := gabs.New()
	_, err := currentCostComponentRow.Set(resourceID, "resource_name")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(costComponent.name, "cost_component")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(costComponent.unit, "unit")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(costComponent.unitPrice, "price")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(costComponent.monthlyQuantity, "monthly_quantity")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(costComponent.monthlyCost, "monthly_cost")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(isPrimarilyUsageBased, "is_usage_based")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	_, err = currentCostComponentRow.Set(subResourceName, "sub_resource_name")
	if err != nil {
		return nil, fmt.Errorf("[currentCostComponentRow.Set()]%v", err)
	}

	return currentCostComponentRow, nil
}
