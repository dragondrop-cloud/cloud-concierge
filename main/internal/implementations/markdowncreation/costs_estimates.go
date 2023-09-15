package markdowncreation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atsushinee/go-markdown-generator/doc"
	log "github.com/sirupsen/logrus"
)

// setCostsEstimatesData sets the costs estimates data in the markdown report
func (m *MarkdownCreator) setCostsEstimatesData(report *doc.MarkDownDoc) {
	report.Write("# Calculable Cloud Costs (Monthly)").Writeln().Writeln()

	if len(m.costEstimates) == 0 {
		report.Write("No cost estimates calculated for scanned resources.").Writeln()
		return
	}

	newResourcesCosts, notNewResourcesCosts := m.getNewAndNotNewResourcesCosts()
	log.Debugf("New resources costs: %s", newResourcesCosts)
	log.Debugf("Not new resources costs: %s", notNewResourcesCosts)

	report.Write("|Uncontrolled Resources Cost")
	report.Write("|Terraform Controlled Resources Cost|").Writeln()
	report.Write("| :---: | :---: |").Writeln()

	report.Write(fmt.Sprintf("|$%s", newResourcesCosts))
	report.Write(fmt.Sprintf("|$%s|", notNewResourcesCosts))

	report.Writeln().Writeln()
}

// getNewAndNotNewResourcesCosts returns the costs of new and not new resources
func (m *MarkdownCreator) getNewAndNotNewResourcesCosts() (string, string) {
	newResources := make(map[string]bool)
	for resourceID := range m.newResources {
		resourceIDComponents := strings.Split(resourceID, ".")
		resourceType := resourceIDComponents[0]
		resourceTerraformerName := resourceIDComponents[1]

		newResources[fmt.Sprintf("%s.%s", resourceType, resourceTerraformerName)] = true
	}

	newResourcesCosts := 0.0
	notNewResourcesCosts := 0.0
	for _, costEstimate := range m.costEstimates {
		totalCost := 0.0

		// If the resource has a definitive monthly cost from infracost, we use that value
		if costEstimate.MonthlyCost != "" && costEstimate.MonthlyCost != "0" {
			monthlyCostFloat, err := strconv.ParseFloat(costEstimate.MonthlyCost, 32)
			if err != nil {
				log.Debugf("[markdown_creator][get_new_and_not_new_resources_costs] error parsing monthly cost: %s", err)
				monthlyCostFloat = 0.0
			}

			totalCost = monthlyCostFloat
			// If no monthly cost, only add the unit price if the resource is not usage based
		} else if !costEstimate.IsUsageBased {
			priceCostFloat, err := strconv.ParseFloat(costEstimate.Price, 32)
			if err != nil {
				log.Debugf("[markdown_creator][get_new_and_not_new_resources_costs] error parsing price: %s", err)
				priceCostFloat = 0.0
			}

			totalCost = priceCostFloat
		}

		if newResources[costEstimate.ResourceName] {
			newResourcesCosts += totalCost
		} else {
			notNewResourcesCosts += totalCost
		}
	}

	return fmt.Sprintf("%.2f", newResourcesCosts), fmt.Sprintf("%.2f", notNewResourcesCosts)
}
