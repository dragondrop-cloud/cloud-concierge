package markdowncreation

import (
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownCreator_setDriftedResourcesManagedByTerraformData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.managedDrift = []ManagedDriftResource{
		{
			ModuleName:            "module1",
			ResourceType:          "resource_type1",
			ResourceName:          "resource_name1",
			StateFileName:         "state_file_name1",
			InstanceID:            "instance_id1",
			RecentActor:           "recent_actor1",
			RecentActionTimestamp: "2023-01-01",
			AttributeName:         "attribute_name1",
			TerraformValue:        "terraform_value1",
			CloudValue:            "cloud_value1",
		},
		{
			ModuleName:            "module1",
			ResourceType:          "resource_type2",
			ResourceName:          "resource_name2",
			StateFileName:         "state_file_name1",
			InstanceID:            "instance_id2",
			RecentActor:           "recent_actor1",
			RecentActionTimestamp: "2023-01-01",
			AttributeName:         "attribute_name2",
			TerraformValue:        "terraform_value2",
			CloudValue:            "cloud_value2",
		},
	}

	// When
	markdownCreator.setDriftedResourcesManagedByTerraformData(report)

	// Then
	title := "# Drifted Resources Managed By Terraform\n\n"

	stateFile := "## State File `state_file_name1`\n\n"

	resourcePath1 := "### Resource: module1 (module) \"resource_type1\" \"resource_name1\"\n\n"
	instanceID1 := "**Instance ID**: `instance_id1`\n\n"
	actor1 := "**Most Recent Non-Terraform Actor**: `recent_actor1`\n"
	date1 := "**Most Recent Action Date**: `2023-01-01`\n\n"

	complete := "- [ ] Completed\n\n"

	tableHeaders := "|Attribute|Terraform Value|Cloud Value|\n| :---: | :---: | :---: |\n"
	tableContent1 := "|attribute_name1|terraform_value1|cloud_value1|\n\n"

	resourcePath2 := "### Resource: module1 (module) \"resource_type2\" \"resource_name2\"\n\n"
	instanceID2 := "**Instance ID**: `instance_id2`\n\n"
	actor2 := "**Most Recent Non-Terraform Actor**: `recent_actor1`\n"
	date2 := "**Most Recent Action Date**: `2023-01-01`\n\n"
	tableContent2 := "|attribute_name2|terraform_value2|cloud_value2|\n\n"

	actualMarkdown := report.String()
	assert.Equal(t, title, actualMarkdown[:len(title)])

	staticContentLength := len(title) + len(stateFile)
	assert.Equal(t, stateFile, actualMarkdown[len(title):staticContentLength])

	if actualMarkdown[len(title)+len(stateFile):staticContentLength+len(resourcePath1)] == resourcePath1 {
		expectedDynamicContentPath1 := resourcePath1 + instanceID1 + actor1 + date1 + complete + tableHeaders + tableContent1
		assert.Equal(t, expectedDynamicContentPath1, actualMarkdown[staticContentLength:staticContentLength+len(expectedDynamicContentPath1)])

		expectedDynamicContentPath2 := resourcePath2 + instanceID2 + actor2 + date2 + complete + tableHeaders + tableContent2
		assert.Equal(t, expectedDynamicContentPath2, actualMarkdown[staticContentLength+len(expectedDynamicContentPath1):])
	} else {
		expectedDynamicContentPath2 := resourcePath2 + instanceID2 + actor2 + date2 + complete + tableHeaders + tableContent2
		assert.Equal(t, expectedDynamicContentPath2, actualMarkdown[staticContentLength:staticContentLength+len(expectedDynamicContentPath2)])

		expectedDynamicContentPath1 := resourcePath1 + instanceID1 + actor1 + date1 + complete + tableHeaders + tableContent1
		assert.Equal(t, expectedDynamicContentPath1, actualMarkdown[staticContentLength+len(expectedDynamicContentPath2):])
	}
}
