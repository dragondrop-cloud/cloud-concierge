package costEstimation

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// CostEstimatorConfig is configuration for the CostEstimator struct that conforms
// to envconfig's format expectations.
type CostEstimatorConfig struct {

	// DivisionCloudCredentials is a map between a division and request cloud credentials.
	DivisionCloudCredentials terraformValueObjects.DivisionCloudCredentialDecoder `required:"true"`

	// InfracostAPIToken is the token for accessing Infracost's API.
	InfracostAPIToken string `required:"true"`
}

// CostEstimator is a struct that implements interfaces.CostEstimation.
type CostEstimator struct {

	// config is a struct of configuration parameters
	config CostEstimatorConfig

	// DivisionToProvider is a map between the string representing a division and the corresponding
	// cloud provider (aws, azurerm, google, etc.).
	// For AWS, an account is the division, for GCP a project name is the division,
	// and for azurerm a resource group is a division.
	divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider `required:"true"`
}

// NewCostEstimator creates a new instance of CostEstimator a struct that implements interfaces.CostEstimation.
func NewCostEstimator(config CostEstimatorConfig, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) interfaces.CostEstimation {
	return &CostEstimator{
		config:             config,
		divisionToProvider: divisionToProvider,
	}
}

// Execute creates structured cost estimation data for the current identified/scanned
// cloud resources.
func (ce *CostEstimator) Execute(ctx context.Context) error {
	if ce.config.InfracostAPIToken == "None" {
		fmt.Println("No Infracost token specified, skipping cost estimation.")
		return nil
	}

	// Setting the Infracost API token
	authArgs := []string{"configure", "set", "api_key", ce.config.InfracostAPIToken}
	_, err := executeCommand("infracost", authArgs...)
	if err != nil {
		return fmt.Errorf("[gcloud_authentication][gcloud auth activate-service-account, failed to authenticate]%w", err)
	}
	fmt.Println("Done setting Infracost API token.")

	err = ce.GetAllCostEstimates()
	if err != nil {
		return fmt.Errorf("[ce.GetAllCostEstimates]%v", err)
	}

	err = ce.FormatAllCostEstimates()
	if err != nil {
		return fmt.Errorf("[ce.FormatAllCostEstimates]%v", err)
	}

	err = ce.AggregateCostEstimates()
	if err != nil {
		return fmt.Errorf("[ce.AggregateCostEstimates]%v", err)
	}

	return nil
}

// AggregateCostEstimates merges all calculated and formatted cost estimations into a single
// json object and outputs it to data maps for end consumption.
func (ce *CostEstimator) AggregateCostEstimates() error {
	outputObj := gabs.New()

	for division := range ce.config.DivisionCloudCredentials {
		divisionFolderName := fmt.Sprintf("%v-%v", ce.divisionToProvider[division], division)

		infracostJSONPath := fmt.Sprintf("./current_cloud/%v/infracost-formatted.json", divisionFolderName)

		divisionCosts, err := os.ReadFile(infracostJSONPath)
		if err != nil {
			return fmt.Errorf("[os.ReadFile]%v", err)
		}

		costContainer, err := gabs.ParseJSON(divisionCosts)
		if err != nil {
			return fmt.Errorf("[gabs.ParseJSON]%v", err)
		}

		_, err = outputObj.Set(costContainer, divisionFolderName)
		if err != nil {
			return fmt.Errorf("[outputObj.Set()]%v", err)
		}
	}

	err := os.WriteFile("mappings/division-to-cost-estimates.json", outputObj.Bytes(), 0400)
	if err != nil {
		return fmt.Errorf("[os.WriteFile]%v", err)
	}
	return nil
}

// GetAllCostEstimates invokes the infracost CLI to generate cost estimates for identified resources
// within a all cloud divisions.
func (ce *CostEstimator) GetAllCostEstimates() error {
	for division := range ce.config.DivisionCloudCredentials {
		err := ce.GetDivisionCostEstimate(division)
		if err != nil {
			return fmt.Errorf("[ce.GetDivisionCostEstimate for division %v]%v", division, err)
		}
	}
	return nil
}

// GetDivisionCostEstimate invokes the infracost CLI to generate cost estimates for identified resources
// within a single, specified, cloud division.
func (ce *CostEstimator) GetDivisionCostEstimate(division terraformValueObjects.Division) error {
	divisionFolderName := fmt.Sprintf("%v-%v", ce.divisionToProvider[division], division)

	infracostEstimationPath := fmt.Sprintf("./current_cloud/%v/", divisionFolderName)
	infracostJSONPath := fmt.Sprintf("./current_cloud/%v/infracost.json", divisionFolderName)

	costEstimateArgs := []string{"breakdown", "--path", infracostEstimationPath, "--format", "json", "--out-file", infracostJSONPath}
	_, err := executeCommand("infracost", costEstimateArgs...)
	if err != nil {
		return fmt.Errorf("[executeCommand]%v", err)
	}

	return nil
}

// executeCommand wraps os.exec.Command with capturing of std output and errors.
func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Setting up logging objects
	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("[error executing command: %s, %s]%w", stderr.String(), out.String(), err)
	}
	fmt.Printf("\n%s Output:\n\n%v\n", command, out.String())
	return out.String(), nil
}
