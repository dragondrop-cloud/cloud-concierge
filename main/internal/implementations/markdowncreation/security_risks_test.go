package markdowncreation

import (
	"fmt"
	"testing"

	"github.com/atsushinee/go-markdown-generator/doc"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownCreator_setSecurityRiskData(t *testing.T) {
	// Given
	report := doc.NewMarkDown()
	markdownCreator := NewMarkdownCreator()
	markdownCreator.securityScan = []SecurityRisk{
		{
			ID:              "CVE-2019-1234",
			RuleDescription: "This is a rule description",
			Severity:        "HIGH",
			Resolution:      "This is a resolution",
			Links:           []string{"https://example.com", "https://example2.com"},
		},
		{
			ID:              "CVE-2019-1235",
			RuleDescription: "This is a rule description",
			Severity:        "CRITICAL",
			Resolution:      "This is a resolution",
			Links:           []string{"https://example.com", "https://example2.com"},
		},
		{
			ID:              "CVE-2019-1235",
			RuleDescription: "This is a 2nd rule description",
			Severity:        "HIGH",
			Resolution:      "This is a 2nd resolution",
			Links:           []string{"https://example3.com", "https://example4.com"},
		},
	}

	// When
	markdownCreator.setSecurityRiskData(report)

	// Then
	title := "# Identified Security Risks\n\n"

	instance1 := "**Instance ID**: `CVE-2019-1234`\n"
	instance2 := "**Instance ID**: `CVE-2019-1235`\n"
	tableHeaders := "|Rule Description|Severity|Resolution|Doc Links|\n| :---: | :---: | :---: | :---: |\n"

	table1 := "|This is a rule description|HIGH|This is a resolution|[Rule](https://example.com), [Tf Doc](https://example2.com)|\n\n"
	table2Line1 := "|This is a rule description|CRITICAL|This is a resolution|[Rule](https://example.com), [Tf Doc](https://example2.com)|\n"
	table2Line2 := "|This is a 2nd rule description|HIGH|This is a 2nd resolution|[Rule](https://example3.com), [Tf Doc](https://example4.com)|\n\n"

	expectedMarkdown := fmt.Sprintf(
		"%s%s%s%s%s%s%s%s",
		title, instance1, tableHeaders, table1, instance2, tableHeaders, table2Line1, table2Line2,
	)
	assert.Equal(t, expectedMarkdown, report.String())
}
