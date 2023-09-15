package markdowncreation

import (
	"fmt"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownCreator_setCostsEstimatesData_ManagedResources(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.costEstimates = []CostEstimate{
		{
			MonthlyCost:  "16.425",
			ResourceName: "aws_instance.example",
		},
		{
			Price:        "5.84",
			ResourceName: "aws_instance.example2",
		},
	}
	markdownCreator.newResources = map[string]string{
		"aws_instance2.example":  "terraform generated resource",
		"aws_instance2.example2": "terraform generated resource",
	}

	// When
	markdownCreator.setCostsEstimatesData(report)

	// Then
	title := "# Calculable Cloud Costs (Monthly)\n\n"

	tableHeaders := "|Uncontrolled Resources Cost|Terraform Controlled Resources Cost|\n| :---: | :---: |\n"
	tableContent := "|$0.00|$22.26|\n\n"

	expectedMarkdown := fmt.Sprintf("%s%s%s", title, tableHeaders, tableContent)
	assert.Equal(t, expectedMarkdown, report.String())
}

func TestMarkdownCreator_setCostsEstimatesData_NewResources(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.costEstimates = []CostEstimate{
		{
			MonthlyCost:  "16.425",
			ResourceName: "aws_instance.example",
		},
		{
			Price:        "5.84",
			ResourceName: "aws_instance.example2",
		},
	}
	markdownCreator.newResources = map[string]string{
		"aws_instance.example":  "terraform generated resource",
		"aws_instance.example2": "terraform generated resource",
	}

	// When
	markdownCreator.setCostsEstimatesData(report)

	// Then
	title := "# Calculable Cloud Costs (Monthly)\n\n"

	tableHeaders := "|Uncontrolled Resources Cost|Terraform Controlled Resources Cost|\n| :---: | :---: |\n"
	tableContent := "|$22.26|$0.00|\n\n"

	expectedMarkdown := fmt.Sprintf("%s%s%s", title, tableHeaders, tableContent)
	assert.Equal(t, expectedMarkdown, report.String())
}
