package markdowncreation

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/atsushinee/go-markdown-generator/doc"
)

// OutputPath is the path where the markdown file will be created
const OutputPath = "state_of_cloud"

// ManagedDriftResource represents a resource that is managed by terraform but has drifted
type ManagedDriftResource struct {
	RecentActor           string `json:"RecentActor"`
	RecentActionTimestamp string `json:"RecentActionTimestamp"`
	AttributeName         string `json:"AttributeName"`
	TerraformValue        string `json:"TerraformValue"`
	CloudValue            string `json:"CloudValue"`
	InstanceID            string `json:"InstanceID"`
	InstanceRegion        string `json:"InstanceRegion"`
	StateFileName         string `json:"StateFileName"`
	ModuleName            string `json:"ModuleName"`
	ResourceType          string `json:"ResourceType"`
	ResourceName          string `json:"ResourceName"`
}

// CostEstimate represents a cost estimate for a resource
type CostEstimate struct {
	CostComponent   string `json:"cost_component"`
	IsUsageBased    bool   `json:"is_usage_based"`
	MonthlyCost     string `json:"monthly_cost"`
	MonthlyQuantity string `json:"monthly_quantity"`
	Price           string `json:"price"`
	ResourceName    string `json:"resource_name"`
	SubResourceName string `json:"sub_resource_name"`
	Unit            string `json:"unit"`
}

// CloudActionDetail represents the details of a cloud action
type CloudActionDetail struct {
	Actor     string `json:"actor"`
	Timestamp string `json:"timestamp"`
}

// DeletedResource represents a resource that was deleted
type DeletedResource struct {
	InstanceID    string `json:"InstanceID"`
	StateFileName string `json:"StateFileName"`
	ModuleName    string `json:"ModuleName"`
	ResourceType  string `json:"ResourceType"`
	ResourceName  string `json:"ResourceName"`
}

// MarkdownCreator is responsible for creating the markdown file with the data from the state of cloud
type MarkdownCreator struct {
	newResources            map[string]string
	resourcesToCloudActions map[string]map[string]CloudActionDetail
	costEstimates           []CostEstimate
	securityScan            []SecurityRisk
	managedDrift            []ManagedDriftResource
	deletedResources        []DeletedResource
}

// NewMarkdownCreator returns a new MarkdownCreator
func NewMarkdownCreator() *MarkdownCreator {
	return &MarkdownCreator{}
}

// CreateMarkdownFile creates a markdown file with the data from the state of cloud
func (m *MarkdownCreator) CreateMarkdownFile(jobName string) error {
	err := m.initData()
	if err != nil {
		return fmt.Errorf("[markdown_creator][create_markdown_file] error initializing data: %w", err)
	}

	report := doc.NewMarkDown()

	m.setGeneralData(report, jobName)
	m.setSecurityRiskData(report)
	m.setCostsEstimatesData(report)
	m.setResourcesOutsideOfTerraformControlData(report)
	m.setDeletedResourcesData(report)
	m.setDriftedResourcesManagedByTerraformData(report)
	m.setRootCausesOfDriftData(report)
	m.setFooter(report)

	err = os.Mkdir(OutputPath, 0o755)
	if err != nil {
		return fmt.Errorf("[markdown_creator][create_markdown_file] error creating output directory: %w", err)
	}

	err = report.Export(fmt.Sprintf("%s/report.md", OutputPath))
	if err != nil {
		return fmt.Errorf("[markdown_creator][create_markdown_file] error creating file: %w", err)
	}

	return nil
}

// initData initializes the data from the files
func (m *MarkdownCreator) initData() error {
	filePathRoot := "outputs/"
	newResourcesBytes, err := readFile(filePathRoot + "new-resources-to-documents.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading new resources file: %w", err)
	}
	var newResources map[string]string
	err = json.Unmarshal(newResourcesBytes, &newResources)
	if err != nil {
		return fmt.Errorf("error parsing JSON from resources new resources: %v", err)
	}

	resourcesToCloudActionsBytes, err := readFile(filePathRoot + "resources-to-cloud-actions.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading resources to cloud actions file: %w", err)
	}
	var resourcesToCloudActions map[string]map[string]CloudActionDetail
	err = json.Unmarshal(resourcesToCloudActionsBytes, &resourcesToCloudActions)
	if err != nil {
		return fmt.Errorf("error parsing JSON from resources to cloud actions: %v", err)
	}

	costEstimatesBytes, err := readFile(filePathRoot + "cost-estimates.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading cost estimates file: %w", err)
	}
	var costEstimates []CostEstimate
	err = json.Unmarshal(costEstimatesBytes, &costEstimates)
	if err != nil {
		return fmt.Errorf("error parsing JSON from cost estimates: %v", err)
	}

	securityScanBytes, err := readFile(filePathRoot + "security-scan.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading security scan file: %w", err)
	}
	var securityScan map[string][]SecurityRisk
	err = json.Unmarshal(securityScanBytes, &securityScan)
	if err != nil {
		return fmt.Errorf("error parsing JSON from security scans: %v", err)
	}

	managedDriftBytes, err := readFile(filePathRoot + "drift-resources-differences.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading drift resources differences file: %w", err)
	}
	var managedDrift []ManagedDriftResource
	err = json.Unmarshal(managedDriftBytes, &managedDrift)
	if err != nil {
		return fmt.Errorf("error parsing JSON from managed drifted sources: %v", err)
	}

	deletedResourcesBytes, err := readFile(filePathRoot + "drift-resources-deleted.json")
	if err != nil {
		return fmt.Errorf("[markdown_creator][init_data] error reading drift resources deleted file: %w", err)
	}
	var deletedResources []DeletedResource
	err = json.Unmarshal(deletedResourcesBytes, &deletedResources)
	if err != nil {
		return fmt.Errorf("error parsing JSON from deleted resources: %v", err)
	}

	m.newResources = newResources
	m.resourcesToCloudActions = resourcesToCloudActions
	m.costEstimates = costEstimates
	m.securityScan = securityScan["results"]
	m.managedDrift = managedDrift
	m.deletedResources = deletedResources

	return nil
}

// setGeneralData sets the general data of the report
func (m *MarkdownCreator) setGeneralData(report *doc.MarkDownDoc, jobName string) {
	report.Write(fmt.Sprintf("%s - State of Scanned Cloud Resources", jobName)).Writeln()
	report.Write("========================================================").Writeln().Writeln()

	report.Write("# How to Read this Report").Writeln().Writeln()
	report.Write(fmt.Sprintf("'%s' has run. Of the resources the execution scans, at least one resource was identified to have drifted or "+
		"be outside of Terraform control. While code has been generated of the Terraform code and corresponding import statements needed to "+
		"bring these resources under Terraform control, below you will find a summary of the gaps identified in your current IaC posture.",
		jobName)).Writeln().Writeln()
}

// setFooter sets the footer of the report
func (m *MarkdownCreator) setFooter(report *doc.MarkDownDoc) {
	report.Write("#### Disclaimer").Writeln().Writeln()

	report.Write(
		"*Indicates that a resource's cost is usage based. Since we currently do not infer/have knowledge of usage, costs may be material although indicated as 0 here.",
	).Writeln().Writeln()

	report.Write(
		"This report presents information on the state of your cloud at a point in time and as best Cloud Concierge is able to determine. Cloud Concierge does not currently scan every cloud resource for every cloud provider. For a list of supported resources, please see our [documentation](https://www.docs.dragondrop.cloud/).",
	).Writeln().Writeln()

	currentTime := time.Now().UTC()
	report.Write(currentTime.Format("Created by Cloud Concierge at 15:04 UTC on 2006-01-02"))
}

// readFile reads a file and returns the bytes
func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", path, err)
	}

	return data, nil
}
