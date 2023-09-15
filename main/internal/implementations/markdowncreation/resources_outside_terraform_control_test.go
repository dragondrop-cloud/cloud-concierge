package markdowncreation

import (
	"strings"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.newResources = map[string]string{
		"aws_db_snapshot.example":      "terraform generated resource",
		"aws_db_subnet_group.example2": "terraform generated resource",
		"aws_lb_listener.example3":     "terraform generated resource",
	}

	// When
	markdownCreator.setResourcesOutsideOfTerraformControlData(report)

	// Then
	title := "# Resources Outside of Terraform Control\n\n"

	tableHeaders := "|Type|# Resources|\n| :---: | :---: |\n"

	reportResult := report.String()
	assert.Equal(t, title, reportResult[:len(title)])
	assert.Equal(t, tableHeaders, reportResult[len(title):len(title)+len(tableHeaders)])

	resourcesValues := strings.Split(reportResult[len(title)+len(tableHeaders):], "\n")
	assert.Equal(t, 5, len(resourcesValues))
	require.Contains(t, resourcesValues, "|aws_db_snapshot|1|")
	require.Contains(t, resourcesValues, "|aws_lb_listener|1|")
	require.Contains(t, resourcesValues, "|aws_db_subnet_group|1|")
}
