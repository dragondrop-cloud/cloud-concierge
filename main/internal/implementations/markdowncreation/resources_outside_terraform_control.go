package markdowncreation

import (
	"fmt"
	"strings"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// setResourcesOutsideOfTerraformControlData sets the resources outside terraform control data in the markdown report
func (m *MarkdownCreator) setResourcesOutsideOfTerraformControlData(report *doc.MarkDownDoc) {
	report.Write("# Resources Outside of Terraform Control").Writeln().Writeln()

	if len(m.newResources) == 0 {
		report.Write("No new resources found!").Writeln()
		return
	}

	if len(m.costEstimates) == 0 {
		m.resourcesWithoutCostEstimates(report)
		return
	}

	resourcesDetailsByType := m.getCostsEstimations()
	if len(resourcesDetailsByType) == 0 || !m.atLeastOneValidResource(resourcesDetailsByType) {
		m.resourcesWithoutCostEstimates(report)
		return
	}

	m.resourcesWithCostEstimates(report, resourcesDetailsByType)
}

// ResourceCostEstimate represents the cost estimate for a resource
type ResourceCostEstimate struct {
	ResourceCount  int
	CostComponents map[string]bool
	IsUsageBased   bool
	MonthlyCost    float64
}

func (m *MarkdownCreator) getCostsEstimations() map[string]*ResourceCostEstimate {
	resourcesDetailsByType := make(map[string]*ResourceCostEstimate)
	for resourceID := range m.newResources {
		resourceType := strings.Split(resourceID, ".")[0]
		if _, ok := resourcesDetailsByType[resourceType]; !ok {
			resourcesDetailsByType[resourceType] = &ResourceCostEstimate{}
		}

		resourcesDetailsByType[resourceType].ResourceCount++
	}

	for _, costEstimate := range m.costEstimates {
		resourceType := strings.Split(costEstimate.ResourceName, ".")[0]

		if _, ok := resourcesDetailsByType[resourceType]; ok {
			if resourcesDetailsByType[resourceType].CostComponents == nil {
				resourcesDetailsByType[resourceType].CostComponents = make(map[string]bool)
			}

			resourcesDetailsByType[resourceType].CostComponents[costEstimate.CostComponent] = true
			resourcesDetailsByType[resourceType].IsUsageBased = costEstimate.IsUsageBased

			resourcesDetailsByType[resourceType].MonthlyCost += m.getActualMonthlyCostEstimate(costEstimate)
		}
	}

	return resourcesDetailsByType
}

func (m *MarkdownCreator) atLeastOneValidResource(resourcesByType map[string]*ResourceCostEstimate) bool {
	validResourcesByType := make(map[string]*ResourceCostEstimate)

	for resourceType, resourceDetail := range resourcesByType {
		if len(resourceDetail.CostComponents) > 0 && ((resourceDetail.MonthlyCost > 0) || (resourceDetail.MonthlyCost == 0 && resourceDetail.IsUsageBased)) {
			validResourcesByType[resourceType] = resourceDetail
		}
	}

	return len(validResourcesByType) > 0
}

// resourcesWithCostEstimates sets the resources outside terraform control data in the markdown report
func (m *MarkdownCreator) resourcesWithCostEstimates(report *doc.MarkDownDoc, resourcesDetailsByType map[string]*ResourceCostEstimate) {
	report.Write("|Type|# Resources|Cost Components|Monthly Cost|Usage Based*|").Writeln()
	report.Write("| :---: | :---: | :---: | :---: | :---: |").Writeln()

	for resourceType, resourceDetail := range resourcesDetailsByType {
		report.Write(fmt.Sprintf("|%s", resourceType))
		report.Write(fmt.Sprintf("|%d", resourceDetail.ResourceCount))

		if len(resourceDetail.CostComponents) > 0 {
			report.Write(fmt.Sprintf("|%d", len(resourceDetail.CostComponents)))
			if resourceDetail.IsUsageBased {
				report.Write("|$0.00*")
				report.Write("|True|")
			} else {
				report.Write(fmt.Sprintf("|$%.2f", resourceDetail.MonthlyCost))
				report.Write("|False|")
			}
		} else {
			report.Write("|No Charge")
			report.Write("|No Charge")
			report.Write("|No Charge|")
		}

		report.Writeln()
	}

	report.Writeln()
}

// resourcesWithoutCostEstimates sets the resources outside terraform control data in the markdown report
func (m *MarkdownCreator) resourcesWithoutCostEstimates(report *doc.MarkDownDoc) {
	resourcesCounterByType := make(map[string]int)
	for resourceID := range m.newResources {
		resourceType := strings.Split(resourceID, ".")[0]
		resourcesCounterByType[resourceType]++
	}

	report.Write("|Type")
	report.Write("|# Resources|").Writeln()
	report.Write("| :---: | :---: |").Writeln()

	rowCounter := 0
	for resourceType, resourceCount := range resourcesCounterByType {
		report.Write(fmt.Sprintf("|%s", resourceType))
		report.Write(fmt.Sprintf("|%d|", resourceCount)).Writeln()
		rowCounter++
	}

	report.Writeln()
}
