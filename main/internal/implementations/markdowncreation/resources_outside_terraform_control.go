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
