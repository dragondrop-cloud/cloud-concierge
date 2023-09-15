package markdowncreation

import (
	"fmt"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownCreator_setDeletedResourcesData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.deletedResources = []DeletedResource{
		{
			InstanceID:    "i-1234567890abcdef0",
			ResourceType:  "aws_instance",
			ResourceName:  "example",
			ModuleName:    "module.example",
			StateFileName: "terraform.tfstate",
		},
		{
			InstanceID:    "i-1234567890abcdef1",
			ResourceType:  "aws_instance",
			ResourceName:  "example2",
			ModuleName:    "module.example2",
			StateFileName: "terraform.tfstate",
		},
	}

	// When
	markdownCreator.setDeletedResourcesData(report)

	// Then
	title := "# Drifted Resources Deleted\n\n"

	tableHeaders := "|Type|Name|Module|State File|\n| :---: | :---: | :---: | :---: |\n"
	tableContent := "|aws_instance|example|module.example|terraform.tfstate|\n|aws_instance|example2|module.example2|terraform.tfstate|\n\n"

	expectedMarkdown := fmt.Sprintf("%s%s%s", title, tableHeaders, tableContent)
	assert.Equal(t, expectedMarkdown, report.String())
}
