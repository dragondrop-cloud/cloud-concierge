package markdowncreation

import (
	"fmt"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownCreator_setRootCausesOfDriftData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.resourcesToCloudActions = map[string]map[string]CloudActionDetail{
		"type1.resource1": {
			"creation": {
				Actor:     "root",
				Timestamp: "2021-01-01",
			},
		},
		"type1.resource2": {
			"creation": {
				Actor:     "root",
				Timestamp: "2021-01-01",
			},
		},
		"type1.resource3": {
			"modified": {
				Actor:     "root",
				Timestamp: "2021-01-01",
			},
		},
	}

	// When
	markdownCreator.setRootCausesOfDriftData(report)

	// Then
	title := "# Root Causes of Drift\n\n"
	subtitle := "## Cloud Actors Causing Changes\n\n"

	tableHeaders := "|Actor|Action|Count|\n| :---: | :---: | :---: |\n"
	tableContent := "|root|Create Resource|2|\n|root|Modify Resource|1|\n\n"

	expectedMarkdown := fmt.Sprintf("%s%s%s%s", title, subtitle, tableHeaders, tableContent)
	assert.Equal(t, expectedMarkdown, report.String())
}
