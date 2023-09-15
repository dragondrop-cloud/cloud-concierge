package markdowncreation

import (
	"fmt"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// setDeletedResourcesData sets the deleted resources data in the markdown report
func (m *MarkdownCreator) setDeletedResourcesData(report *doc.MarkDownDoc) {
	report.Write("# Drifted Resources Deleted").Writeln().Writeln()

	if len(m.deletedResources) == 0 {
		report.Write("No deleted resources found!").Writeln()
		return
	}

	report.Write("|Type|Name|Module|State File|\n| :---: | :---: | :---: | :---: |\n")
	for _, deletedResource := range m.deletedResources {
		report.Write(fmt.Sprintf("|%s", deletedResource.ResourceType))
		report.Write(fmt.Sprintf("|%s", deletedResource.ResourceName))
		report.Write(fmt.Sprintf("|%s", deletedResource.ModuleName))
		report.Write(fmt.Sprintf("|%s|", deletedResource.StateFileName)).Writeln()
	}

	report.Writeln()
}
