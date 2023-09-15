package terraformercli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

var defaultAzureRegions = []string{"eastus"}

// AzureScanner implements the Scanner interface for use with Azure cloud environments.
type AzureScanner struct {
	// credential needed to scan an Azure cloud environment.
	credential terraformValueObjects.Credential

	// terraformer is the TerraformerCLI interface used to scan the Azure cloud environment.
	terraformer TerraformerCLI

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions []terraformValueObjects.CloudRegion `required:"true"`
}

// NewAzureScanner creates and returns a new instance of AzureScanner.
func NewAzureScanner(credential terraformValueObjects.Credential, cliConfig Config, cloudRegions []terraformValueObjects.CloudRegion) (Scanner, error) {
	return &AzureScanner{
		CloudRegions: cloudRegions,
		credential:   credential,
		terraformer:  newTerraformerCLI(cliConfig),
	}, nil
}

// AzureEnvironment represents the configuration to run terraformer for Azure
type AzureEnvironment struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantID       string `json:"tenant_id"`
	SubscriptionID string `json:"subscription_id"`
}

// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
func (azureScanner *AzureScanner) Scan(_ terraformValueObjects.Division, credential terraformValueObjects.Credential, _ ...string) error {
	logrus.Debugf("[AzureScanner][Scan] Scanning Azure environment %v", credential)
	env := new(AzureEnvironment)
	credentialBytes := bytes.TrimPrefix([]byte(credential), []byte("\xef\xbb\xbf"))

	err := json.Unmarshal(credentialBytes, &env)
	if err != nil {
		return fmt.Errorf("[azure_scanner][configure_environment][error unmarshalling credentials] %w", err)
	}

	err = azureScanner.configureEnvironment(*env)
	if err != nil {
		return fmt.Errorf("[Azure Scanner] Error configuring environment %w", err)
	}

	filterValue := fmt.Sprintf("/subscriptions/%s", env.SubscriptionID)

	err = azureScanner.terraformer.Import(TerraformImportMigrationGeneratorParams{
		Provider:       "azurerm",
		Resources:      []string{},
		AdditionalArgs: []string{fmt.Sprintf("--filter=resource_group=%s", filterValue)},
		Regions:        getValidRegions(azureScanner.CloudRegions, terraformValueObjects.AzureRegions, defaultAzureRegions),
		IsCompact:      true,
	})

	if err != nil {
		return fmt.Errorf("[Scan] Error in terraformer.Import(): %v", err)
	}

	err = azureScanner.terraformer.UpdateState("azurerm")

	if err != nil {
		return fmt.Errorf("[Scan] Error in terraformer.UpdateState(): %v", err)
	}

	return nil
}

func (azureScanner *AzureScanner) configureEnvironment(env AzureEnvironment) error {
	logrus.Debugf("[AzureScanner][configureEnvironment] Configuring environment %v", env)

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
