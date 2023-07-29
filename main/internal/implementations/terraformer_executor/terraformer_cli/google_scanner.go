package terraformerCLI

import (
	"fmt"
	"os"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

var defaultGoogleRegions = []string{"us-east4"}

// GoogleScanner implements the Scanner interface for use with GCP cloud environments.
type GoogleScanner struct {
	// credential is the credential needed to scan a GCP cloud environment.
	credential terraformValueObjects.Credential

	// terraformer is the TerraformerCLI interface used to scan the GCP cloud environment.
	terraformer TerraformerCLI

	// CloudRegions represents the list of cloud regions that will be considered for inclusion in the import statement.
	CloudRegions []terraformValueObjects.CloudRegion `required:"true"`
}

// NewGoogleScanner creates and returns a new instance of GCPScanner.
func NewGoogleScanner(credential terraformValueObjects.Credential, cliConfig Config, cloudRegions []terraformValueObjects.CloudRegion) (Scanner, error) {
	return &GoogleScanner{
		CloudRegions: cloudRegions,
		credential:   credential,
		terraformer:  newTerraformerCLI(cliConfig),
	}, nil
}

// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
func (gcpScan *GoogleScanner) Scan(project terraformValueObjects.Division, credential terraformValueObjects.Credential, options ...string) (terraformValueObjects.Path, error) {
	_ = os.MkdirAll("credentials", 0660)

	err := os.WriteFile(fmt.Sprintf("credentials/google-%v.json", project), []byte(credential), 0400)

	if err != nil {
		return "", fmt.Errorf("[Scan] error saving credential file: %v", err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", fmt.Sprintf("credentials/google-%s.json", project))

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in setting GOOGLE_APPLICATION_CREDENTIALS value: %v", err)
	}

	projectsFlag := fmt.Sprintf("--projects=%v", project)
	path, err := gcpScan.terraformer.Import(TerraformImportMigrationGeneratorParams{
		Provider:       "google",
		Division:       project,
		Regions:        getValidRegions(gcpScan.CloudRegions, terraformValueObjects.GoogleRegions, defaultGoogleRegions),
		Resources:      []string{"us-east4", "global"},
		AdditionalArgs: []string{projectsFlag},
		IsCompact:      true,
	})

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.Import(): %v", err)
	}

	err = gcpScan.terraformer.UpdateState("google", string(path))

	if err != nil {
		return "", fmt.Errorf("[Scan] Error in terraformer.UpdateState(): %v", err)
	}

	return path, nil
}
