package terraformerCLI

import (
	"encoding/json"
	"fmt"
	"os"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

var defaultAzureRegions = []string{"eastus"}

// AzureScanner implements the Scanner interface for use with Azure cloud environments.
type AzureScanner struct {
	// Config is the needed configuration of a mapping between Division name and the corresponding
	// Credential needed to access that environment.
	config map[terraformValueObjects.Division]terraformValueObjects.Credential

	// terraformer is the TerraformerCLI interface used to scan the Azure cloud environment.
	terraformer TerraformerCLI

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions []terraformValueObjects.CloudRegion `required:"true"`
}

// NewAzureScanner creates and returns a new instance of AzureScanner.
func NewAzureScanner(config map[terraformValueObjects.Division]terraformValueObjects.Credential, cliConfig Config, cloudRegions []terraformValueObjects.CloudRegion) (Scanner, error) {
	return &AzureScanner{
		CloudRegions: cloudRegions,
		config:       config,
		terraformer:  newTerraformerCLI(cliConfig),
	}, nil
}

// ScanAll wraps Scan to scan each division for the provider.
func (azureScanner *AzureScanner) ScanAll(options ...string) (*MultiScanResult, error) {
	fmt.Println("Scanning all specified azure divisions.")
	scanMap := make(map[terraformValueObjects.Division]terraformValueObjects.Path)

	for division, credential := range azureScanner.config {
		path, err := azureScanner.Scan(division, credential)
		if err != nil {
			return nil, fmt.Errorf("[ScanAll] Error in azureScanner.Scan: %v", err)
		}
		scanMap[division] = path
	}

	return &MultiScanResult{scanMap}, nil
}

// AzureEnvironment represents the configuration to run terraformer for Azure
type AzureEnvironment struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantID       string `json:"tenant_id"`
	SubscriptionID string `json:"subscription_id"`
}

// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
func (azureScanner *AzureScanner) Scan(resourceGroup terraformValueObjects.Division, credential terraformValueObjects.Credential, options ...string) (terraformValueObjects.Path, error) {
	env := new(AzureEnvironment)
	err := json.Unmarshal([]byte(credential), &env)
	if err != nil {
		return "", fmt.Errorf("[azure_scanner][configure_environment][error unmarshalling credentials] %w", err)
	}

	err = azureScanner.configureEnvironment(*env)
	if err != nil {
		return "", fmt.Errorf("[Azure Scanner] Error configuring environment %w", err)
	}

	filterValue := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", env.SubscriptionID, resourceGroup)

	path, err := azureScanner.terraformer.Import(TerraformImportMigrationGeneratorParams{
		Provider:       "azurerm",
		Division:       resourceGroup,
		Resources:      []string{},
		AdditionalArgs: []string{fmt.Sprintf("--filter=resource_group=%s", filterValue)},
		Regions:        []string{},
		IsCompact:      true,
	})

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.Import(): %v", err)
	}

	err = azureScanner.terraformer.UpdateState("azurerm", string(path))

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.UpdateState(): %v", err)
	}

	return path, nil
}

func (azureScanner *AzureScanner) configureEnvironment(env AzureEnvironment) error {
	err := os.Setenv("ARM_CLIENT_ID", env.ClientID)
	if err != nil {
		return fmt.Errorf("[azure_scanner][configure_environment][error setting client_id credential] %w", err)
	}

	err = os.Setenv("ARM_CLIENT_SECRET", env.ClientSecret)
	if err != nil {
		return fmt.Errorf("[azure_scanner][configure_environment][error setting client_secret credential] %w", err)
	}

	err = os.Setenv("ARM_TENANT_ID", env.TenantID)
	if err != nil {
		return fmt.Errorf("[azure_scanner][configure_environment][error setting tenant_id credential] %w", err)
	}

	err = os.Setenv("ARM_SUBSCRIPTION_ID", env.SubscriptionID)
	if err != nil {
		return fmt.Errorf("[azure_scanner][configure_environment][error setting subscription_id credential] %w", err)
	}

	return nil
}
