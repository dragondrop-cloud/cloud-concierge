package costestimation

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// CostEstimatorConfig is configuration for the CostEstimator struct that conforms
// to envconfig's format expectations.
type CostEstimatorConfig struct {
	// CloudCredential is a cloud credential with read-only access to a cloud division and, if applicable, access to read Terraform state files.
	CloudCredential terraformValueObjects.Credential

	// InfracostAPIToken is the token for accessing Infracost's API.
	InfracostAPIToken string
}

// CostEstimator is a struct that implements interfaces.CostEstimation.
type CostEstimator struct {
	// config is a struct of configuration parameters
	config CostEstimatorConfig

	// DivisionToProvider is a map between the string representing a division and the corresponding
	// cloud provider (aws, azurerm, google, etc.).
	// For AWS, an account is the division, for GCP a project name is the division,
	// and for azurerm a resource group is a division.
	provider terraformValueObjects.Provider `required:"true"`
}

// NewCostEstimator creates a new instance of CostEstimator a struct that implements interfaces.CostEstimation.
func NewCostEstimator(config CostEstimatorConfig, provider terraformValueObjects.Provider) interfaces.CostEstimation {
	return &CostEstimator{
		config:   config,
		provider: provider,
	}
}

// SetInfracostAPIToken sets the Infracost API token.
func (ce *CostEstimator) SetInfracostAPIToken(token string) {
	ce.config.InfracostAPIToken = token
}

// Execute creates structured cost estimation data for the current identified/scanned
// cloud resources.
func (ce *CostEstimator) Execute(_ context.Context) error {
	logrus.Debugf("Executing cost estimation for %s", ce.provider)

	// Setting the Infracost API token
	authArgs := []string{"configure", "set", "api_key", ce.config.InfracostAPIToken}
	_, err := executeCommand("infracost", authArgs...)
	if err != nil {
		return fmt.Errorf("[failed to set infracost api_key value]%w", err)
	}
	logrus.Info("Done setting Infracost API token.")

	err = ce.GetCostEstimate()
	if err != nil {
		return fmt.Errorf("[ce.GetAllCostEstimates]%v", err)
	}

	err = ce.FormatCostEstimate()
	if err != nil {
		return fmt.Errorf("[ce.FormatAllCostEstimates]%v", err)
	}

	err = ce.WriteCostEstimates()
	if err != nil {
		return fmt.Errorf("[ce.WriteCostEstimates]%v", err)
	}

	return nil
}

// WriteCostEstimates outputs cost estimates into data maps for end consumption.
func (ce *CostEstimator) WriteCostEstimates() error {
	infracostJSONPath := "./current_cloud/infracost-formatted.json"

	costs, err := os.ReadFile(infracostJSONPath)
	if err != nil {
		return fmt.Errorf("[os.ReadFile]%v", err)
	}

	err = os.WriteFile("outputs/cost-estimates.json", costs, 0o400)
	if err != nil {
		return fmt.Errorf("[os.WriteFile]%v", err)
	}
	return nil
}

// GetCostEstimate invokes the infracost CLI to generate cost estimates for identified resources
// within a single, specified, cloud division.
func (ce *CostEstimator) GetCostEstimate() error {
	infracostEstimationPath := "./current_cloud/"
	infracostJSONPath := "./current_cloud/infracost.json"

	costEstimateArgs := []string{"breakdown", "--path", infracostEstimationPath, "--format", "json", "--out-file", infracostJSONPath}
	output, err := executeCommand("infracost", costEstimateArgs...)
	if err != nil {
		return fmt.Errorf("[executeCommand]%v", err)
	}

	logrus.Debugf("Infracost output: %s", output)
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
