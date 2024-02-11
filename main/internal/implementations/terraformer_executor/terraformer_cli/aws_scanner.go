package terraformercli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

var defaultAwsRegions = []string{"us-east-1"}

// AWSScanner implements the Scanner interface for use with AWS cloud environments.
type AWSScanner struct {
	// credential needed to scan an AWS cloud environment.
	credential terraformValueObjects.Credential

	// terraformer is the TerraformerCLI interface used to scan the AWS cloud environment.
	terraformer TerraformerCLI

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions []terraformValueObjects.CloudRegion `required:"true"`
}

// NewAWSScanner creates and returns a new instance of AWSScanner.
func NewAWSScanner(credential terraformValueObjects.Credential, cliConfig Config, cloudRegions []terraformValueObjects.CloudRegion) (Scanner, error) {
	return &AWSScanner{
		CloudRegions: cloudRegions,
		credential:   credential,
		terraformer:  newTerraformerCLI(cliConfig),
	}, nil
}

// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
func (awsScanner *AWSScanner) Scan(_ terraformValueObjects.Division, credential terraformValueObjects.Credential, _ ...string) error {
	logrus.Debugf("[AWSScanner][Scan] Scanning AWS account %v", credential)
	err := awsScanner.configureEnvironment(credential)
	if err != nil {
		return fmt.Errorf("[AWS Scanner] Error configuring environment %w", err)
	}

	err = awsScanner.terraformer.Import(TerraformImportMigrationGeneratorParams{
		Provider:       "aws",
		Resources:      []string{},
		AdditionalArgs: []string{"--profile="},
		Regions:        getValidRegions(awsScanner.CloudRegions, terraformValueObjects.AwsRegions, defaultAwsRegions),
		IsCompact:      true,
	})
	if err != nil {
		return fmt.Errorf("[Scan] Error in terraformer.Import(): %v", err)
	}

	err = awsScanner.terraformer.UpdateState("aws")
	if err != nil {
		return fmt.Errorf("[Scan] Error in terraformer.UpdateState(): %v", err)
	}

	return nil
}

// AWSEnvironment is a struct defining the credential values needed for authenticating with an AWS account.
type AWSEnvironment struct {
	AWSAccessKeyID     string `json:"awsAccessKeyID"`
	AWSSecretKeyAccess string `json:"awsSecretAccessKey"`
}

// configureEnvironment loads and sets as environment variables AWS credentials for a given AWS account.
func (awsScanner *AWSScanner) configureEnvironment(credential terraformValueObjects.Credential) error {
	logrus.Debugf("[AWSScanner][configure_environment] Configuring environment for AWS account %v", credential)
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
