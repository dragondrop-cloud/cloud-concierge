package terraformerCLI

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// TerraformerExecutorConfig is a struct containing the variables that determine the specific
// behavior of the TerraformerExecutor.
type TerraformerExecutorConfig struct {
	// DivisionCloudCredentials is a map between a division and request cloud credentials.
	DivisionCloudCredentials terraformValueObjects.DivisionCloudCredentialDecoder `required:"true"`

	// Providers is a map between a cloud provider and the version for that provider.
	Providers map[terraformValueObjects.Provider]string `required:"true"`

	// TerraformVersion is the version of Terraform used.
	TerraformVersion terraformValueObjects.Version `required:"true"`

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions terraformValueObjects.CloudRegionsDecoder `required:"true"`
}

// TerraformerExecutor is a struct that implements interfaces.TerraformerExecutor
type TerraformerExecutor struct {
	// hclCreate implements the hclcreate.HCLCreate interface
	hclCreate hclcreate.HCLCreate

	// scanners is a map between each provider and an instantiation of that provider's scanner.
	scanners map[terraformValueObjects.Provider]Scanner

	// dragonDrop is needed to inform scanned
	dragonDrop interfaces.DragonDrop

	// config contains the variables that determine the specific behavior of the TerraformerExecutor
	config TerraformerExecutorConfig
}

// NewTerraformerExecutor creates and returns a new instance of TerraformerExecutor.
func NewTerraformerExecutor(ctx context.Context, hclCreate hclcreate.HCLCreate, dragonDrop interfaces.DragonDrop, config TerraformerExecutorConfig, cliConfig Config, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) (interfaces.TerraformerExecutor, error) {
	scanners, err := getScanners(config, cliConfig, divisionToProvider)
	if err != nil {
		return nil, err
	}

	dragonDrop.PostLog(ctx, "Created TFExec.")
	return &TerraformerExecutor{hclCreate: hclCreate, scanners: scanners, config: config, dragonDrop: dragonDrop}, nil
}

// getScanners provisions all needed cloud environment scanners by Terraform provider to scan.
func getScanners(config TerraformerExecutorConfig, cliConfig Config, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) (map[terraformValueObjects.Provider]Scanner, error) {
	providerSet := make(map[terraformValueObjects.Provider]bool)
	scanners := make(map[terraformValueObjects.Provider]Scanner)
	for p := range config.Providers {
		providerSet[p] = true
	}

	for p := range providerSet {
		switch p {
		case "google":
			googleScannerConfig := subsetMapOfDivisionToCredentials(config.DivisionCloudCredentials, divisionToProvider, p)
			googleScanner, err := NewGoogleScanner(googleScannerConfig, cliConfig, config.CloudRegions)

			if err != nil {
				log.Errorf("[NewTerraformerExec] Error in NewGoogleScanner(): %s", err.Error())
				return nil, fmt.Errorf("[NewTerraformerExec] Error in NewGoogleScanner(): %w", err)
			}

			scanners[p] = googleScanner
		case "aws":
			awsScannerConfig := subsetMapOfDivisionToCredentials(config.DivisionCloudCredentials, divisionToProvider, p)
			awsScanner, err := NewAWSScanner(awsScannerConfig, cliConfig, config.CloudRegions)

			if err != nil {
				log.Errorf("[NewTerraformerExec] Error in NewAWSScanner(): %s", err.Error())
				return nil, fmt.Errorf("[NewTerraformerExec] Error in NewAWSScanner(): %w", err)
			}

			scanners[p] = awsScanner
		case "azurerm":
			azureScannerConfig := subsetMapOfDivisionToCredentials(config.DivisionCloudCredentials, divisionToProvider, p)
			azureScanner, err := NewAzureScanner(azureScannerConfig, cliConfig, config.CloudRegions)

			if err != nil {
				log.Errorf("[NewTerraformerExec] Error in NewAzureScanner(): %s", err.Error())
				return nil, fmt.Errorf("[NewTerraformerExec] Error in NewAzureScanner(): %w", err)
			}

			scanners[p] = azureScanner
		default:
			log.Errorf("currently only a scanner for [google, aws, azurerm] is supported. Specified %s", p)
			return nil, fmt.Errorf("currently only a scanner for [google, aws, azurerm] is supported. Specified %s", p)
		}
	}

	return scanners, nil
}

// subsetMapOfDivisionToCredentials extracts all the division to cloud credentials pairings for
// a given cloud provider.
func subsetMapOfDivisionToCredentials(divisionCloudCredentials terraformValueObjects.DivisionCloudCredentialDecoder, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider, provider terraformValueObjects.Provider) map[terraformValueObjects.Division]terraformValueObjects.Credential {
	divisionSet := make(map[terraformValueObjects.Division]bool)

	// Unique divisions for the current provider.
	for div, p := range divisionToProvider {
		if p == provider {
			divisionSet[div] = true
		}
	}

	// Filter down division to credential dictionary to subset corresponding with the unique divisions.
	subsetDivisionToCredentials := make(map[terraformValueObjects.Division]terraformValueObjects.Credential)

	for div, cred := range divisionCloudCredentials {
		if divisionSet[div] {
			subsetDivisionToCredentials[div] = cred
		}
	}
	return subsetDivisionToCredentials
}

// Execute runs the workflow needed to capture the current state of an
// external cloud environment via the terraformer package.
func (e *TerraformerExecutor) Execute(ctx context.Context) error {
	e.dragonDrop.PostLog(ctx, "Beginning to make main.tf file.")

	err := e.makeProviderVersionFile()
	if err != nil {
		return fmt.Errorf("[terraformer_executor][set_up][error making provider version file]%w", err)
	}

	e.dragonDrop.PostLog(ctx, "Done with main.tf file. Running `tfswitch`.")

	err = e.setTerraformVersion()
	if err != nil {
		return fmt.Errorf("[terraformer_executor][set_up][error setting terraform version]%w", err)
	}

	e.dragonDrop.PostLog(ctx, "Done with `tfswitch`. Running `terraform init`.")

	err = e.initializeTerraform()
	if err != nil {
		return fmt.Errorf("[terraformer_executor][set_up][error initializing terraform]%w", err)
	}

	e.dragonDrop.PostLog(ctx, "Done with running `terraform init`.\n Beginning to scan existing cloud environment.")

	err = e.scanAllProviders()
	if err != nil {
		return fmt.Errorf("[terraformer_executor][set_up][error scanning all providers]%w", err)
	}

	e.dragonDrop.PostLog(ctx, "Executed terraformer scan.")

	err = e.dragonDrop.InformCloudEnvironmentScanned(ctx)
	if err != nil {
		return fmt.Errorf("[terraformer_executor][set_up][error informing cloud environment scanned]%w", err)
	}

	return nil
}

// scanAllProviders runs terraformer against all specified providers and all divisions,
// within each provider.
func (e *TerraformerExecutor) scanAllProviders() error {
	scanOutput := make(map[terraformValueObjects.Provider]*MultiScanResult)

	for provider, s := range e.scanners {
		currentMultiScan, err := s.ScanAll()

		if err != nil {
			return fmt.Errorf(
				"[scan_all_providers][error in s.ScanAll() for provider %s]%w", provider, err,
			)
		}
		scanOutput[provider] = currentMultiScan
	}

	return nil
}

// initializeTerraform initializes Terraform within the current working directory.
func (e *TerraformerExecutor) initializeTerraform() error {
	err := os.Chdir("current_cloud/")
	if err != nil {
		return fmt.Errorf("[initialize_terraform][error changing working directory]%w", err)
	}

	cmd := exec.Command("terraform", "init")
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("[initialize_terraform][error in running 'terraform init': %s]%w", out.String(), err)
	}
	fmt.Printf("%v", out.String())

	return nil
}

// setTerraformVersion uses tfswitch to install the user-specified version of terraform
func (e *TerraformerExecutor) setTerraformVersion() error {
	tfVersion := string(e.config.TerraformVersion)
	cmd := exec.Command("tfswitch", tfVersion)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"[set_terraform_version][error in running 'tfswitch %s' %s]%w",
			e.config.TerraformVersion, out.String(), err,
		)
	}
	fmt.Printf("%v", out.String())

	return nil
}

// makeProviderVersionFile writes an HCL file which defines all provider versions.
func (e *TerraformerExecutor) makeProviderVersionFile() error {
	genericProviders := make(map[string]string)

	for provider, version := range e.config.Providers {
		genericProviders[string(provider)] = string(version)
	}

	mainTF, err := e.hclCreate.CreateMainTF(genericProviders)

	if err != nil {
		return fmt.Errorf("[make_provider_version_file][error in creating main terraform file]%w", err)
	}

	err = os.MkdirAll("current_cloud", 0660)
	if err != nil {
		return err
	}

	err = os.WriteFile("current_cloud/main.tf", mainTF, 0400)
	if err != nil {
		return fmt.Errorf("[make_provider_version_file][error saving file]%w", err)
	}

	return nil
}
