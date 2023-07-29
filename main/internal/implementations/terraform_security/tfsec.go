package terraformSecurity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// TFSecFileBytesPerDivision is a map that relates the division to the bytes of the tfsec command results
type TFSecFileBytesPerDivision map[terraformValueObjects.Division][]byte

// TFSecParsedFilePerDivision is a map that relates the division to the parsed TFSecFile
type TFSecParsedFilePerDivision map[terraformValueObjects.Division]TFSecFile

// TFSecResultsPerDivision is a map that relates the division to a Result array
type TFSecResultsPerDivision map[terraformValueObjects.Division][]Result

// TFSecFile is a structure that represents the content of a tfsec command output file
type TFSecFile struct {
	Results []Result `json:"results"`
}

// Result is a struct that represents a single result in the array of results of a tfsec command output file
type Result struct {
	ID              string   `json:"id"`
	RuleID          string   `json:"rule_id"`
	LongID          string   `json:"long_id"`
	RuleDescription string   `json:"rule_description"`
	RuleProvider    string   `json:"rule_provider"`
	RuleService     string   `json:"rule_service"`
	Impact          string   `json:"impact"`
	Resolution      string   `json:"resolution"`
	Links           []string `json:"links"`
	Description     string   `json:"description"`
	Severity        string   `json:"severity"`
	Warning         bool     `json:"warning"`
	Status          int      `json:"status"`
	Resource        string   `json:"resource"`
	Location        Location `json:"location"`
}

// Location is a struct that represents the location of the file that contains the security mistake
type Location struct {
	FileName  string `json:"file_name"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

// TFSec is a struct that implements the interfaces.TerraformSecurity but
// executing the tfsec command
type TFSec struct {
	// DivisionToProvider is a map between the string representing a division and the corresponding
	// cloud provider (aws, azurerm, google, etc.).
	// For AWS, an account is the division, for GCP a project name is the division,
	// and for azurerm a resource group is a division.
	divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider
}

// NewTFSec generates a new instance from TFSec
func NewTFSec(divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) *TFSec {
	return &TFSec{
		provider: provider,
	}
}

// ExecuteScan is called from the main job flow to execute the tfsec command and save the output
// to show to the user in the PR
func (s *TFSec) ExecuteScan(ctx context.Context) error {
	contentResults, err := s.runTFSec()
	if err != nil {
		return fmt.Errorf("[tfsec][execute_scan][error running tfsec command][%v]", err)
	}

	parsedContentResults, err := s.parseContentResults(contentResults)
	if err != nil {
		return fmt.Errorf("[tfsec][execute_scan][error parsing tfsec results][%v]", err)
	}

	mergedResults, err := s.mergeTFSecResultsPerDivision(parsedContentResults)
	if err != nil {
		return fmt.Errorf("[tfsec][execute_scan][error merging tfsec results][%v]", err)
	}

	mergedResultsWithID, err := s.addIDToResources(mergedResults)
	if err != nil {
		return fmt.Errorf("[tfsec][execute_scan][error adding the id to the tfsec results][%v]", err)
	}

	err = s.writeResultsToMappingFile(mergedResultsWithID)
	if err != nil {
		return fmt.Errorf("[tfsec][execute_scan][error writing tfsec results][%v]", err)
	}

	return nil
}

// runTFSec runs the tfsec command through the directories from the divisions configured by the user
func (s *TFSec) runTFSec() ([]byte, error) {
	tfsecScanningPath := "./current_cloud/"
	outLocationFlag := "./current_cloud/tfsec.json"
	outFlag := fmt.Sprintf("--out=%s", outLocationFlag)

	cmd := exec.Command("tfsec", outFlag, "--format=json", "--soft-fail", tfsecScanningPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", out.String(), err)
	}

	results, err := os.ReadFile(outLocationFlag)
	if err != nil {
		return nil, fmt.Errorf("[os.ReadFile][%v]", err)
	}

	return results, nil
}

// parseContentResults takes the bytes of the output tfsec results and returns the same bytes parsed
func (s *TFSec) parseContentResults(contentResults TFSecFileBytesPerDivision) (TFSecParsedFilePerDivision, error) {
	tfSecFiles := map[terraformValueObjects.Division]TFSecFile{}

	for division, results := range contentResults {
		var tfSecFile TFSecFile
		err := json.Unmarshal(results, &tfSecFile)
		if err != nil {
			return nil, err
		}

		tfSecFiles[division] = tfSecFile
	}

	return tfSecFiles, nil
}

// mergeTFSecResultsPerDivision takes the parsed files and returns only the results of the files grouped by division
func (s *TFSec) mergeTFSecResultsPerDivision(tfSecFiles TFSecParsedFilePerDivision) (TFSecResultsPerDivision, error) {
	mergedResults := map[terraformValueObjects.Division][]Result{}

	for division, tfSecFile := range tfSecFiles {
		mergedResults[division] = tfSecFile.Results
	}

	return mergedResults, nil
}

// writeResultsToMappingFile takes the results grouped by division and writes in the mapping file
func (s *TFSec) writeResultsToMappingFile(results TFSecResultsPerDivision) error {
	differencesJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("mappings/division-to-security-scan.json", differencesJSON, 0400)
}

// TODO: Lots to do here
// addIDToResources takes the results grouped by division and adds the id of the resource
func (s *TFSec) addIDToResources(resultsPerDivision TFSecResultsPerDivision) (TFSecResultsPerDivision, error) {
	mergedResults := map[terraformValueObjects.Division][]Result{}

	for division, results := range resultsPerDivision {
		fileContent, err := os.ReadFile("current_cloud/terraform.tfstate")
		if err != nil {
			return nil, err
		}

		stateFile, err := driftDetector.ParseTerraformerStateFile(fileContent)
		if err != nil {
			return nil, fmt.Errorf("[ParseTerraformerStateFile]%v", err)
		}

		resources := s.mapResourceIDsFromStateFile(stateFile)

		mergedResults[division] = s.getResultsWithResourceID(results, resources)
	}

	return mergedResults, nil
}

// mapResourceIDsFromStateFile maps the resources with its resource identifier
func (s *TFSec) mapResourceIDsFromStateFile(file driftDetector.TerraformerStateFile) map[driftDetector.ResourceIdentifier]string {
	resourcesMap := map[driftDetector.ResourceIdentifier]string{}
	for _, resource := range file.Resources {
		resourceKey := fmt.Sprintf("%s.%s", resource.Type, resource.Name)

		// check if "arn" attribute exists within AttributesFlat
		resourceIDValue, okay := resource.Instances[0].AttributesFlat["arn"]

		// check if the resource instance has an "arn" attribute
		if !okay {
			resourceIDValue = resource.Instances[0].AttributesFlat["id"]
		}

		resourcesMap[driftDetector.ResourceIdentifier(resourceKey)] = resourceIDValue
	}
	return resourcesMap
}

// getResultsWithResourceID maps the Resource IDs into the ID variable of the resources array
func (s *TFSec) getResultsWithResourceID(results []Result, resources map[driftDetector.ResourceIdentifier]string) []Result {
	resultsWithID := make([]Result, 0)
	for _, result := range results {

		// split tfSecResourceName by "." and join the first two elements by "."
		// to match what is in the state file
		tfSecResourceName := strings.Join(strings.Split(result.Resource, ".")[:2], ".")

		result.ID = resources[driftDetector.ResourceIdentifier(tfSecResourceName)]
		resultsWithID = append(resultsWithID, result)
	}
	return resultsWithID
}
