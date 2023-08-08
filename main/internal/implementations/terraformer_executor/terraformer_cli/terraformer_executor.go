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
	// CloudCredential is a cloud credential with read-only access to a cloud division and, if applicable, access to read Terraform state files.
	CloudCredential terraformValueObjects.Credential `required:"true"`

	// Provider is a map between a cloud provider and the version for that provider.
	Provider map[terraformValueObjects.Provider]string `required:"true"`

	// TerraformVersion is the version of Terraform used.
	TerraformVersion terraformValueObjects.Version `required:"true"`

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions terraformValueObjects.CloudRegionsDecoder `required:"true"`
}

// TerraformerExecutor is a struct that implements interfaces.TerraformerExecutor
type TerraformerExecutor struct {
	// hclCreate implements the hclcreate.HCLCreate interface
	hclCreate hclcreate.HCLCreate

	// scanner is an instantiation of the current provider's scanner.
	scanner Scanner

	// dragonDrop is needed to inform scanned
	dragonDrop interfaces.DragonDrop

	// config contains the variables that determine the specific behavior of the TerraformerExecutor
	config TerraformerExecutorConfig
}

// NewTerraformerExecutor creates and returns a new instance of TerraformerExecutor.
func NewTerraformerExecutor(ctx context.Context, hclCreate hclcreate.HCLCreate, dragonDrop interfaces.DragonDrop, config TerraformerExecutorConfig, cliConfig Config, provider terraformValueObjects.Provider) (interfaces.TerraformerExecutor, error) {
	scanner, err := getScanner(config, cliConfig, provider)
	if err != nil {
		return nil, err
	}

	dragonDrop.PostLog(ctx, "Created TFExec.")
	return &TerraformerExecutor{hclCreate: hclCreate, scanner: scanner, config: config, dragonDrop: dragonDrop}, nil
}

// getScanner provisions the cloud environment scanner for the specified provider.
func getScanner(config TerraformerExecutorConfig, cliConfig Config, provider terraformValueObjects.Provider) (Scanner, error) {

	switch provider {
	case "google":
		googleScanner, err := NewGoogleScanner(config.CloudCredential, cliConfig, config.CloudRegions)

		if err != nil {
			log.Errorf("[NewTerraformerExec] Error in NewGoogleScanner(): %s", err.Error())
			return nil, fmt.Errorf("[NewTerraformerExec] Error in NewGoogleScanner(): %w", err)
		}

		return googleScanner, nil
	case "aws":
		awsScanner, err := NewAWSScanner(config.CloudCredential, cliConfig, config.CloudRegions)

		if err != nil {
			log.Errorf("[NewTerraformerExec] Error in NewAWSScanner(): %s", err.Error())
			return nil, fmt.Errorf("[NewTerraformerExec] Error in NewAWSScanner(): %w", err)
		}

		return awsScanner, nil
	case "azurerm":
		azureScanner, err := NewAzureScanner(config.CloudCredential, cliConfig, config.CloudRegions)

		if err != nil {
			log.Errorf("[NewTerraformerExec] Error in NewAzureScanner(): %s", err.Error())
			return nil, fmt.Errorf("[NewTerraformerExec] Error in NewAzureScanner(): %w", err)
		}

		return azureScanner, nil
	default:
		log.Errorf("currently only a scanner for [google, aws, azurerm] is supported. Specified %s", provider)
		return nil, fmt.Errorf("currently only a scanner for [google, aws, azurerm] is supported. Specified %s", provider)
	}
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

	err = e.scanner.Scan("", e.config.CloudCredential) // TODO: Will need to load division name from config in a future PR.
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

	for provider, version := range e.config.Provider {
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
