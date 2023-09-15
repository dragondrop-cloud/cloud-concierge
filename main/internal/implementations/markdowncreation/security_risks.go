package markdowncreation

import (
	"fmt"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// SecurityRisk represents a security risk identified by tfsec
type SecurityRisk struct {
	ID              string   `json:"id"`
	RuleDescription string   `json:"rule_description"`
	Severity        string   `json:"severity"`
	Resolution      string   `json:"resolution"`
	Links           []string `json:"links"`
}

// setSecurityRiskData sets the security risk data in the markdown report
func (m *MarkdownCreator) setSecurityRiskData(report *doc.MarkDownDoc) {
	report.Write("# Identified Security Risks").Writeln().Writeln()

	if len(m.securityScan) == 0 {
		report.Write("No security risks identified within scanned resources.").Writeln()
		return
	}

	securityRisksByID := make(map[string][]SecurityRisk)
	for _, securityRisk := range m.securityScan {
		securityRisksByID[securityRisk.ID] = append(securityRisksByID[securityRisk.ID], securityRisk)
	}

	for _, securityRisk := range securityRisksByID {
		report.Write(fmt.Sprintf("**Instance ID**: `%s`", securityRisk[0].ID)).Writeln()

		report.Write("|Rule Description")
		report.Write("|Severity")
		report.Write("|Resolution")
		report.Write("|Doc Links|").Writeln()
		report.Write("| :---: | :---: | :---: | :---: |").Writeln()

		for _, risk := range securityRisk {
			report.Write(fmt.Sprintf("|%s", risk.RuleDescription))
			report.Write(fmt.Sprintf("|%s", risk.Severity))
			report.Write(fmt.Sprintf("|%s", risk.Resolution))
			report.Write(fmt.Sprintf("|[Rule](%s), [Tf Doc](%s)|", risk.Links[0], risk.Links[1])).Writeln()
		}

		report.Writeln()
	}
}
