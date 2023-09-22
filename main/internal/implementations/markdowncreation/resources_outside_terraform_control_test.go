package markdowncreation

import (
	"strings"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData_WithCosts(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.newResources = map[string]string{
		"aws_db_snapshot.example":      "terraform generated resource",
		"aws_db_subnet_group.example2": "terraform generated resource",
		"aws_lb_listener.example3":     "terraform generated resource",
		"aws_s3_bucket.example4":       "terraform generated resource",
		"aws_vpc.example5":             "terraform generated resource",
		"aws_vpc.example6":             "terraform generated resource",
	}
	markdownCreator.costEstimates = []CostEstimate{
		{
			CostComponent: "Storage",
			MonthlyCost:   "16.425",
			ResourceName:  "aws_db_snapshot.example",
		},
		{
			CostComponent: "Storage",
			Price:         "5.84",
			ResourceName:  "aws_db_subnet_group.example2",
			IsUsageBased:  false,
		},
		{
			CostComponent: "Container",
			Price:         "12.847",
			ResourceName:  "aws_s3_bucket.example3",
			IsUsageBased:  true,
		},
		{
			CostComponent: "VPC Resource1",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example5",
			IsUsageBased:  false,
		},
		{
			CostComponent: "VPC Resource2",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example6",
			IsUsageBased:  false,
		},
	}

	// When
	markdownCreator.setResourcesOutsideOfTerraformControlData(report)

	// Then
	title := "# Resources Outside of Terraform Control\n\n"

	tableHeaders := "|Type|# Resources|Cost Components|Monthly Cost|Usage Based*|" +
		"\n| :---: | :---: | :---: | :---: | :---: |\n"

	reportResult := report.String()
	assert.Equal(t, title, reportResult[:len(title)])
	assert.Equal(t, tableHeaders, reportResult[len(title):len(title)+len(tableHeaders)])

	resourcesValues := strings.Split(reportResult[len(title)+len(tableHeaders):], "\n")
	assert.Equal(t, 7, len(resourcesValues))
	require.Contains(t, resourcesValues, "|aws_db_snapshot|1|1|$16.42|False|")
	require.Contains(t, resourcesValues, "|aws_db_subnet_group|1|1|$5.84|False|")
	require.Contains(t, resourcesValues, "|aws_lb_listener|1|No Charge|No Charge|No Charge|")
	require.Contains(t, resourcesValues, "|aws_s3_bucket|1|1|$0.00*|True|")
	require.Contains(t, resourcesValues, "|aws_vpc|2|2|$25.69|False|")
}

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData_(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.newResources = map[string]string{
		"aws_vpc.example5": "terraform generated resource",
		"aws_vpc.example6": "terraform generated resource",
	}
	markdownCreator.costEstimates = []CostEstimate{
		{
			CostComponent: "VPC Resource1",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example5",
			IsUsageBased:  true,
		},
		{
			CostComponent: "VPC Resource2",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example6",
			IsUsageBased:  true,
		},
	}

	// When
	markdownCreator.setResourcesOutsideOfTerraformControlData(report)

	// Then
	title := "# Resources Outside of Terraform Control\n\n"

	tableHeaders := "|Type|# Resources|Cost Components|Monthly Cost|Usage Based*|" +
		"\n| :---: | :---: | :---: | :---: | :---: |\n"

	reportResult := report.String()
	assert.Equal(t, title, reportResult[:len(title)])
	assert.Equal(t, tableHeaders, reportResult[len(title):len(title)+len(tableHeaders)])

	resourcesValues := strings.Split(reportResult[len(title)+len(tableHeaders):], "\n")
	assert.Equal(t, 3, len(resourcesValues))
	require.Contains(t, resourcesValues, "|aws_vpc|2|2|$0.00*|True|")
}

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData_NoCostEstimationFound(t *testing.T) {
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

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData_NoCostEstimationForNewResources(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.newResources = map[string]string{
		"aws_db_snapshot.example":      "terraform generated resource",
		"aws_db_subnet_group.example2": "terraform generated resource",
		"aws_lb_listener.example3":     "terraform generated resource",
	}
	markdownCreator.costEstimates = []CostEstimate{
		{
			CostComponent: "Container",
			Price:         "12.847",
			ResourceName:  "aws_s3_bucket.example3",
			IsUsageBased:  true,
		},
		{
			CostComponent: "VPC Resource1",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example5",
			IsUsageBased:  false,
		},
		{
			CostComponent: "VPC Resource2",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example6",
			IsUsageBased:  false,
		},
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

func TestMarkdownCreator_setResourcesOutsideOfTerraformControlData_OnlyOneValidCostEstimationCharge(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.newResources = map[string]string{
		"aws_db_snapshot.example":      "terraform generated resource",
		"aws_db_subnet_group.example2": "terraform generated resource",
		"aws_lb_listener.example3":     "terraform generated resource",
	}
	markdownCreator.costEstimates = []CostEstimate{
		{
			CostComponent: "Container",
			Price:         "12.847",
			ResourceName:  "aws_s3_bucket.example3",
			IsUsageBased:  true,
		},
		{
			CostComponent: "VPC Resource1",
			Price:         "12.847",
			ResourceName:  "aws_vpc.example5",
			IsUsageBased:  false,
		},
		{
			CostComponent: "lb listener",
			Price:         "12.843",
			ResourceName:  "aws_lb_listener.example3",
			IsUsageBased:  false,
		},
	}

	// When
	markdownCreator.setResourcesOutsideOfTerraformControlData(report)

	// Then
	title := "# Resources Outside of Terraform Control\n\n"

	tableHeaders := "|Type|# Resources|Cost Components|Monthly Cost|Usage Based*|" +
		"\n| :---: | :---: | :---: | :---: | :---: |\n"

	reportResult := report.String()
	assert.Equal(t, title, reportResult[:len(title)])
	assert.Equal(t, tableHeaders, reportResult[len(title):len(title)+len(tableHeaders)])

	resourcesValues := strings.Split(reportResult[len(title)+len(tableHeaders):], "\n")
	assert.Equal(t, 5, len(resourcesValues))
	require.Contains(t, resourcesValues, "|aws_db_snapshot|1|No Charge|No Charge|No Charge|")
	require.Contains(t, resourcesValues, "|aws_db_subnet_group|1|No Charge|No Charge|No Charge|")
	require.Contains(t, resourcesValues, "|aws_lb_listener|1|1|$12.84|False|")
}
