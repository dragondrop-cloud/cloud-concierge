package hclcreate

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// costs is a map between the complete resource name and the corresponding
// resourceCosts struct.
type costs map[terraformValueObjects.ResourceName]resourceCosts

// resourceCosts is a struct containing cost descriptions for resource costs from
// costComponents and subResources.
type resourceCosts struct {
	// costComponents is a list of individual cost components for a resource.
	costComponents []costComponent

	// subResources is a list of subResources for a resource.
	subResources map[string][]costComponent
}

// costComponent struct for organizing information on a resource's single
// component cost.
type costComponent struct {
	// componentName is the name of the component.
	componentName string

	// isUsageBased is whether the component's price is usage based
	isUsageBased bool

	// monthlyCost is the total cost for the component.
	monthlyCost string

	// price is the unit price for the costComponent
	price string

	// unit is the name of unit that the costComponent is priced in.
	unit string
}

// gabsContainerToCostsStruct converts a gabs container to the divisionCosts struct.
func gabsContainerToCostsStruct(c *gabs.Container) (costs, error) {
	divCosts := costs{}

	for _, cost := range c.Children() {
		resourceName := cost.Search("resource_name").Data().(string)
		intermediateString := strings.Replace(resourceName, "tfer--", "", -1)
		resourceName = strings.Replace(intermediateString, "-", "_", -1)

		var resourceCostEntry resourceCosts
		existingResourceCost, ok := divCosts[terraformValueObjects.ResourceName(resourceName)]
		if ok {
			resourceCostEntry = existingResourceCost
		} else {
			resourceCostEntry = resourceCosts{
				costComponents: []costComponent{},
				subResources:   map[string][]costComponent{},
			}
		}

		subResourceName := cost.Search("sub_resource_name").Data().(string)
		if subResourceName == "" {
			// cost component definition
			currentComponent := costComponent{
				componentName: cost.Search("cost_component").Data().(string),
				isUsageBased:  cost.Search("is_usage_based").Data().(bool),
				monthlyCost:   cost.Search("monthly_cost").Data().(string),
				price:         cost.Search("price").Data().(string),
				unit:          cost.Search("unit").Data().(string),
			}
			resourceCostEntry.costComponents = append(resourceCostEntry.costComponents, currentComponent)
		} else {
			existingSubComponentList, ok := resourceCostEntry.subResources[subResourceName]
			if !ok {
				existingSubComponentList = []costComponent{}
			}

			// sub resource cost component definition
			currentComponent := costComponent{
				componentName: cost.Search("cost_component").Data().(string),
				isUsageBased:  cost.Search("is_usage_based").Data().(bool),
				monthlyCost:   cost.Search("monthly_cost").Data().(string),
				price:         cost.Search("price").Data().(string),
				unit:          cost.Search("unit").Data().(string),
			}
			existingSubComponentList = append(existingSubComponentList, currentComponent)
			resourceCostEntry.subResources[subResourceName] = existingSubComponentList
		}

		divCosts[terraformValueObjects.ResourceName(resourceName)] = resourceCostEntry
	}

	return divCosts, nil
}

// generateHCLCloudCostComment generates data on Cloud Actor actions for the specified resource.
func (h *hclCreate) generateHCLCloudCostComment(
	resourceType string, resourceName string,
	costEstimates costs,
) hclwrite.Tokens {
	completeResourceName := fmt.Sprintf("%v.%v", resourceType, resourceName)
	cloudCostStatement := ""
	cloudCostCurrentResource, ok := costEstimates[terraformValueObjects.ResourceName(completeResourceName)]
	if ok {
		cloudCostStatement += "# Identified Resource Cost Components:\n"
		for _, costComponent := range cloudCostCurrentResource.costComponents {

			monthlyCost := "0"
			if costComponent.monthlyCost != "" {
				monthlyCost = costComponent.monthlyCost
			}

			if costComponent.isUsageBased {
				cloudCostStatement += fmt.Sprintf("## %v (Usage-based)\n###### Price / Unit: $%v / %v\n", costComponent.componentName, costComponent.price, costComponent.unit)
			} else {
				cloudCostStatement += fmt.Sprintf("## %v\n###### Monthly Cost: $%v\n###### Price / Unit: $%v / %v\n", costComponent.componentName, monthlyCost, costComponent.price, costComponent.unit)
			}
		}

		for name, subResource := range cloudCostCurrentResource.subResources {
			cloudCostStatement += fmt.Sprintf("\n# This resource has additional cost components for %v:\n", name)
			for _, costComponent := range subResource {

				monthlyCost := "0"
				if costComponent.monthlyCost != "" {
					monthlyCost = costComponent.monthlyCost
				}

				if costComponent.isUsageBased {
					cloudCostStatement += fmt.Sprintf("## %v (Usage-based)\n###### Price / Unit: $%v / %v \n", costComponent.componentName, costComponent.price, costComponent.unit)
				} else {
					cloudCostStatement += fmt.Sprintf("## %v\n###### Monthly Cost: $%v\n###### Price / Unit: $%v / %v\n", costComponent.componentName, monthlyCost, costComponent.price, costComponent.unit)
				}
			}
		}
	} else {
		cloudCostStatement = "# This resource has no identified cost"
	}
	return hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(cloudCostStatement),
			SpacesBefore: 0,
		},
	}
}
