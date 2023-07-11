package terraformerCLI

import (
	"encoding/json"
	"fmt"
	"os"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

var defaultAwsRegions = []string{"us-east-1"}

// AWSScanner implements the Scanner interface for use with AWS cloud environments.
type AWSScanner struct {
	// Config is the needed configuration of a mapping between Division name and the corresponding
	// Credential needed to access that environment.
	config map[terraformValueObjects.Division]terraformValueObjects.Credential

	// terraformer is the TerraformerCLI interface used to scan the AWS cloud environment.
	terraformer TerraformerCLI

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions []terraformValueObjects.CloudRegion `required:"true"`
}

// NewAWSScanner creates and returns a new instance of AWSScanner.
func NewAWSScanner(config map[terraformValueObjects.Division]terraformValueObjects.Credential, cliConfig Config, cloudRegions []terraformValueObjects.CloudRegion) (Scanner, error) {
	return &AWSScanner{
		CloudRegions: cloudRegions,
		config:       config,
		terraformer:  newTerraformerCLI(cliConfig),
	}, nil
}

// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
func (awsScanner *AWSScanner) Scan(project terraformValueObjects.Division, credential terraformValueObjects.Credential, options ...string) (terraformValueObjects.Path, error) {
	err := awsScanner.configureEnvironment(credential)
	if err != nil {
		return "", fmt.Errorf("[AWS Scanner] Error configuring environment %w", err)
	}

	path, err := awsScanner.terraformer.Import(TerraformImportMigrationGeneratorParams{
		Provider:       "aws",
		Division:       project,
		Resources:      []string{},
		AdditionalArgs: []string{"--profile="},
		Regions:        getValidRegions(awsScanner.CloudRegions, terraformValueObjects.AwsRegions, defaultAwsRegions),
		IsCompact:      true,
	})

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.Import(): %v", err)
	}

	err = awsScanner.terraformer.UpdateState("aws", string(path))

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.UpdateState(): %v", err)
	}

	return path, nil
}

// ScanAll wraps Scan to scan each division for the provider.
func (awsScanner *AWSScanner) ScanAll(options ...string) (*MultiScanResult, error) {
	fmt.Println("Scanning all specified AWS divisions.")
	scanMap := make(map[terraformValueObjects.Division]terraformValueObjects.Path)

	for div, credential := range awsScanner.config {
		path, err := awsScanner.Scan(div, credential)
		if err != nil {
			return nil, fmt.Errorf("[ScanAll] Error in awsScanner.Scan: %v", err)
		}
		scanMap[div] = path
	}

	return &MultiScanResult{scanMap}, nil
}

// AWSEnvironment is a struct defining the credential values needed for authenticating with an AWS account.
type AWSEnvironment struct {
	AWSAccessKeyID     string `json:"awsAccessKeyID"`
	AWSSecretKeyAccess string `json:"awsSecretAccessKey"`
}

// configureEnvironment loads and sets as environment variables AWS credentials for a given AWS account.
func (awsScanner *AWSScanner) configureEnvironment(credential terraformValueObjects.Credential) error {
	env := new(AWSEnvironment)
	err := json.Unmarshal([]byte(credential), &env)
	if err != nil {
		return fmt.Errorf("[aws_scanner][configure_environment][error unmarshalling credentials] %w", err)
	}

	err = os.Setenv("AWS_ACCESS_KEY_ID", env.AWSAccessKeyID)
	if err != nil {
		return fmt.Errorf("[aws_scanner][configure_environment][error setting access_key_id credential] %w", err)
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", env.AWSSecretKeyAccess)
	if err != nil {
		return fmt.Errorf("[aws_scanner][configure_environment][error setting secret_access_key credential] %w", err)
	}

	return nil
}
